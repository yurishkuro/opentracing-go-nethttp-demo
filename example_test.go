package main

import (
	"net/http"
	"net/http/httptrace"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
)

func ExampleAskGoogle() {
	AskGoogle(context.Background())
}

// We will talk about this later
var tracer opentracing.Tracer

func AskGoogle(ctx context.Context) error {
	// retrieve current Span from Context
	var parentCtx opentracing.SpanContext
	parentSpan := opentracing.SpanFromContext(ctx)
	if parentSpan != nil {
		parentCtx = parentSpan.Context()
	}

	// start a new Span to wrap HTTP request
	span := tracer.StartSpan(
		"ask google",
		opentracing.ChildOf(parentCtx),
	)

	// make sure the Span is finished once we're done
	defer span.Finish()

	// make the Span current in the context
	ctx = opentracing.ContextWithSpan(ctx, span)

	// now prepare the request
	req, err := http.NewRequest("GET", "http://google.com", nil)
	if err != nil {
		return err
	}

	// attach client trace to the Context, and Context to request
	trace := NewClientTrace(span)
	ctx = httptrace.WithClientTrace(ctx, trace)
	req = req.WithContext(ctx)

	// execute the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	// Google home page is not too exciting, so ignore the result
	res.Body.Close()
	return nil
}

func NewClientTrace(span opentracing.Span) *httptrace.ClientTrace {
	trace := &clientTrace{span: span}
	return &httptrace.ClientTrace{
		DNSStart: trace.dnsStart,
		DNSDone:  trace.dnsDone,
	}
}

// clientTrace holds a reference to the Span and
// provides methods used as ClientTrace callbacks
type clientTrace struct {
	span opentracing.Span
}

func (h *clientTrace) dnsStart(info httptrace.DNSStartInfo) {
	h.span.LogKV(
		log.String("event", "DNS start"),
		log.Object("host", info.Host),
	)
}

func (h *clientTrace) dnsDone(httptrace.DNSDoneInfo) {
	h.span.LogKV(log.String("event", "DNS done"))
}
