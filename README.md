# Kubernetes Admission Controller for jenkins

This repository contains a small HTTP server that can be used as a Kubernetes
[MutatingAdmissionWebhook](https://kubernetes.io/docs/admin/admission-controllers/#mutatingadmissionwebhook-beta-in-19).

The hook is intended to stop pods ASAP when a wrong registry is provided and notify jenkins about which containers needs to be fixed.

This is needed as we often see developers forget to specify the JNLP container, and in those cases the kubernetes plugin in Jenkins will default to a image in docker hub. Which often is ok, but in our case not as we are using a private registry and disallowing others.

Ideally in next version we will start mutating the container image, to include a valid registry.

## Prerequisites

A cluster on which this can be tested must be running Kubernetes 1.9.0 or above,
with the `admissionregistration.k8s.io/v1beta1` API enabled. You can verify that by observing that the
following command produces a non-empty output:
```
kubectl api-versions | grep admissionregistration.k8s.io/v1beta1
```
In addition, the `MutatingAdmissionWebhook` admission controller should be added and listed in the admission-control
flag of `kube-apiserver`.

For building the image [Go](https://golang.org) are required.

## Deploying the Webhook Server

1. Bring up a Kubernetes cluster satisfying the above prerequisites, and make
sure it is active (i.e., either via the configuration in the default location, or by setting
the `KUBECONFIG` environment variable).

In dev run
```
kind/bootcluster.sh
```

2. Run `./deploy.sh`. This will create a CA, a certificate and private key for the webhook server,
and deploy the resources in the newly created `webhook-for-jenkins` namespace in your Kubernetes cluster.

In dev run 
```
deploy-dev.sh
```

## Verify

1. The `webhook-server` pod in the `webhook-for-jenkins` namespace should be running:
```
$ kubectl -n webhook-for-jenkins get pods
NAME                             READY     STATUS    RESTARTS   AGE
webhook-server-6f976f7bf-hssc9   1/1       Running   0          35m
```

2. A `MutatingWebhookConfiguration` named `mutating-webhook-for-jenkins` should exist:
```
$ kubectl get mutatingwebhookconfigurations
```

3. Deploy [a pod](examples/dev/pod-with-error.yaml) that contains an error:
```
$ kubectl create -f examples/dev/pod-with-error.yaml -n webhook-for-jenkins
```
Should yield
```
Error from server: error when creating "examples/dev/pod-with-error.yaml": admission webhook "webhook-server.webhook-for-jenkins.svc" denied the request: Error in image, you need to fix remote-busybox:busybox
```

4. Deploy [a pod](examples/dev/pod-working.yaml) that work. 

```
$ kubectl create -f examples/dev/pod-working.yaml
```

Should yield
```
pod/pod-working created
```

## Build the Image from Sources locally

cross compile from windows to linux 

```
export GOOS="linux"
cd cmd\webhook-server\
go build
```

## cleanup cluster

```
kind delete cluster --name kind
```

## cleanup namespacec and webhook

```
kubectl delete namespace webhook-for-jenkins
kubectl delete MutatingWebhookConfiguration mutating-webhook-for-jenkins
```

## Helpful LINKS
Based on code from this [repo](https://github.com/stackrox/admission-controller-webhook-demo) 
* https://medium.com/ovni/writing-a-very-basic-kubernetes-mutating-admission-webhook-398dbbcb63ec 
* [youtube.com/watch?v=r_v07P8Go6w](https://www.youtube.com/watch?v=r_v07P8Go6w)
* https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/
* https://kind.sigs.k8s.io/docs/user/local-registry/
* https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/

Consider fixing retryBackoff - https://stackoverflow.com/questions/57417027/kubernetes-jobs-and-back-off-limit-values-is-the-value-a-number-of-retries-or-m