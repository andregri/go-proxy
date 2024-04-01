https://medium.com/@mlowicki/http-s-proxy-in-golang-in-less-than-100-lines-of-code-6a51c2f2c38c

https://corte.si/posts/code/mitmproxy/howitworks/index.html

When the client makes a https request, the proxy receives a CONNECT http request that means to open a bi-directional tunnel with the target server.
Since the request is protected by TLS, the proxy server can't read the actual data of the HTTPS request.

To make a proper man-in-the-middle proxy server, the proxy server should generate certificates on-the-fly.

https://www.alexedwards.net/blog/an-introduction-to-handlers-and-servemuxes-in-go

Default HTTP handler in Go is the [ServeMux](https://cs.opensource.google/go/go/+/refs/tags/go1.22.1:src/net/http/server.go;l=2432) that matches the url against a list of patterns.
In case of CONNECT requests, it redirects the connection to the same host and path specified in the request.
The CONNECT request to google.com/ points to the "google.com" host and the "/" root path.
Since my handler was handling just the root path "/", the handler didn't match the host.

https://stackoverflow.com/questions/75418196/correct-usage-of-io-copy-to-proxy-data-between-two-net-conn-tcp-connections-in-g
https://stackoverflow.com/questions/32460618/golang-1-5-io-copy-blocked-with-two-tcpconn?rq=4
https://stackoverflow.com/questions/75418196/correct-usage-of-io-copy-to-proxy-data-between-two-net-conn-tcp-connections-in-g
https://gist.github.com/jbardin/821d08cb64c01c84b81a
How to proxy 2 tcp connections?
Many discussions on who should close the connection but after some tests I guess it is important to close them once after all streams are copied.