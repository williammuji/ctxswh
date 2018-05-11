package main

import (
	"io"
	"math/rand"
	"time"

	"github.com/golang/glog"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/williammuji/ctxswh/gatewaypb"
	"github.com/williammuji/ctxswh/util"

	"sync"
)

type session struct {
	//config
	gatewayServerAddr string
	opts              []grpc.DialOption
	userName          string
	passwd            string
	secs              int
	ppMsgSize         int
	gatewayStream     gatewaypb.GatewayService_ServeClient
	gatewayReqC       chan *gatewaypb.GatewayRequest
	wg                sync.WaitGroup
	cancel            context.CancelFunc
	startWg           *sync.WaitGroup
}

func NewSession(gatewayServerAddr string, opts []grpc.DialOption, userName, passwd string, secs, ppMsgSize int, startWg *sync.WaitGroup) *session {
	return &session{
		gatewayServerAddr: gatewayServerAddr,
		opts:              opts,
		userName:          userName,
		passwd:            passwd,
		secs:              secs,
		ppMsgSize:         ppMsgSize,
		gatewayReqC:       make(chan *gatewaypb.GatewayRequest),
		startWg:           startWg,
	}
}

func (sess *session) start() {
	defer sess.startWg.Done()

	gatewayConn, err := grpc.Dial(sess.gatewayServerAddr, sess.opts...)
	if err != nil {
		glog.Fatalf("session.start failed sess:%#v err:%s", sess, err)
		return
	}
	defer gatewayConn.Close()
	glog.Infof("session.start grpc.Dial success sess:%#v gatewayConn:%#v", sess, gatewayConn)

	gatewayClient := gatewaypb.NewGatewayServiceClient(gatewayConn)
	gatewayCtx, cancel := context.WithCancel(context.Background())
	gatewayStream, err := gatewayClient.Serve(gatewayCtx)
	if err != nil {
		cancel()
		glog.Fatalf("session.start gatewayClient.Serve failed sess:%#v gatewayConn:%#v gatewayClient:%#v gatewayCtx:%#v", sess, gatewayConn, gatewayClient, gatewayCtx)
		return
	}
	glog.Infof("session.start gatewayClient.Serve success sess:%#v gatewayConn:%#v gatewayClient:%#v gatewayCtx:%#v gatewayStream:%#v", sess, gatewayConn, gatewayClient, gatewayCtx, gatewayStream)

	sess.gatewayStream = gatewayStream
	sess.cancel = cancel
	sess.wg.Add(1)
	go func() {
		sess.wg.Wait()
		close(sess.gatewayReqC)
		glog.Infof("session.start sess.wg.Wait sess:%#v", sess)
	}()

	go sess.recvLoop()
	go sess.login()
	go sess.scheduleTimeout()
	sess.sendLoop()
}

