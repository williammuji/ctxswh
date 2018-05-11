# Command For Add Account

## Add mysql user and database

```
$ mysql -u root -prootpasswd
mysql> create database testdb;
mysql> use testdb;
mysql> CREATE TABLE ACCOUNT (username varchar(255) NOT NULL, passwdHash varchar(255) NOT NULL, email varchar(255) NOT NULL, gm BOOL, PRIMARY KEY (username));
mysql> create user 'testuser'@'%' identified by 'testpasswd' REQUIRE SSL;
mysql> grant all on testdb.* to 'testuser';
```

## User 'testuser' login mysql with SSL

```
# mysql -u testuser -p --ssl-ca=${CTXSWH_PATH}/mysqltls/ca.pem --ssl-cert=${CTXSWH_PATH}/mysqltls/server.pem --ssl-key=${CTXSWH_PATH}/mysqltls/server-key.pem
```

## Run pod account

```
# kubectl create -f account/accountcontroller.yaml
# kubectl get pods
# kubectl exec -it account-ypwjk -- /account -dbHost=192.168.98.170  -dbServerCACert=/mysqltls/ca.pem -dbClientCert=/mysqltls/client.pem -dbClientKey=/mysqltls/client-key.pem -u=testuser1 -e=testuser1@testuser1.com
```
