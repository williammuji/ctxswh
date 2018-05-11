package main

import (
	"flag"
	"net"
	"net/http"

	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"golang.org/x/net/trace"

	"github.com/golang/glog"
	"github.com/williammuji/ctxswh/chatpb"
	"github.com/williammuji/ctxswh/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
)

func main() {
	var (
		port      = flag.String("port", "8301", "Chat service port.")
		debugPort = flag.String("debugPort", "8302", "Chat service debug port.")
	)
	flag.Parse()
	defer glog.Flush()

	glog.Info("Chat service starting...")

	pool := util.NewSentinelPool()
	defer pool.Close()

	s := grpc.NewServer()
	cs, err := NewChatServer(pool)
	if err != nil {
		glog.Fatal(err)
	}
	chatpb.RegisterChatServiceServer(s, cs)

	hs := health.NewServer()
	hs.SetServingStatus("grpc.health.v1.chatservice", 1)
	healthpb.RegisterHealthServer(s, hs)

	ln, err := net.Listen("tcp", ":"+(*port))
	if err != nil {
		glog.Fatal(err)
	}
	go s.Serve(ln)

	trace.AuthRequest = func(req *http.Request) (any, sensitive bool) { return true, true }
	glog.Info("Chat service started successfully.")
	glog.Fatal(http.ListenAndServe(":"+(*debugPort), nil))
}
