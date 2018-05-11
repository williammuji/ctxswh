package main

import (
	"flag"
	"net"
	"net/http"

	"github.com/golang/glog"

	"golang.org/x/net/trace"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"

	"github.com/williammuji/ctxswh/gatewaypb"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	var (
		port            = flag.String("port", "8001", "gateway listen port")
		debugPort       = flag.String("debugPort", "8002", "gateway debug http listen port")
		serverCert      = flag.String("cert", "/home/.auth/gateway.pem", "gateway server tls certificate")
		serverKey       = flag.String("key", "/home/.auth/gateway-key.pem", "gateway server tls key")
		authServerAddr  = flag.String("authServerAddr", "authservice:8101", "auth server addr")
		gameServerAddr  = flag.String("gameServerAddr", "gameservice:8201", "game server addr")
		chatServerAddr  = flag.String("chatServerAddr", "chatservice:8301", "chat server addr")
		loginServerAddr = flag.String("loginServerAddr", "loginservice:8401", "login server addr")
	)
	flag.Parse()
	defer glog.Flush()

	glog.Info("Gateway service starting...")

	creds, err := credentials.NewServerTLSFromFile(*serverCert, *serverKey)
	if err != nil {
		glog.Fatal(err)
	}
	s := grpc.NewServer(grpc.Creds(creds))

	gs, err := NewGatewayServer(*authServerAddr, *gameServerAddr, *chatServerAddr, *loginServerAddr)
	if err != nil {
		glog.Fatal(err)
	}
	defer gs.CloseConn()
	gatewaypb.RegisterGatewayServiceServer(s, gs)

	hs := health.NewServer()
	hs.SetServingStatus("grpc.health.v1.gatewayservice", 1)
	healthpb.RegisterHealthServer(s, hs)

	ln, err := net.Listen("tcp", ":"+(*port))
	if err != nil {
		glog.Fatal(err)
	}
	go s.Serve(ln)

	trace.AuthRequest = func(req *http.Request) (any, sensitive bool) { return true, true }
	glog.Info("Gateway service started successfully.")
	glog.Fatal(http.ListenAndServe(":"+(*debugPort), nil))
}
