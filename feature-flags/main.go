package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var version = "0.0.1"

func main() {
	log.Println("Starting kv-sr demo...")

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	httpAddr := os.Getenv("HTTP_ADDR")
	if httpAddr == "" {
		httpAddr = fmt.Sprintf("%s:80", hostname)
		log.Printf("using default host: %s ", httpAddr)
	}

	log.Printf("Registering with consul as %s", httpAddr)
	registerService("demo", httpAddr)

	// Some channel madness for watching KV stuff

	// Set up some signal handling here
	sigs := make(chan os.Signal, 1)
	done := make(chan bool)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for _ = range sigs {
			fmt.Println("\nReceived an interrupt, deregistering services...\n")
			deRegisterService("demo", httpAddr)
			os.Exit(0)
			done <- true
		}
	}()

	//magic := "who knows"
	//var index uint64 = 0
	magic, index := watchKey("/flags/magic", 0)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//magic := getKey("/flags/magic")
		//log.Printf("HTTP trying to serve %s", magic)
		m2 := magic
		fmt.Fprintf(w, html, m2, hostname, version)
	})

	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		magic, index = watchKey("/flags/magic", 0)
		log.Printf("updating magic to %s at index %d", magic, index)
		fmt.Fprintf(w, "<html>cool story bro</html>")
	})

	log.Printf("HTTP service listening on %s", httpAddr)
	http.ListenAndServe(httpAddr, nil)
	<-done

}
