package main

import (
	"io"
	"sync"

	"github.com/golang/glog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/williammuji/ctxswh/gamepb"
)

type gameServer struct {
}

func NewGameServer() (*gameServer, error) {
	return &gameServer{}, nil
}

type session struct {
	gs *gameServer

	gRPCStream gamepb.GameService_ServeServer
	wg         sync.WaitGroup
	respC      chan *gamepb.GameResponse
}

func (gs *gameServer) Serve(stream gamepb.GameService_ServeServer) error {
	sess := session{
		gs:         gs,
		gRPCStream: stream,
		respC:      make(chan *gamepb.GameResponse),
	}

	sess.wg.Add(1)
	go func() {
		sess.wg.Wait()
		close(sess.respC)
		glog.Infof("gameServer.Serve sess.wg.Wait gameServer:%#v sess:%#v", gs, sess)
	}()

	go sess.recvLoop()
	return sess.sendLoop()
}

func (sess *session) recvLoop() error {
	defer sess.wg.Done()

	for {
		in, err := sess.gRPCStream.Recv()
		if err == io.EOF {
			glog.Infof("session.recvLoop sess.gRPCStream.Recv io.EOF sess:%#v", sess)
			return nil
		}
		if err != nil {
			glog.Infof("session.recvLoop sess.gRPCStream.Recv error sess:%#v err:%v", sess, err)
			return err
		}

		switch req := in.MessageRequest.(type) {
		case *gamepb.GameRequest_PingpongRequest:
			glog.Infof("session.recvLoop recv GameRequest_PingpongRequest sess:%#v req:%#v", sess, req)

			out := &gamepb.GameResponse{
				MessageResponse: &gamepb.GameResponse_PingpongResponse{
					PingpongResponse: req.PingpongRequest,
				},
			}
			err := sess.writeGatewayChan(out)
			if err != nil {
				glog.Infof("session.recvLoop writeGatewayChan error sess:%#v out:%#v err:%v", sess, out, err)
				return nil
			}
		default:
			glog.Infof("session.recvLoop recv invalid MessageRequest sess:%#v %v", sess, in.MessageRequest)
			return grpc.Errorf(codes.InvalidArgument, "invalid MessageRequest: %v", in.MessageRequest)
		}
	}
}

func (sess *session) sendLoop() error {
	for resp := range sess.respC {
		if err := sess.gRPCStream.Send(resp); err != nil {
			glog.Infof("session.sendLoop sess.gRPCStream.Send error sess:%#v resp:%#v err:%v", sess, resp, err)
			return err
		}
	}
	glog.Infof("session.sendLoop exit sess:%#v", sess)
	return nil
}

func (sess *session) writeGatewayChan(resp *gamepb.GameResponse) error {
	select {
	case sess.respC <- resp:
	case <-sess.gRPCStream.Context().Done():
		glog.Infof("session.writeGatewayChan sess.gRPCStream.Context().Done() sess:%#v", sess)
		return sess.gRPCStream.Context().Err()
	}
	return nil
}
