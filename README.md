### TCP and UDP proxy in Go

This version is heavily modified from the original source by [arkadijs](https://github.com/arkadijs/goproxy).

- Proxies inbound port to target:port
- More than one route can be defined
- UDP is supported (and is bi-directional)

The added UDP code is ugly as it's a bit of a mish-mash from other proxies.

Usage:

    $ goproxy [flags] listen-port:upstream-hostname:upstream-port listen-port2:upstream-hostname2:upstream-port2 ...
      -debug=false: Print every connection information
      -dns="": DNS server address, supply host:port; will use system default if not set
      -dns-interval=20s: Time interval between DNS queries
      -timeout=10s: TCP connect timeout
      -udp=false: UDP mode, will start listeners on BOTH TCP and UDP
      -verbose=false: Print noticeable info

