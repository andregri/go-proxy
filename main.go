package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func proxy(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v\n", r.URL.String())

	req, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		log.Printf("Error %s: %s\n", r.URL.String(), err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error during NewRequest(): %s", err)
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
		log.Printf("Error %s: %s\n", r.URL.String(), err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error during Do(): %s", err)
		return
	}
	defer resp.Body.Close()

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Error %s: %s\n", r.URL.String(), err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error during Copy(): %s", err)
		return
	}

	log.Printf("%d - %v\n", resp.StatusCode, r.URL.String())
}

func main() {
	http.HandleFunc("/", proxy)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
