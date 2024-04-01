#!/usr/bin/env bash

case $(uname -s) in
    Linux*)     sslConfig=/etc/ssl/openssl.cnf;;
    Darwin*)    sslConfig=/System/Library/OpenSSL/openssl.cnf;;
esac

mkdir -p certs

openssl req \
    -newkey rsa:2048 \
    -x509 \
    -nodes \
    -keyout certs/server.key \
    -new \
    -out certs/server.pem \
    -subj /CN=localhost \
    -reqexts SAN \
    -extensions SAN \
    -config <(cat "$sslConfig" \
        <(printf '[SAN]\nsubjectAltName=DNS:localhost')) \
    -sha256 \
    -days 365