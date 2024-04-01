package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

var (
	CertFilePath = "certs/server.pem"
	KeyFilePath  = "certs/server.key"
)

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		log.Printf("Error during NewRequest() %s: %s\n", r.URL.String(), err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// copy headers
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error during Do() %s: %s\n", r.URL.String(), err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	written, err := io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Error during Copy() %s: %s\n", r.URL.String(), err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("%s - %s - %s - %d - %dKB\n", r.Proto, r.Method, r.Host, resp.StatusCode, written/1000)
}

func handleTunnel(w http.ResponseWriter, r *http.Request) {
	dest_conn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)

	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
		return
	}
	src_conn, _, err := hj.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	go transfer(src_conn.LocalAddr().String(), src_conn, dest_conn.LocalAddr().String(), dest_conn)
	go transfer(dest_conn.LocalAddr().String(), dest_conn, src_conn.LocalAddr().String(), src_conn)

	log.Printf("%s - %s - %s\n", r.Proto, r.Method, r.Host)
}

func transfer(destName string, destination io.WriteCloser, sourceName string, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	written, err := io.Copy(destination, source)
	if err != nil {
		fmt.Printf("Error during copy: %v\n", err)
	}
	log.Printf("copied %d bytes %s -> %s\n", written, sourceName, destName)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("health check")
	fmt.Fprintf(w, "OK")
}

func main() {
	log.Fatal(http.ListenAndServeTLS(
		":8443",
		CertFilePath,
		KeyFilePath,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				handleTunnel(w, r)
			} else {
				if r.URL.Path == "/health" {
					healthCheck(w, r)
				} else {
					handleHTTP(w, r)
				}
			}
		})),
	)
}