func (sess *session) recvLoop() {
	defer sess.wg.Done()

	for {
		in, err := sess.gatewayStream.Recv()
		if err == io.EOF {
			glog.Infof("session.recvLoop sess.gatewayStream.Recv io.EOF sess:%#v", sess)
			return
		}
		if err != nil {
			glog.Infof("session.recvLoop sess.gatewayStream.Recv sess:%#v err:%s", sess, err)
			return
		}

		switch resp := in.MessageResponse.(type) {
		case *gatewaypb.GatewayResponse_LoginResponse:
			if resp.LoginResponse.Status != gatewaypb.LoginResponse_SUCCESS {
				glog.Infof("session.recvLoop resp.LoginResponse.Status NOT SUCCESS sess:%#v resp:%#v", sess, resp)
				return
			}
			glog.Infof("session.recvLoop recv GatewayResponse_LoginResponse SUCCESS sess:%#v resp:%#v", sess, resp)

			//start pingpong
			size := uint16(rand.Int31n(int32(sess.ppMsgSize)))
			req := &gatewaypb.GatewayRequest{
				MessageRequest: &gatewaypb.GatewayRequest_PingpongRequest{
					PingpongRequest: &gatewaypb.PingpongData{
						Size: int32(size),
						Data: util.RandBytes(size),
					},
				},
			}
			err = sess.writeGatewayChan(req)
			if err != nil {
				glog.Infof("session.recvLoop sess.writeGatewayChan error sess:%#v req:%#v err:%s", sess, req, err)
				return
			}

			//start sub
			subReq := &gatewaypb.GatewayRequest{
				MessageRequest: &gatewaypb.GatewayRequest_SubRequest{
					SubRequest: &gatewaypb.SubRequest{
						Topic: "foo",
					},
				},
			}
			err = sess.writeGatewayChan(subReq)
			if err != nil {
				glog.Infof("session.recvLoop sess.writeGatewayChan error sess:%#v req:%#v err:%s", sess, req, err)
				return
			}

			//in future unsub
			go func() {
				secs := rand.Intn(sess.secs / 2)
				select {
				case <-time.After(time.Duration(secs) * time.Second):
					unsubReq := &gatewaypb.GatewayRequest{
						MessageRequest: &gatewaypb.GatewayRequest_UnsubRequest{
							UnsubRequest: &gatewaypb.UnsubRequest{
								Topic: "foo",
							},
						},
					}
					err = sess.writeGatewayChan(unsubReq)
					if err != nil {
						glog.Infof("session.recvLoop sess.writeGatewayChan error sess:%#v unsubReq:%#v err:%s", sess, unsubReq, err)
						return
					}
					glog.Infof("session.recvLoop sess.writeGatewayChan sess:%#v unsubReq:%#v", sess, unsubReq)
				}
			}()
		case *gatewaypb.GatewayResponse_PingpongResponse:
			glog.Infof("session.recvLoop recv GatewayResponse_PingpongResponse sess:%#v resp:%#v", sess, resp)

			go func() {
				select {
				case <-time.After(1 * time.Second):
					//pingpong every second
					req := &gatewaypb.GatewayRequest{
						MessageRequest: &gatewaypb.GatewayRequest_PingpongRequest{
							PingpongRequest: resp.PingpongResponse,
						},
					}
					err = sess.writeGatewayChan(req)
					if err != nil {
						glog.Infof("session.recvLoop writeGatewayChan error sess:%#v req:%#v err:%v", sess, req, err)
						return
					}
					glog.Infof("session.recvLoop time.After 1 second pingpong  sess:%#v req:%#v", sess, req)

					//pub every second
					t := time.Now()
					ts := sess.userName + "_" + t.String()
					pubReq := &gatewaypb.GatewayRequest{
						MessageRequest: &gatewaypb.GatewayRequest_PubRequest{
							PubRequest: &gatewaypb.PubRequest{
								Topic:   "foo",
								Content: []byte(ts),
							},
						},
					}
					err = sess.writeGatewayChan(pubReq)
					if err != nil {
						glog.Infof("session.recvLoop writeGatewayChan error sess:%#v pubReq:%#v err:%v", sess, pubReq, err)
						return
					}
					glog.Infof("session.recvLoop time.After 1 second pub  sess:%#v pubReq:%#v", sess, pubReq)
				}
			}()
		case *gatewaypb.GatewayResponse_PubResponse:
			for i, d := range resp.PubResponse.Datas {
				glog.Infof("session.recvLoop recv GatewayResponse_PubResponse sess:%#v resp:%#v content:%s pubTime:%d i:%d", sess, resp, string(d.Content), d.PubTime, i)
			}
		default:
			glog.Info(grpc.Errorf(codes.InvalidArgument, "invalid MessageResponse: %v", in.MessageResponse))
		}
	}
}

func (sess *session) sendLoop() error {
	for req := range sess.gatewayReqC {
		if err := sess.gatewayStream.Send(req); err != nil {
			glog.Infof("session.sendLoop sess.gatewayStream.Send error sess:%#v req:%#v err:%v", sess, req, err)
			return err
		}
	}
	glog.Infof("session.sendLoop exit sess:%#v", sess)
	return nil
}

func (sess *session) writeGatewayChan(req *gatewaypb.GatewayRequest) error {
	select {
	case sess.gatewayReqC <- req:
	case <-sess.gatewayStream.Context().Done():
		glog.Infof("session.writeGatewayChan sess.gatewayStream.Context().Done() sess:%#v req:%#v", sess, req)
		return sess.gatewayStream.Context().Err()
	}
	return nil
}

func (sess *session) login() {
	req := &gatewaypb.GatewayRequest{
		MessageRequest: &gatewaypb.GatewayRequest_LoginRequest{
			LoginRequest: &gatewaypb.LoginRequest{
				Username: sess.userName,
				Password: sess.passwd,
			},
		},
	}
	err := sess.writeGatewayChan(req)
	if err != nil {
		glog.Infof("session.login sess.writeGatewayChan error sess:%#v req:%#v err:%v", sess, req, err)
	} else {
		glog.Infof("session.login sess.writeGatewayChan SUCCESS sess:%#v req:%#v", sess, req)
	}
}

func (sess *session) scheduleTimeout() {
	select {
	case <-time.After(time.Duration(sess.secs) * time.Second):
		glog.Infof("session.scheduleTimeout sess:%#v secs:%d", sess, sess.secs)
		sess.cancel()
	}
}
