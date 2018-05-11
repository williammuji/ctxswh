package main

import (
	"flag"
	"net"
	"net/http"

	"github.com/golang/glog"

	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/williammuji/ctxswh/gamepb"

	"golang.org/x/net/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
)

func main() {
	var (
		port      = flag.String("port", "8201", "Auth service port.")
		debugPort = flag.String("debugPort", "8202", "Auth service debug port.")
	)
	flag.Parse()
	defer glog.Flush()

	var err error
	glog.Info("Game service starting...")

	s := grpc.NewServer()

	as, err := NewGameServer()
	if err != nil {
		glog.Fatal(err)
	}
	gamepb.RegisterGameServiceServer(s, as)

	hs := health.NewServer()
	hs.SetServingStatus("grpc.health.v1.gameservice", 1)
	healthpb.RegisterHealthServer(s, hs)

	ln, err := net.Listen("tcp", ":"+(*port))
	if err != nil {
		glog.Fatal(err)
	}
	go s.Serve(ln)

	trace.AuthRequest = func(req *http.Request) (any, sensitive bool) { return true, true }
	glog.Info("Game service started successfully.")
	glog.Fatal(http.ListenAndServe(":"+(*debugPort), nil))
}
