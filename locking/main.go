package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var version = "0.0.1"

func signalHandler(sigs chan os.Signal, done chan bool, httpAddr string, sid string) {
	for _ = range sigs {
		fmt.Println("\nReceived an interrupt, deregistering services...\n")
		releaseLock(sid, "service/demo/worklock")
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
	sid := consulSession(httpAddr, "180s")

	// Set up some signal handling here
	sigs := make(chan os.Signal, 1)
	signalHandlerDone := make(chan bool)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go signalHandler(sigs, signalHandlerDone, httpAddr, sid)

	// Go find the current index value
	_, index := watchKey("service/demo/worklock", 0, 0)

	hasLock := acquireLock(sid, "service/demo/worklock", httpAddr)
	log.Printf("sid is %s, hasLock is %b", sid, hasLock)

	// Keep watching for updates to the lock
	// This is a blocking query, so it's cheapish.
	go func() {
		for {
			// TODO: need to put bounds in here to wake up
			// and renew the session
			_, index = watchKey("service/demo/worklock", index, 1)
			log.Printf("woke up, checking the lock")
			hasLock = acquireLock(sid, "service/demo/worklock", httpAddr)
			if !hasLock {
				log.Printf(" no lock, still trying to acquire")
			} else {
				log.Printf("I have the lock.  snoozing...")
				time.Sleep(6 * time.Second)
				releaseLock(sid, "service/demo/worklock")
				hasLock = false
				log.Printf("released lock.  snooze a bit more...")
				time.Sleep(2 * time.Second)

			}
			sid = renewSession(sid)
			log.Printf("sid is %s, index is %d, leader is %t", sid, index, hasLock)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// optimized out without this?
		if hasLock {
			log.Printf("HTTP trying to serve with lock")
			fmt.Fprintf(w, htmllocked, hostname, version)
		} else {
			log.Printf("HTTP trying to serve without lock")
			fmt.Fprintf(w, htmlunlocked, hostname, version)
		}
	})

	http.HandleFunc("/release", func(w http.ResponseWriter, r *http.Request) {
		if releaseLock(sid, "service/demo/worklock") {
			log.Printf("releasing lock from sid %s", sid)
			hasLock = false
		} else {
			fmt.Printf("trying to release lock when we don't have it, whoops")
		}
		fmt.Fprintf(w, "<html>cool story bro</html>")
	})

	log.Printf("HTTP service listening on %s", httpAddr)
	http.ListenAndServe(httpAddr, nil)

	// wait for cleanup, we deregister with this hook
	<-signalHandlerDone

}
