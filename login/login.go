package main

import (
	"io"
	"math"
	"net"
	"sync"

	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/williammuji/ctxswh/loginpb"
	"google.golang.org/grpc/peer"
)

type session struct {
	ls          *loginServer
	gGRPCStream loginpb.LoginService_ServeServer
	wg          sync.WaitGroup
	respC       chan *loginpb.LoginResponse
	ip          string
}

type loginServer struct {
	sync.RWMutex
	ipCount map[string]uint64
}

func NewLoginServer() (*loginServer, error) {
	return &loginServer{
		ipCount: make(map[string]uint64),
	}, nil
}

func (ls *loginServer) Serve(stream loginpb.LoginService_ServeServer) error {
	sess := &session{
		ls:          ls,
		gGRPCStream: stream,
		respC:       make(chan *loginpb.LoginResponse),
	}

	sess.wg.Add(1)
	go func() {
		sess.wg.Wait()
		close(sess.respC)
		glog.Infof("loginServer.Serve sess.wg.Wait loginServer:%#v sess:%#v", ls, sess)
	}()

	go sess.recvLoop()
	return sess.sendLoop()
}

func (ls *loginServer) setIpCount(ip string, count uint64) {
	ls.Lock()
	defer ls.Unlock()

	ls.ipCount[ip] = count
}

func (ls *loginServer) removeIpCount(ip string) {
	ls.Lock()
	defer ls.Unlock()

	delete(ls.ipCount, ip)
}

func (ls *loginServer) getMinCountIp() string {
	ls.Lock()
	defer ls.Unlock()

	var ip string
	var minCount uint64 = math.MaxUint64
	for k, v := range ls.ipCount {
		if v <= minCount {
			ip = k
			minCount = v
		}
	}
	return ip
}

func (sess *session) recvLoop() error {
	defer sess.wg.Done()

	for {
		in, err := sess.gGRPCStream.Recv()
		if err == io.EOF {
			glog.Infof("session.recvLoop sess.gGRPCStream.Recv io.EOF sess:%#v", sess)
			return nil
		}
		if err != nil {
			glog.Infof("session.recvLoop sess.gGRPCStream.Recv error sess:%#v err:%v", sess, err)
			return err
		}

		switch req := in.MessageRequest.(type) {
		case *loginpb.LoginRequest_IpCount:
			pr, ok := peer.FromContext(sess.gGRPCStream.Context())
			if !ok {
				glog.Infof("session.recvLoop peer.FromContext error sess:%#v req:%#v", sess, req)
				return nil
			}
			if pr.Addr == net.Addr(nil) {
				glog.Infof("session.recvLoop peer.FromContext pr.Addr=nil sess:%#v req:%#v", sess, req)
				return nil
			}

			sess.ip = pr.Addr.String()
			sess.ls.setIpCount(pr.Addr.String(), req.IpCount.GetCount())
			glog.Infof("session.recvLoop recv loginpb.LoginRequest_IpCount sess:%#v req:%#v in:%#v sess.gGRPCStream:%#v IP:%s count:%d", sess, req, in, sess.gGRPCStream, pr.Addr.String(), req.IpCount.GetCount())
		default:
			glog.Infof("session.recvLoop recv invalid MessageRequest sess:%#v req:%#v in:%#v sess.gGRPCStream:%#v", sess, req, in, sess.gGRPCStream)
			return grpc.Errorf(codes.InvalidArgument, "invalid MessageRequest: %v", in.MessageRequest)
		}
	}
}

func (sess *session) sendLoop() error {
	for resp := range sess.respC {
		if err := sess.gGRPCStream.Send(resp); err != nil {
			glog.Infof("session.sendLoop sess.gGRPCStream.Send error sess:%#v resp:%#v err:%v", sess, resp, err)
			return err
		}
	}
	glog.Infof("session.sendLoop exit sess:%#v", sess)
	sess.ls.removeIpCount(sess.ip)
	return nil
}

func (sess *session) writeChan(resp *loginpb.LoginResponse) error {
	select {
	case sess.respC <- resp:
	case <-sess.gGRPCStream.Context().Done():
		glog.Infof("session.writeChan sess.gGRPCStream.Context().Done() sess:%#v resp:%#v", sess, resp)
		return sess.gGRPCStream.Context().Err()
	}
	return nil
}
