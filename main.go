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

func proxy(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v\n", r.URL.String())

	req, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		log.Printf("Error during NewRequest() %s: %s\n", r.URL.String(), err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("%v\n", req)

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

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Error during Copy() %s: %s\n", r.URL.String(), err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("%d - %v\n", resp.StatusCode, r.URL.String())
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("health check")
	fmt.Fprintf(w, "OK")
}

func main() {
	http.HandleFunc("/", proxy)
	http.HandleFunc("/health", healthCheck)
	log.Fatal(http.ListenAndServeTLS(":8443", CertFilePath, KeyFilePath, nil))
}
