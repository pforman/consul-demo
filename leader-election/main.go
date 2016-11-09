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

func signalHandler(sigs chan os.Signal, done chan bool, httpAddr string, sid string) {
	for _ = range sigs {
		fmt.Println("\nReceived an interrupt, deregistering services...\n")
		releaseLock(sid, "service/demo/leader")
		deRegisterService("demo", httpAddr)
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

	httpAddr := os.Getenv("HTTP_ADDR")
	if httpAddr == "" {
		httpAddr = fmt.Sprintf("%s:80", hostname)
		log.Printf("using default host: %s ", httpAddr)
	}

	log.Printf("Registering with consul as %s", httpAddr)
	registerService("demo", httpAddr)
	sid := consulSession(httpAddr, "60s")

	// Set up some signal handling here
	sigs := make(chan os.Signal, 1)
	signalHandlerDone := make(chan bool)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go signalHandler(sigs, signalHandlerDone, httpAddr, sid)

	// Go find the current index value
	_, index := watchKey("service/demo/leader", 0, 0)

	isLeader := acquireLock(sid, "service/demo/leader", httpAddr)
	log.Printf("sid is %s, leader is %b", sid, isLeader)

	// Keep watching for updates to the lock
	// This is a blocking query, so it's cheapish.
	go func() {
		for {
			// TODO: need to put bounds in here to wake up
			// and renew the session
			_, index = watchKey("service/demo/leader", index, 30)
			log.Printf("woke up, checking the lock")
			if !isLeader {
				log.Printf(" not leader, trying to acquire the lock")
				isLeader = acquireLock(sid, "service/demo/leader", httpAddr)
				if isLeader {
					log.Printf("acquired lock.  now what?")
				}

			}
			sid = renewSession(sid)
			log.Printf("sid is %s, index is %d, leader is %t", sid, index, isLeader)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// optimized out without this?
		leading := isLeader
		log.Printf("HTTP trying to serve as leader=%t", leading)
		fmt.Fprintf(w, html, leading, hostname, version)
	})

	http.HandleFunc("/release", func(w http.ResponseWriter, r *http.Request) {
		if releaseLock(sid, "service/demo/leader") {
			log.Printf("releasing lock from sid %s", sid)
			isLeader = false
		} else {
			fmt.Printf("trying to release lock when we aren't leader, whoops")
		}
		fmt.Fprintf(w, "<html>cool story bro</html>")
	})

	log.Printf("HTTP service listening on %s", httpAddr)
	http.ListenAndServe(httpAddr, nil)

	// wait for cleanup, we deregister with this hook
	<-signalHandlerDone

}
