#!/usr/bin/env bash
set -x

: ${1?'missing key directory'}

key_dir="$1"

chmod 0700 "$key_dir"
cd "$key_dir"

cat >server.conf <<EOF
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name
prompt = no
[req_distinguished_name]
CN = webhook-server.webhook-for-jenkins.svc
[ v3_req ]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage = clientAuth, serverAuth
subjectAltName = @alt_names
[alt_names]
DNS.1 = webhook-server.webhook-for-jenkins.svc
EOF

echo here / //
# Generate the CA cert and private key
openssl req -nodes -new -x509 -keyout ca.key -out $key_dir/ca.crt -subj "//CN=Admission Controller Webhook for Jenkins CA"

echo now here
# Generate the private key for the webhook server
openssl genrsa -out $key_dir/webhook-server-tls.key 2048

ls -latr
pwd
echo now down here
# Generate a Certificate Signing Request (CSR) for the private key, and sign it with the private key of the CA.
openssl req -new -key $key_dir/webhook-server-tls.key -subj "/CN=webhook-server.webhook-for-jenkins.svc" -config $key_dir/server.conf \
    | openssl x509 -req -CA $key_dir/ca.crt -CAkey $key_dir/ca.key -CAcreateserial -out $key_dir/webhook-server-tls.crt -extensions v3_req -extfile $key_dir/server.conf
