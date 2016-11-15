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
var consulAddr = "localhost:8500"

func signalHandler(sigs chan os.Signal, done chan bool, svcAddr string) {
	for _ = range sigs {
		fmt.Println("\nReceived an interrupt, deregistering services...\n")
		deRegisterService("demo", svcAddr)
		os.Exit(0)
		done <- true
	}
}

func main() {
	log.Println("Starting kv-sr demo...")

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	consulEnv := os.Getenv("CONSUL_ADDR")
	if consulEnv != "" {
		consulAddr = consulEnv
		log.Printf("using environment for consul address: %s ", consulAddr)
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "80"
		log.Printf("using default port: %s ", httpPort)
	}

	svcAddr := os.Getenv("HTTP_ADDR")
	if svcAddr == "" {
		svcAddr = fmt.Sprintf("%s", hostname)
		log.Printf("using default host: %s ", svcAddr)
	}

	svcAddr = fmt.Sprintf("%s:%s", svcAddr, httpPort)

	// we want to listen everywhere
	listenAddr := fmt.Sprintf("0.0.0.0:%s", httpPort)

	log.Printf("Registering with consul as %s", svcAddr)
	registerService("demo", svcAddr)

	// Set up some signal handling here
	sigs := make(chan os.Signal, 1)
	signalHandlerDone := make(chan bool)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go signalHandler(sigs, signalHandlerDone, svcAddr)

	// Go find the initial flag value
	magic, index := watchKey("/flags/magic", 0)

	// Keep watching for updates to the magic flag
	// This is a blocking query, so it's cheap.
	go func() {
		for {
			magic, index = watchKey("/flags/magic", index)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// optimized out without this?
		mymagic := magic
		log.Printf("HTTP trying to serve %s and %d", mymagic, index)
		if magic != "true" {
			fmt.Fprintf(w, html, mymagic, hostname, version)
		} else {
			fmt.Fprintf(w, htmltrue)
		}
	})

	/*
		http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
			magic, index = watchKey("/flags/magic", 0)
			log.Printf("updating magic to %s at index %d", magic, index)
			fmt.Fprintf(w, "<html>cool story bro</html>")
		})
	*/

	log.Printf("HTTP service listening on %s", listenAddr)
	http.ListenAndServe(listenAddr, nil)

	// wait for cleanup, we deregister with this hook
	<-signalHandlerDone

}
