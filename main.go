package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
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

	log.Printf("%s - %d - %dKB\n", r.Host, resp.StatusCode, written/1000)
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
				log.Println("CONNECT not implemented")
				w.WriteHeader(http.StatusNotImplemented)
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
