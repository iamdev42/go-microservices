package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"github.com/opentracing/opentracing-go"
	"time"
)

func main() {
	log.Println("I am server")

	// create Jaeger tracer
	// Sample configuration for testing. Use constant sampling to sample every trace
	// and enable LogSpan to log every span via configured Logger.
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}

	// Example logger and metrics factory. Use github.com/uber/jaeger-client-go/log
	// and github.com/uber/jaeger-lib/metrics respectively to bind to real logging and metrics
	// frameworks.
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	// Initialize tracer with a logger and a metrics factory
	closer, err := cfg.InitGlobalTracer(
		"go-micro-interface",
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		log.Printf("Could not initialize jaeger tracer: %s", err.Error())
		return
	}
	defer closer.Close()

	// Start server and expose endpoints
	http.HandleFunc("/", home)
	http.HandleFunc("/wscall", wsCall)
	http.HandleFunc("/db-write", callDatabaseService)

	if err := http.ListenAndServe(":8086", nil); err != nil {
		panic(err)
	}
}

// Function used to showcase how to do simple tests -> interface_test.go file
func Sum(x, y int) (sum int){
	return x + y
}

// Simple external GET REST call to google
func wsCall(w http.ResponseWriter, r *http.Request) {

	log.Println("Calling google")

	var responseMessage string
	response, err := http.Get("https://www.google.com")
	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
		responseMessage = "The HTTP request failed with error"
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		responseMessage = string(data) + string(data)
	}
	log.Print("Response message is: ")
	log.Print(responseMessage)
}

// Function showcases tracing functionality
func home(w http.ResponseWriter, r *http.Request) {
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("render-home")

	firstExecution(span)
	secondExecution(span)

	message := "Home page of awesome interface"
	w.Write([]byte(message))

	span.Finish()
}

func firstExecution(rootSpan opentracing.Span) {
	span := rootSpan.Tracer().StartSpan("first-execution", opentracing.ChildOf(rootSpan.Context()))
	time.Sleep(1 * time.Second)
	span.Finish()
}

func secondExecution(rootSpan opentracing.Span) {
	span := rootSpan.Tracer().StartSpan("second-execution", opentracing.ChildOf(rootSpan.Context()))
	time.Sleep(3 * time.Second)
	span.Finish()
}

// Function showcases call to another microservice
func callDatabaseService(w http.ResponseWriter, r *http.Request) {

	var responseMessage string
	response, err := http.Get("http://localhost:8082/get")
	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
		responseMessage = "The HTTP request failed with error"
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		responseMessage = string(data) + string(data)
	}

	w.Write([]byte(responseMessage))
}
