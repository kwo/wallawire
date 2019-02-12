#!/bin/bash

mkdir -p certs
CERTDIR=./certs

echo '{"signing":{"default":{"expiry":"168h"},"profiles":{"server":{"expiry":"26280h","usages":["signing","key encipherment","server auth"]},"client":{"expiry":"26280h","usages":["signing","key encipherment","client auth"]}}}}' > ca-config.json

# CA
echo '{"CN":"Wallawire UI CA","key":{"algo":"ecdsa","size":256}}' | cfssl gencert -initca - | cfssljson -bare ca -
mv ca.pem ca.crt
mv ca-key.pem ca.key
rm ca.csr

# Web Server
echo '{"CN":"'wallawire'","hosts":[""],"key":{"algo":"ecdsa","size":256}}' | cfssl gencert -config=ca-config.json -profile=server -ca=ca.crt -ca-key=ca.key -hostname="localhost,127.0.0.1" - | cfssljson -bare server
mv server.pem ${CERTDIR}/server.crt
mv server-key.pem ${CERTDIR}/server.key
rm server.csr

mv ca.crt ${CERTDIR}/
mv ca.key ${CERTDIR}/
rm ca-config.json
