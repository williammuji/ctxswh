package main

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"time"

	"github.com/golang/glog"

	"github.com/go-sql-driver/mysql"
	"github.com/williammuji/ctxswh/authpb"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
)

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
		username       = flag.String("u", "", "New account username.")
		email          = flag.String("e", "", "New account email address.")
		gm             = flag.Bool("gm", false, "New account game master flag.")
	)
	flag.Parse()
	defer glog.Flush()

	fmt.Println("enter password:")
	password, err := terminal.ReadPassword(0)
	if err != nil {
		glog.Fatal(err)
	}

	passwdHash, err := bcrypt.GenerateFromPassword(password, 16)
	if err != nil {
		glog.Fatal(err)
	}

	account := authpb.Account{
		Username:   *username,
		PasswdHash: string(passwdHash),
		Email:      *email,
		Gm:         *gm,
	}

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

	var db *sql.DB
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

	_, err = db.Exec("INSERT INTO ACCOUNT (username, passwdHash, email, gm) VALUES (?, ?, ?, ?)", account.Username, account.PasswdHash, account.Email, account.Gm)
	if err != nil {
		glog.Fatal(err)
	}
}
