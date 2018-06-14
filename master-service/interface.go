package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

func main() {
	log.Println("I am server")

	// create Jaeger tracer
	setupTracing()

	// http.HandleFunc("/", home)
	// http.HandleFunc("/call", callDatabaseService)
	// //http.HandleFunc("/db-write", writeToDb)

	// if err := http.ListenAndServe(":8080", nil); err != nil {
	// 	panic(err)
	// }
}

func home(w http.ResponseWriter, r *http.Request) {
	message := "Home page of awesome interface"
	w.Write([]byte(message))
}

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

func setupTracing() {
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
		"serviceName",
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		log.Printf("Could not initialize jaeger tracer: %s", err.Error())
		return
	}
	defer closer.Close()
}
