package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var version = "0.0.1"

func main() {
	log.Println("Starting frame demo...")

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	httpAddr := os.Getenv("HTTP_ADDR")
	if httpAddr == "" {
		httpAddr = fmt.Sprintf("%s:80", hostname)
		log.Printf("using default host: %s ", httpAddr)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, html, time.Now())
	})

	log.Printf("HTTP service listening on %s", httpAddr)
	http.ListenAndServe(httpAddr, nil)

}
