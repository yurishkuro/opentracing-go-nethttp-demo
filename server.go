package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
)

func getTime(w http.ResponseWriter, r *http.Request) {
	log.Print("Received getTime request")
	t := time.Now()
	ts := t.Format("Mon Jan _2 15:04:05 2006")
	io.WriteString(w, fmt.Sprintf("The time is %s", ts))
}

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r,
		fmt.Sprintf("http://localhost:%s/gettime", *serverPort), 301)
}

func runServer(tracer opentracing.Tracer) {
	http.HandleFunc("/gettime", getTime)
	http.HandleFunc("/", redirect)
	log.Printf("Starting server on port %s", *serverPort)
	err := http.ListenAndServe(
		fmt.Sprintf(":%s", *serverPort),
		// use nethttp.Middleware to enable OpenTracing for server
		nethttp.Middleware(tracer, http.DefaultServeMux))
	if err != nil {
		log.Fatalf("Cannot start server: %s", err)
	}
}
