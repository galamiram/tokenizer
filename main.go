package main

import (
	"log"
	"net/http"
	"os"

	"github.com/galamiram/tokenizer/cmd/serve"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	phm := serve.InitPrometheusHttpMetric("tokenizer", prometheus.LinearBuckets(0, 5, 20))

	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/tokenize", phm.WrapHandler("tokenize", serve.Tokenize))

	port := os.Getenv("LISTENING_PORT")

	if port == "" {
		port = "8080"
	}
	log.Printf("listening on port:%s", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Failed to start server:%v", err)
	}
}
