# net³

net³ [netkube] is a tool to debug and understand network traffic in Kubernetes.

## Usage

### Topo

Show the network topology of a specific connection

```shell
$ net3 topo mypod123 myservice123

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

## Roadmap

- [ ] Kubernetes support for `topo` command
- [ ] Istio support for `topo` command
