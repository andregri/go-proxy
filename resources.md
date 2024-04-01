https://corte.si/posts/code/mitmproxy/howitworks/index.html

When the client makes a https request, the proxy receives a CONNECT http request that means to open a bi-directional tunnel with the target server.
Since the request is protected by TLS, the proxy server can't read the actual data of the HTTPS request.

To make a proper man-in-the-middle proxy server, the proxy server should generate certificates on-the-fly.

https://www.alexedwards.net/blog/an-introduction-to-handlers-and-servemuxes-in-go

Default HTTP handler in Go is the [ServeMux](https://cs.opensource.google/go/go/+/refs/tags/go1.22.1:src/net/http/server.go;l=2432) that matches the url against a list of patterns.
In case of CONNECT requests, it redirects the connection to the same host and path specified in the request.
The CONNECT request to google.com/ points to the "google.com" host and the "/" root path.
Since my handler was handling just the root path "/", the handler didn't match the host.