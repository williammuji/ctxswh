package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	var (
		gatewayServerAddr = flag.String("gatewayServerAddr", "gatewayservice:8001", "gateway server addr")
		gatewayCaCert     = flag.String("gatewayCaCert", "/etc/auth/mysql/ca.pem", "gateway server ca certificate")
		gatewayClientCert = flag.String("gatewayClientCert", "/home/.auth/client-cert.pem", "gateway client certificate.")
		gatewayClientKey  = flag.String("gatewayClientKey", "/home/.auth/client-key.pem", "gateway client key.")
		user              = flag.String("user", "testuser1:testuser2", "user names.")
		passwd            = flag.String("passwd", "testpasswd1:testpasswd2", "user passwds.")
		secs              = flag.Int("secs", 30, "secs after client stop.")
		ppMsgSize         = flag.Int("ppMsgSize", 32, "pingpong msg size.")
	)
	flag.Parse()
	defer glog.Flush()

	glog.Info("Client service starting...")

	cert, err := tls.LoadX509KeyPair(*gatewayClientCert, *gatewayClientKey)
	if err != nil {
		glog.Fatal(err)
	}

	rawCACert, err := ioutil.ReadFile(*gatewayCaCert)
	if err != nil {
		glog.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(rawCACert)

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	})

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(creds))

	glog.Info("Client service started successfully.")
	users := strings.Split(*user, ":")
	passwds := strings.Split(*passwd, ":")
	if len(users) != len(passwds) {
		glog.Fatalf("users:%#v length not equal passwds:%#v length", users, passwds)
	}

	rand.Seed(time.Now().UnixNano())

	var wg sync.WaitGroup
	wg.Add(len(users))
	for i, userName := range users {
		sess := NewSession(*gatewayServerAddr, opts, userName, passwds[i], *secs, *ppMsgSize, &wg)
		go sess.start()
	}
	wg.Wait()
	glog.Info("Client service stopped.")
}
