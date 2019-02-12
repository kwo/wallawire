#!/bin/bash

mkdir -p walladata/certs/{dbserver,dbclient,webserver,private}
CERTDIR=walladata/certs

echo '{"signing":{"default":{"expiry":"168h"},"profiles":{"server":{"expiry":"26280h","usages":["signing","key encipherment","server auth"]},"client":{"expiry":"26280h","usages":["signing","key encipherment","client auth"]}}}}' > ca-config.json

# CA
echo '{"CN":"Wallawire CA","key":{"algo":"ecdsa","size":256}}' | cfssl gencert -initca - | cfssljson -bare ca -
mv ca.pem ca.crt
mv ca-key.pem ca.key
rm ca.csr

# Database Server
echo '{"CN":"'node'","hosts":[""],"key":{"algo":"ecdsa","size":256}}' | cfssl gencert -config=ca-config.json -profile=server -ca=ca.crt -ca-key=ca.key -hostname="localhost,db,127.0.0.1" - | cfssljson -bare node
mv node.pem ${CERTDIR}/dbserver/node.crt
mv node-key.pem ${CERTDIR}/dbserver/node.key
rm node.csr

# Database Client root
echo '{"CN":"'root'","hosts":[""],"key":{"algo":"ecdsa","size":256}}' | cfssl gencert -config=ca-config.json -profile=client -ca=ca.crt -ca-key=ca.key - | cfssljson -bare root
mv root.pem ${CERTDIR}/dbclient/client.root.crt
mv root-key.pem ${CERTDIR}/dbclient/client.root.key
rm root.csr

# Database Client wallawire
echo '{"CN":"'wallawire'","hosts":[""],"key":{"algo":"ecdsa","size":256}}' | cfssl gencert -config=ca-config.json -profile=client -ca=ca.crt -ca-key=ca.key - | cfssljson -bare wallawire
mv wallawire.pem ${CERTDIR}/dbclient/client.wallawire.crt
mv wallawire-key.pem ${CERTDIR}/dbclient/client.wallawire.key
rm wallawire.csr

# Web Server
echo '{"CN":"'wallawire'","hosts":[""],"key":{"algo":"ecdsa","size":256}}' | cfssl gencert -config=ca-config.json -profile=server -ca=ca.crt -ca-key=ca.key -hostname="localhost,127.0.0.1" - | cfssljson -bare server
mv server.pem ${CERTDIR}/webserver/server.crt
mv server-key.pem ${CERTDIR}/webserver/server.key
rm server.csr

cp ca.crt ${CERTDIR}/dbclient/
cp ca.crt ${CERTDIR}/dbserver/
cp ca.crt ${CERTDIR}/webserver/
mv ca.key ${CERTDIR}/private/
rm ca-config.json
mv ca.crt wallawire-ca.crt
