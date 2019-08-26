# Authenticate cluster users through authproxy

## Deploy authproxy on Kubernetes

It is recommended to run authproxy as DaemonSet on the master nodes in the cluster. The corresponding deployment files can be found in the projects for the specific provider implementation.

## Kubernetes setup

Next, configure kube apiserver to verify bearer token using this authenticator.
There are two configuration options you need to set:

* `--authentication-token-webhook-config-file` a kubeconfig file describing how to
  access the remote webhook service (this must be the address of your authproxy service running inside or outside of your cluster).
* `--authentication-token-webhook-cache-ttl` how long to cache authentication
  decisions. Defaults to two minutes.

Check the example config file below and save this file on the Kubernetes master node(s). In case you need to provide additional tls certs, save them under `/etc/kubernetes/pki/certs`.

```
# Kubernetes API version
apiVersion: v1
# kind of the API object
kind: Config
# clusters refers to the remote service.
clusters:
  - name: name-of-remote-authz-service
    cluster:
      # CA for verifying the remote service.
      certificate-authority: /path/to/ca.pem
      # URL of remote service to query. Must use 'https'. May not include parameters.
      server: https://authz.example.com/authorize

# users refers to the API Server's webhook configuration.
users:
  - name: name-of-api-server
    user:
      client-certificate: /path/to/cert.pem # cert for the webhook plugin to use
      client-key: /path/to/key.pem          # key matching the cert

# kubeconfig files require a context. Provide one for the API Server.
current-context: webhook
contexts:
- context:
    cluster: name-of-remote-authz-service
    user: name-of-api-server
  name: webhook
```

It is recommended you read the [Kubernetes
documentation](https://kubernetes.io/docs/admin/authentication/#webhook-token-authentication) for how to configure
webhook token authentication.
