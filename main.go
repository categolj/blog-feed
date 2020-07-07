package main

import (
	"fmt"
	"github.com/categolj/blog-feed/handler"
	"github.com/gorilla/mux"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	"github.com/openzipkin/zipkin-go/model"
	reporterhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"log"
	"net/http"
	"os"
	"strconv"
)

func newTracer(port uint16) (*zipkin.Tracer, error) {
	var zipkinUrl string
	if zipkinUrl = os.Getenv("ZIPKIN_URL"); len(zipkinUrl) == 0 {
		zipkinUrl = "http://localhost:9411"
	}

	// The reporter sends traces to zipkin server
	reporter := reporterhttp.NewReporter(zipkinUrl + "/api/v2/spans")

	// Local endpoint represent the local service information
	localEndpoint := &model.Endpoint{ServiceName: "blog:blog-feed", Port: port}

	// Sampler tells you which traces are going to be sampled or not. In this case we will record 100% (1.00) of traces.
	sampler, err := zipkin.NewCountingSampler(1)
	if err != nil {
		return nil, err
	}

	t, err := zipkin.NewTracer(
		reporter,
		zipkin.WithSampler(sampler),
		zipkin.WithLocalEndpoint(localEndpoint),
	)
	if err != nil {
		return nil, err
	}

	return t, err
}

func main() {
	var port string
	if port = os.Getenv("PORT"); len(port) == 0 {
		port = "4000"
	}
	p, _ := strconv.Atoi(port)
	tracer, err := newTracer(uint16(p))
	if err != nil {
		log.Fatal(err)
	}

	client, err := zipkinhttp.NewClient(tracer, zipkinhttp.ClientTrace(true))
	if err != nil {
		log.Fatalf("unable to create client: %+v\n", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", handler.FeedFactory(client))
	r.Use(zipkinhttp.NewServerMiddleware(
		tracer,
		zipkinhttp.SpanName("rss"),
	))

	log.Printf(fmt.Sprintf("Listening at %s", port))
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("ListenAndServe:", port)
	}
}
