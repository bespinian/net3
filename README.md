# net³

net³ [netkube] is a tool to debug and understand network traffic in Kubernetes.

## Usage

### Topo

Show the network topology of a specific connection

```shell
$ net3 topo -n source mypod123 myservice123

┌─────────────────────────────────────────────────────┐
│ Pod                                                 │
│ Name: mypod123                                      │
│ Namespace: my-ns                                    │
│ Status: Running                                     │
│ Ports: TCP:8080 (http), TCP:15090 (http-envoy-prom) │
└─────────────────────────────────────────────────────┘
      │
      │
      │ TCP:80
      │
      V
┌────────────────────────┐
│ Service                │
│ Name: myservice123     │
│ Namespace: other-ns    │
│ Ports: TCP:80 (http)   │
│ Target Ports: TCP:http │
└────────────────────────┘
      │
      │
      │ http
      │
      V
┌─────────────────────────────────────────────────┐
│ Ingress Network Policy                          │
│ Name: allow-ingress                             │
│ Namespace: other-ns                             │
│ Rule: Allow all traffic from all pods           │
└─────────────────────────────────────────────────┘
      │
      │
      │ http
      │
      V
┌─────────────────────────────────────────────────────┐
│ Pod                                                 │
│ Name: myotherpod123                                 │
│ Namespace: other-ns                                 │
│ Status: Running                                     │
│ Ports: TCP:8080 (http), TCP:15090 (http-envoy-prom) │
└─────────────────────────────────────────────────────┘

```

### Proxy

Add a logging proxy to an existing service. Currently, only the HTTP protocol is supported.

To add a proxy, run

```shell
$ net3 proxy add -n mynamespace123 myservice123 80
```

and in another shell run

```shell
$ kubectl logs -n mynamespace123 mypod123 -c net3-proxy -f
```

When you're done, remove the logging proxy from a service by running

```shell
$ net3 proxy remove -n mynamespace123 myservice123 80
```
