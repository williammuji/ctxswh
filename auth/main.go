package main

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"flag"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/golang/glog"

	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/go-sql-driver/mysql"
	"github.com/williammuji/ctxswh/authpb"
	"golang.org/x/net/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
)

var db *sql.DB

func main() {
	var (
		dbUser         = flag.String("dbUser", "testuser", "Mysql username.")
		dbPasswd       = flag.String("dbPasswd", "testpasswd", "Mysql password.")
		dbHost         = flag.String("dbHost", "MySQL_Server_5.7.16_Auto_Generated_Server_Certificate", "Mysql host.")
		dbPort         = flag.String("dbPort", "3306", "Mysql port.")
		dbDatabaseName = flag.String("dbDatabaseName", "testdb", "Mysql database name.")
		dbServerCACert = flag.String("dbServerCACert", "/etc/auth/mysql/ca.pem", "Mysql server ca certificate")
		dbClientCert   = flag.String("dbClientCert", "/etc/auth/mysql/client-cert.pem", "Mysql client certificate.")
		dbClientKey    = flag.String("dbClientKey", "/etc/auth/mysql/client-key.pem", "Mysql client key.")
		port           = flag.String("port", "8101", "Auth service port.")
		debugPort      = flag.String("debugPort", "8102", "Auth service debug port.")
	)
	flag.Parse()
	defer glog.Flush()

	var err error
	glog.Info("Auth service starting...")

	certPool := x509.NewCertPool()
	pem, err := ioutil.ReadFile(*dbServerCACert)
	if err != nil {
		glog.Fatal(err)
	}
	if ok := certPool.AppendCertsFromPEM(pem); !ok {
		glog.Fatal("Failed to append PEM.")
	}
	clientCert := make([]tls.Certificate, 0, 1)
	certs, err := tls.LoadX509KeyPair(*dbClientCert, *dbClientKey)
	if err != nil {
		glog.Fatal(err)
	}
	clientCert = append(clientCert, certs)
	mysql.RegisterTLSConfig("custom", &tls.Config{
		RootCAs:      certPool,
		Certificates: clientCert,
	})

	dbAddr := net.JoinHostPort(*dbHost, *dbPort)
	dbConfig := mysql.Config{
		User:      *dbUser,
		Passwd:    *dbPasswd,
		Net:       "tcp",
		Addr:      dbAddr,
		DBName:    *dbDatabaseName,
		TLSConfig: "custom",
	}

	for {
		db, err = sql.Open("mysql", dbConfig.FormatDSN())
		if err != nil {
			goto dberror
		}
		err = db.Ping()
		if err != nil {
			goto dberror
		}
		break

	dberror:
		glog.Warning(err)
		glog.Warning("error connecting to the auth database, retrying in 5 secs.")
		time.Sleep(5 * time.Second)
	}

	s := grpc.NewServer()

	as, err := NewAuthServer()
	if err != nil {
		glog.Fatal(err)
	}
	authpb.RegisterAuthServiceServer(s, as)

	hs := health.NewServer()
	hs.SetServingStatus("grpc.health.v1.authservice", 1)
	healthpb.RegisterHealthServer(s, hs)

	ln, err := net.Listen("tcp", ":"+(*port))
	if err != nil {
		glog.Fatal(err)
	}
	go s.Serve(ln)

	trace.AuthRequest = func(req *http.Request) (any, sensitive bool) { return true, true }
	glog.Info("Auth service started successfully.")
	glog.Fatal(http.ListenAndServe(":"+(*debugPort), nil))
}
