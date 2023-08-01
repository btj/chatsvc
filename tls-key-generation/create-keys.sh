#!/bin/sh
set -e -x
openssl genrsa -aes256 -out ../../rootCA.key 2048
openssl req -new -key ../../rootCA.key -out ../../rootCA.csr -subj '/CN=chatsvc.mooo.com trust CA'
openssl x509 -req -sha256 -days 1024 -in ../../rootCA.csr -signkey ../../rootCA.key -extfile rootCAext.ini -out ../../rootCA.pem

openssl req -new -sha256 -nodes -out ../../server.csr -newkey rsa:2048 -keyout ../../server.key -config server.csr.cnf
openssl x509 -req -in ../../server.csr -CA ../../rootCA.pem -CAkey ../../rootCA.key -CAcreateserial -out ../../server.crt -days 500 -sha256 -extfile v3.ext
cat ../../server.crt ../../rootCA.pem > ../../server_and_ca.crt
