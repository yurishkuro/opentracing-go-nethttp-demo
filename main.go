package main

import (
	"flag"
	"log"

	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport/zipkin"
)

var (
	zipkinURL = flag.String("url",
		"http://localhost:9411/api/v1/spans", "Zipkin server URL")
	serverPort = flag.String("port", "8000", "server port")
	actorKind  = flag.String("actor", "server", "server or client")
)

const (
	server = "server"
	client = "client"
)

func main() {
	flag.Parse()

	if *actorKind != server && *actorKind != client {
		log.Fatal("Please specify '-actor server' or '-actor client'")
	}

	// Jaeger tracer can be initialized with a transport that will
	// report tracing Spans to a Zipkin backend
	transport, err := zipkin.NewHTTPTransport(
		*zipkinURL,
		zipkin.HTTPBatchSize(1),
		zipkin.HTTPLogger(jaeger.StdLogger),
	)
	if err != nil {
		log.Fatalf("Cannot initialize HTTP transport: %v", err)
	}
	// create Jaeger tracer
	tracer, closer := jaeger.NewTracer(
		*actorKind,
		jaeger.NewConstSampler(true), // sample all traces
		jaeger.NewRemoteReporter(transport, nil),
	)

	if *actorKind == server {
		runServer(tracer)
		return
	}

	runClient(tracer)

	// Close the tracer to guarantee that all spans that could
	// be still buffered in memory are sent to the tracing backend
	closer.Close()
}
