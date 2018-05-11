package main

import (
	"flag"
	"net"
	"net/http"

	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"

	"github.com/golang/glog"
	"github.com/williammuji/ctxswh/loginpb"

	"golang.org/x/net/trace"
)

var (
	ls *loginServer
)

func MinLoadIPHandler(w http.ResponseWriter, req *http.Request) {
	ip := ls.getMinCountIp()
	w.Write([]byte(ip))

	glog.Infof("MinLoadIPHandler w:%#v req:%#v ip:%s", w, req, ip)
}

func main() {
	var (
		port      = flag.String("port", "7001", "Login service inner port.")
		httpPort  = flag.String("httpPort", "7002", "Login http service port.")
		debugPort = flag.String("debugPort", "7003", "Login service debug port.")
	)
	flag.Parse()
	defer glog.Flush()

	glog.Info("Login service starting...")

	var err error
	s := grpc.NewServer()
	ls, err = NewLoginServer()
	if err != nil {
		glog.Fatal(err)
	}
	loginpb.RegisterLoginServiceServer(s, ls)

	hs := health.NewServer()
	hs.SetServingStatus("grpc.health.v1.loginservice", 1)
	healthpb.RegisterHealthServer(s, hs)

	ln, err := net.Listen("tcp", ":"+(*port))
	if err != nil {
		glog.Fatal(err)
	}
	go s.Serve(ln)

	serverMuxOut := http.NewServeMux()
	serverMuxOut.HandleFunc("/lb", MinLoadIPHandler)
	go func() {
		http.ListenAndServe(":"+(*httpPort), serverMuxOut)
	}()

	trace.AuthRequest = func(req *http.Request) (any, sensitive bool) { return true, true }
	glog.Info("Login service started successfully.")
	glog.Fatal(http.ListenAndServe(":"+(*debugPort), nil))
}
