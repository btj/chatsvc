#!/bin/sh
openssl x509 -noout -fingerprint -sha256 -inform pem -in ../../server.crt
