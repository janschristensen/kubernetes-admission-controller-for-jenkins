#!/usr/bin/env bash
kubectl delete namespace webhook-for-jenkins
kubectl delete MutatingWebhookConfiguration mutating-webhook-for-jenkins
bash -c 'export GOOS="linux" && echo go build ... && cd cmd\\webhook-server && go build && echo go build succes'
IMAGE=localhost:5000/dk-bankdata-web/admission-controller-webhook-for-jenkins:latest
docker build . -t $IMAGE -f ./image/Dockerfile
docker push $IMAGE

set -euo pipefail

basedir="$(dirname "$0")/deployment"
keydir="$(mktemp -d)"

# Generate keys into a temporary directory.
echo "Generating TLS keys ..."
"${basedir}/generate-keys.sh" "$keydir"

# Create the `webhook-for-jenkins` namespace. This cannot be part of the YAML file as we first need to create the TLS secret,
# which would fail otherwise.
echo "Creating Kubernetes objects ..."
kubectl create namespace webhook-for-jenkins

# Create the TLS secret for the generated keys.
kubectl -n webhook-for-jenkins create secret tls webhook-server-tls \
    --cert "${keydir}/webhook-server-tls.crt" \
    --key "${keydir}/webhook-server-tls.key"

# Read the PEM-encoded CA certificate, base64 encode it, and replace the `${CA_PEM_B64}` placeholder in the YAML
# template with it. Then, create the Kubernetes resources.
ca_pem_b64="$(openssl base64 -A <"${keydir}/ca.crt")"
sed -e 's@${CA_PEM_B64}@'"$ca_pem_b64"'@g' <"${basedir}/deployment.yaml.template" \
    | sed 's@${registry}@localhost:5000@' |   
    kubectl create -f -

# Delete the key directory to prevent abuse (DO NOT USE THESE KEYS ANYWHERE ELSE).
rm -rf "$keydir"

echo "The webhook server has been deployed and configured!"
