package main

import (
	"io"
	"sync"

	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/williammuji/ctxswh/authpb"
	"github.com/williammuji/ctxswh/gatewaypb"

	"golang.org/x/crypto/bcrypt"
)

type authServer struct {
}

func NewAuthServer() (*authServer, error) {
	return &authServer{}, nil
}

func (as *authServer) Serve(stream authpb.AuthService_ServeServer) error {
	sess := session{
		as:         as,
		gRPCStream: stream,
		respC:      make(chan *authpb.AuthResponse),
	}

	sess.wg.Add(1)
	go func() {
		sess.wg.Wait()
		close(sess.respC)
		glog.Infof("authServer.Serve sess.wg.Wait authServer:%#v sess:%#v", as, sess)
	}()

	go sess.recvLoop()
	return sess.sendLoop()
}

type session struct {
	as *authServer

	gRPCStream authpb.AuthService_ServeServer
	wg         sync.WaitGroup
	respC      chan *authpb.AuthResponse
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
		case *authpb.AuthRequest_LoginRequest:
			glog.Infof("session.recvLoop recv AuthRequest_LoginRequest sess:%#v req:%#v", sess, req)

			resp, err := sess.Auth(req.LoginRequest)
			out := &authpb.AuthResponse{
				MessageResponse: &authpb.AuthResponse_LoginResponse{
					LoginResponse: resp,
				},
			}
			sess.writeAuthChan(out)
			return err
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

func (sess *session) writeAuthChan(resp *authpb.AuthResponse) error {
	select {
	case sess.respC <- resp:
	case <-sess.gRPCStream.Context().Done():
		glog.Infof("session.writeAuthChan sess.gRPCStream.Context().Done() sess:%#v", sess)
		return sess.gRPCStream.Context().Err()
	}
	return nil
}

func (sess *session) Auth(req *gatewaypb.LoginRequest) (*gatewaypb.LoginResponse, error) {
	rows, err := db.Query("SELECT * FROM ACCOUNT WHERE username=?", req.Username)
	if err != nil {
		glog.Warning("session.Auth db.Query failed", sess, req, err)
		return &gatewaypb.LoginResponse{gatewaypb.LoginResponse_MYSQL_MAINTAIN}, nil
	}
	defer rows.Close()

	account := authpb.Account{}
	for rows.Next() {
		if err := rows.Scan(&account.Username, &account.PasswdHash, &account.Email, &account.Gm); err != nil {
			glog.Warning("session.Auth rows.Scan failed", sess, req, err)
			return &gatewaypb.LoginResponse{gatewaypb.LoginResponse_INCORRECT_ACCOUNT_PASSWD}, nil
		}
	}

	if err := rows.Err(); err != nil {
		glog.Warning("session.Auth rows.Err", sess, req, err)
		return &gatewaypb.LoginResponse{gatewaypb.LoginResponse_INCORRECT_ACCOUNT_PASSWD}, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.PasswdHash), []byte(req.Password))
	if err != nil {
		glog.Warning("session.Auth CompareHashAndPassword failed", sess, req, err, account.PasswdHash, req.Password)
		return &gatewaypb.LoginResponse{gatewaypb.LoginResponse_INCORRECT_ACCOUNT_PASSWD}, nil
	}

	glog.Infof("session.Auth CompareHashAndPassword SUCCESS authServer:%#v req:%#v PasswdHash:%v Password:%v", sess, req, account.PasswdHash, req.Password)
	return &gatewaypb.LoginResponse{gatewaypb.LoginResponse_SUCCESS}, nil
}
