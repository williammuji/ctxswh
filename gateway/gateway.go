package main

import (
	"io"
	"sync"

	"github.com/golang/glog"

	"golang.org/x/net/context"

	"github.com/williammuji/ctxswh/authpb"
	"github.com/williammuji/ctxswh/chatpb"
	"github.com/williammuji/ctxswh/gamepb"
	"github.com/williammuji/ctxswh/gatewaypb"
	"github.com/williammuji/ctxswh/loginpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

type gatewayServer struct {
	authConn    *grpc.ClientConn
	gameConn    *grpc.ClientConn
	chatConn    *grpc.ClientConn
	loginConn   *grpc.ClientConn
	authClient  authpb.AuthServiceClient
	gameClient  gamepb.GameServiceClient
	chatClient  chatpb.ChatServiceClient
	loginClient loginpb.LoginServiceClient
	sync.RWMutex
	sequence     uint64
	sessionCount uint64
	loginStream  loginpb.LoginService_ServeClient
	loginReqChan chan *loginpb.LoginRequest
	loginWg      sync.WaitGroup
}

func NewGatewayServer(authServerAddr, gameServerAddr, chatServerAddr, loginServerAddr string) (*gatewayServer, error) {
	authConn, err := grpc.Dial(authServerAddr, grpc.WithInsecure())
	if err != nil {
		glog.Fatalf("NewGatewayServer grpc.Dial authServerAddr:%s err:%v", authServerAddr, err)
		return nil, err
	}
	gameConn, err := grpc.Dial(gameServerAddr, grpc.WithInsecure())
	if err != nil {
		glog.Fatalf("NewGatewayServer grpc.Dial gameServerAddr:%s err:%v", gameServerAddr, err)
		return nil, err
	}
	chatConn, err := grpc.Dial(chatServerAddr, grpc.WithInsecure())
	if err != nil {
		glog.Fatalf("NewGatewayServer grpc.Dial chatServerAddr:%s err:%v", chatServerAddr, err)
		return nil, err
	}
	loginConn, err := grpc.Dial(loginServerAddr, grpc.WithInsecure())
	if err != nil {
		glog.Fatalf("NewGatewayServer grpc.Dial loginServerAddr:%s err:%v", loginServerAddr, err)
		return nil, err
	}

	gs := &gatewayServer{
		authConn:     authConn,
		authClient:   authpb.NewAuthServiceClient(authConn),
		gameConn:     gameConn,
		gameClient:   gamepb.NewGameServiceClient(gameConn),
		chatConn:     chatConn,
		chatClient:   chatpb.NewChatServiceClient(chatConn),
		loginConn:    loginConn,
		loginClient:  loginpb.NewLoginServiceClient(loginConn),
		sequence:     0,
		sessionCount: 0,
	}
	go gs.loginLoop()

	return gs, nil
}

func (gs *gatewayServer) CloseConn() {
	gs.authConn.Close()
	gs.gameConn.Close()
	gs.chatConn.Close()
	gs.loginConn.Close()

	glog.Infof("gatewayServer.CloseConn %#v", gs)
}

func (gs *gatewayServer) incSequence() uint64 {
	var seq uint64
	gs.Lock()
	gs.sequence++
	gs.sessionCount++
	seq = gs.sequence
	gs.Unlock()
	return seq
}

func (gs *gatewayServer) decSessionCount() uint64 {
	var sc uint64
	gs.Lock()
	gs.sessionCount--
	sc = gs.sessionCount
	gs.Unlock()
	return sc
}

func (gs *gatewayServer) getSequence() uint64 {
	var seq uint64
	gs.Lock()
	seq = gs.sequence
	gs.Unlock()
	return seq
}

func (gs *gatewayServer) getSessionCount() uint64 {
	var sc uint64
	gs.Lock()
	sc = gs.sessionCount
	gs.Unlock()
	return sc
}

func (gs *gatewayServer) notifySessionCount() {
	out := &loginpb.LoginRequest{
		MessageRequest: &loginpb.LoginRequest_IpCount{
			IpCount: &loginpb.IpCount{
				Count: gs.getSessionCount(),
			},
		},
	}
	err := gs.writeLoginChan(out)
	if err != nil {
		glog.Infof("gatewayServer.notifySessionCount gs.writeLoginChan error gs:%#v out:%#v err:%v", gs, out, err)
	}
}

func (gs *gatewayServer) loginLoop() {
	loginStream, err := gs.loginClient.Serve(context.Background())
	if err != nil {
		glog.Infof("gatewayServer.loginLoop gs.loginClient.Serve error:%#v", err)
		return
	}

	gs.loginStream = loginStream
	gs.loginReqChan = make(chan *loginpb.LoginRequest)

	gs.loginWg.Add(1)
	go func() {
		gs.loginWg.Wait()
		close(gs.loginReqChan)
	}()
	go gs.recvLoginLoop()
	go gs.notifySessionCount()
	gs.sendLoginLoop()
}

func (gs *gatewayServer) recvLoginLoop() {
	defer gs.loginWg.Done()

	for {
		in, err := gs.loginStream.Recv()
		if err == io.EOF {
			glog.Infof("gatewayServer.recvLoginLoop gs.loginStream.Recv io.EOF gs:%#v", gs)
			return
		}
		if err != nil {
			glog.Infof("gatewayServer.recvLoginLoop gs.loginStream.Recv error gs:%#v err:%v", gs, err)
			return
		}

		switch resp := in.MessageResponse.(type) {
		case *loginpb.LoginResponse_Empty:
			glog.Infof("gatewayServer.recvLoginLoop recv LoginResponse_Empty gs:%#v in:%#v resp:%#v", gs, in, resp)
		default:
			glog.Infof("gatewayServer.recvLoginLoop recv invalid MessageResponse:%v gs:%#v", in.MessageResponse, gs)
			return
		}
	}
}

func (gs *gatewayServer) sendLoginLoop() error {
	for req := range gs.loginReqChan {
		if err := gs.loginStream.Send(req); err != nil {
			glog.Infof("gatewayServer.sendLoginLoop gs.loginStream.Send error gs:%#v req:%#v err:%v", gs, req, err)
			return err
		}
	}
	return nil
}

func (gs *gatewayServer) writeLoginChan(req *loginpb.LoginRequest) error {
	select {
	case gs.loginReqChan <- req:
	case <-gs.loginStream.Context().Done():
		glog.Infof("gatewayServer.writeLoginChan gs.loginStream.Context().Done() error gs:%#v req:%#v", gs, req)
		return gs.loginStream.Context().Err()
	}
	return nil
}

func (gs *gatewayServer) Serve(stream gatewaypb.GatewayService_ServeServer) error {
	in, err := stream.Recv()
	if err == io.EOF {
		glog.Infof("gatewayServer.Serve stream.Recv io.EOF gatewayServer:%#v stream:%#v err:%v", gs, stream, err)
		return nil
	}
	if err != nil {
		glog.Infof("gatewayServer.Serve stream.Recv err!=nil gatewayServer:%#v stream:%#v err:%v", gs, stream, err)
		return err
	}

	switch req := in.MessageRequest.(type) {
	case *gatewaypb.GatewayRequest_LoginRequest:
		seq := gs.incSequence()
		gs.notifySessionCount()
		glog.Infof("gatewayServer.Serve recv GatewayRequest_LoginRequest before Auth gatewayServer:%#v stream:%#v req:%#v seq:%d", gs, stream, req, seq)

		sess := &session{
			seq:           seq,
			name:          req.LoginRequest.Username,
			gs:            gs,
			gatewayStream: stream,
		}

		authStart := make(chan bool)
		go func() {
			<-authStart

			out := &authpb.AuthRequest{
				MessageRequest: &authpb.AuthRequest_LoginRequest{
					LoginRequest: req.LoginRequest,
				},
			}
			err = sess.writeAuthChan(out)
			if err != nil {
				glog.Infof("gatewayServer.Serve sess.writeAuthChan error sess:%#v out:%#v err:%v", sess, out, err)
			}
		}()

		sess.authLoop(authStart)
		return nil
	default:
		glog.Infof("gatewayServer.Serve recv invalid MessageRequest gatewayServer:%#v stream:%#v req:%#v", gs, stream, req)
		return grpc.Errorf(codes.InvalidArgument, "grpc: must be LoginRequest msg")
	}
}

type session struct {
	seq             uint64
	name            string
	gs              *gatewayServer
	gatewayStream   gatewaypb.GatewayService_ServeServer
	gatewayWg       sync.WaitGroup
	gatewayRespChan chan *gatewaypb.GatewayResponse
	authWg          sync.WaitGroup
	authStream      authpb.AuthService_ServeClient
	authReqChan     chan *authpb.AuthRequest
	gameWg          sync.WaitGroup
	gameStream      gamepb.GameService_ServeClient
	gameReqChan     chan *gamepb.GameRequest
	chatWg          sync.WaitGroup
	chatStream      chatpb.ChatService_ServeClient
	chatReqChan     chan *chatpb.ChatRequest
}

func (sess *session) gatewayLoop() {
	sess.gatewayRespChan = make(chan *gatewaypb.GatewayResponse)

	sess.gatewayWg.Add(1)
	go func() {
		sess.gatewayWg.Wait()
		close(sess.gatewayRespChan)
	}()
	go sess.recvGatewayLoop()
	sess.sendGatewayLoop()
}

func (sess *session) recvGatewayLoop() {
	defer sess.gatewayWg.Done()

	for {
		in, err := sess.gatewayStream.Recv()
		if err == io.EOF {
			glog.Infof("session.recvGatewayLoop sess.gatewayStream.Recv io.EOF sess:%#v", sess)
			return
		}
		if err != nil {
			glog.Infof("session.recvGatewayLoop sess.gatewayStream.Recv error sess:%#v error:%v", sess, err)
			return
		}

		switch req := in.MessageRequest.(type) {
		case *gatewaypb.GatewayRequest_PingpongRequest:
			out := &gamepb.GameRequest{
				MessageRequest: &gamepb.GameRequest_PingpongRequest{
					PingpongRequest: req.PingpongRequest,
				},
			}
			err = sess.writeGameChan(out)
			if err != nil {
				glog.Infof("session.recvGatewayLoop recv GatewayRequest_PingpongRequest writeGameChan error sess:%#v out:%#v err:%v", sess, out, err)
			}
		case *gatewaypb.GatewayRequest_SubRequest:
			out := &chatpb.ChatRequest{
				MessageRequest: &chatpb.ChatRequest_SubRequest{
					SubRequest: req.SubRequest,
				},
			}
			err = sess.writeChatChan(out)
			if err != nil {
				glog.Infof("session.recvGatewayLoop recv GatewayRequest_SubRequest writeGameChan error sess:%#v out:%#v err:%v", sess, out, err)
			}
		case *gatewaypb.GatewayRequest_UnsubRequest:
			out := &chatpb.ChatRequest{
				MessageRequest: &chatpb.ChatRequest_UnsubRequest{
					UnsubRequest: req.UnsubRequest,
				},
			}
			err = sess.writeChatChan(out)
			if err != nil {
				glog.Infof("session.recvGatewayLoop recv GatewayRequest_UnsubRequest writeGameChan error sess:%#v out:%#v err:%v", sess, out, err)
			}
		case *gatewaypb.GatewayRequest_PubRequest:
			out := &chatpb.ChatRequest{
				MessageRequest: &chatpb.ChatRequest_PubRequest{
					PubRequest: req.PubRequest,
				},
			}
			err = sess.writeChatChan(out)
			if err != nil {
				glog.Infof("session.recvGatewayLoop recv GatewayRequest_PubRequest writeGameChan error sess:%#v out:%#v err:%v", sess, out, err)
			}
		default:
			glog.Infof("session.recvGatewayLoop recv invalid MessageRequest:%v sess:%#v", in.MessageRequest, sess)
			return
		}
	}
}

func (sess *session) sendGatewayLoop() error {
	for resp := range sess.gatewayRespChan {
		if err := sess.gatewayStream.Send(resp); err != nil {
			glog.Infof("session.sendGatewayLoop sess.gatewayStream.Send error sess:%#v resp:%#v err:%v", sess, resp, err)
			return err
		}
	}
	return nil
}

func (sess *session) writeGatewayChan(resp *gatewaypb.GatewayResponse) error {
	select {
	case sess.gatewayRespChan <- resp:
	case <-sess.gatewayStream.Context().Done():
		glog.Infof("session.writeGatewayChan sess.gatewayStream.Context().Done() error sess:%#v resp:%#v", sess, resp)
		return sess.gatewayStream.Context().Err()
	}
	return nil
}

func (sess *session) authLoop(authStart chan<- bool) {
	authCtx := metadata.NewContext(context.Background(), metadata.New(map[string]string{"userName": sess.name}))
	authStream, err := sess.gs.authClient.Serve(authCtx)
	if err != nil {
		glog.Infof("session.authLoop sess.gs.authClient.Serve error sess:%#v authCtx:%#v seq:%#v", sess, authCtx, sess.seq)
		return
	}

	sess.authStream = authStream
	sess.authReqChan = make(chan *authpb.AuthRequest)

	sess.authWg.Add(1)
	go func() {
		sess.authWg.Wait()
		close(sess.authReqChan)
	}()
	go sess.recvAuthLoop()
	go func() {
		authStart <- true
	}()
	sess.sendAuthLoop()
}

func (sess *session) recvAuthLoop() {
	defer sess.authWg.Done()

	for {
		in, err := sess.authStream.Recv()
		if err == io.EOF {
			glog.Infof("session.recvAuthLoop sess.authStream.Recv io.EOF sess:%#v", sess)
			return
		}
		if err != nil {
			glog.Infof("session.recvAuthLoop sess.authStream.Recv error sess:%#v err:%v", sess, err)
			return
		}

		switch resp := in.MessageResponse.(type) {
		case *authpb.AuthResponse_LoginResponse:
			glog.Infof("session.recvAuthLoop recv AuthResponse_LoginResponse sess:%#v in:%#v resp%#v", sess, in, resp)

			out := &gatewaypb.GatewayResponse{
				MessageResponse: &gatewaypb.GatewayResponse_LoginResponse{
					LoginResponse: resp.LoginResponse,
				},
			}
			go sess.gatewayStream.Send(out)

			if resp.LoginResponse.Status == gatewaypb.LoginResponse_SUCCESS {
				go sess.gameLoop()
				go sess.chatLoop()
				sess.gatewayLoop()
			}

			return
		default:
			glog.Infof("session.recvAuthLoop recv invalid MessageResponse:%v sess:%#v", in.MessageResponse, sess)
			return
		}
	}
}

func (sess *session) sendAuthLoop() error {
	for req := range sess.authReqChan {
		if err := sess.authStream.Send(req); err != nil {
			glog.Infof("session.sendAuthLoop sess.authStream.Send error sess:%#v req:%#v err:%v", sess, req, err)
			return err
		}
	}
	return nil
}

func (sess *session) writeAuthChan(req *authpb.AuthRequest) error {
	select {
	case sess.authReqChan <- req:
	case <-sess.authStream.Context().Done():
		glog.Infof("session.writeAuthChan sess.authStream.Context().Done() error sess:%#v req:%#v", sess, req)
		return sess.authStream.Context().Err()
	}
	return nil
}

func (sess *session) gameLoop() {
	gameCtx := metadata.NewContext(context.Background(), metadata.New(map[string]string{"userName": sess.name}))
	gameStream, err := sess.gs.gameClient.Serve(gameCtx)
	if err != nil {
		glog.Infof("session.gameLoop sess.gs.gameClient.Serve error sess:%#v gameCtx:%#v seq:%v name:%v", sess, gameCtx, sess.seq, sess.name)
		return
	}

	sess.gameStream = gameStream
	sess.gameReqChan = make(chan *gamepb.GameRequest)

	sess.gameWg.Add(1)
	go func() {
		sess.gameWg.Wait()
		close(sess.gameReqChan)
	}()
	go sess.recvGameLoop()
	sess.sendGameLoop()
}

func (sess *session) recvGameLoop() {
	defer sess.gameWg.Done()

	for {
		in, err := sess.gameStream.Recv()
		if err == io.EOF {
			glog.Infof("session.recvGameLoop sess.gameStream.Recv io.EOF sess:%#v", sess)
			return
		}
		if err != nil {
			glog.Infof("session.recvGameLoop sess.gameStream.Recv error sess:%#v err:%v", sess, err)
			return
		}

		switch resp := in.MessageResponse.(type) {
		case *gamepb.GameResponse_PingpongResponse:
			glog.Infof("session.recvGameLoop recv GameResponse_PingpongResponse sess:%#v in:%#v resp%#v", sess, in, resp)

			out := &gatewaypb.GatewayResponse{
				MessageResponse: &gatewaypb.GatewayResponse_PingpongResponse{
					PingpongResponse: resp.PingpongResponse,
				},
			}
			err = sess.writeGatewayChan(out)
			if err != nil {
				glog.Infof("session.recvGameLoop sess.writeGatewayChan error sess:%#v resp:%#v out:%#v err:%s", sess, resp, out, err)
			}
		default:
			glog.Infof("session.recvGameLoop recv invalid MessageResponse:%v sess:%#v", in.MessageResponse, sess)
			return
		}
	}
}

func (sess *session) sendGameLoop() error {
	for req := range sess.gameReqChan {
		if err := sess.gameStream.Send(req); err != nil {
			glog.Infof("session.sendGameLoop sess.gameStream.Send error sess:%#v req:%#v err:%v", sess, req, err)
			return err
		}
	}
	return nil
}

func (sess *session) writeGameChan(req *gamepb.GameRequest) error {
	select {
	case sess.gameReqChan <- req:
	case <-sess.gameStream.Context().Done():
		glog.Infof("session.writeGameChan sess.gameStream.Context().Done() error sess:%#v req:%#v", sess, req)
		return sess.gameStream.Context().Err()
	}
	return nil
}

func (sess *session) chatLoop() {
	chatCtx := metadata.NewContext(context.Background(), metadata.New(map[string]string{"userName": sess.name}))
	chatStream, err := sess.gs.chatClient.Serve(chatCtx)
	if err != nil {
		glog.Infof("session.chatLoop sess.gs.chatClient.Serve error sess:%#v chatCtx:%#v seq:%v name:%v", sess, chatCtx, sess.seq, sess.name)
		return
	}

	sess.chatStream = chatStream
	sess.chatReqChan = make(chan *chatpb.ChatRequest)

	sess.chatWg.Add(1)
	go func() {
		sess.chatWg.Wait()
		close(sess.chatReqChan)
	}()
	go sess.recvChatLoop()
	sess.sendChatLoop()
}

func (sess *session) recvChatLoop() {
	defer sess.chatWg.Done()

	for {
		in, err := sess.chatStream.Recv()
		if err == io.EOF {
			glog.Infof("session.recvChatLoop sess.chatStream.Recv io.EOF sess:%#v", sess)
			return
		}
		if err != nil {
			glog.Infof("session.recvChatLoop sess.chatStream.Recv error sess:%#v err:%v", sess, err)
			return
		}

		switch resp := in.MessageResponse.(type) {
		case *chatpb.ChatResponse_PubResponse:
			glog.Infof("session.recvChatLoop recv ChatResponse_PubResponse sess:%#v in:%#v resp:%#v", sess, in, resp)

			out := &gatewaypb.GatewayResponse{
				MessageResponse: &gatewaypb.GatewayResponse_PubResponse{
					PubResponse: resp.PubResponse,
				},
			}
			err = sess.writeGatewayChan(out)
			if err != nil {
				glog.Infof("session.recvChatLoop recv ChatResponse_PubResponse sess.writeGatewayChan sess:%#v in:%#v resp:%#v out:%#v err:%v", sess, in, resp, out, err)
			}
		default:
			glog.Infof("session.recvChatLoop recv invalid MessageResponse:%v sess:%#v", in.MessageResponse, sess)
			return
		}
	}
}

func (sess *session) sendChatLoop() error {
	for req := range sess.chatReqChan {
		if err := sess.chatStream.Send(req); err != nil {
			glog.Infof("session.sendChatLoop sess.chatStream.Send error sess:%#v req:%#v err:%v", sess, req, err)
			return err
		}
	}
	return nil
}

func (sess *session) writeChatChan(req *chatpb.ChatRequest) error {
	select {
	case sess.chatReqChan <- req:
	case <-sess.chatStream.Context().Done():
		glog.Infof("session.writeChatChan sess.chatStream.Context().Done() error sess:%#v req:%#v", sess, req)
		return sess.chatStream.Context().Err()
	}
	return nil
}
