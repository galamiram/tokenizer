package serve

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jdkato/prose/v2"
	"github.com/mpvl/unique"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusHttpMetric struct {
	Prefix                string
	ClientConnected       prometheus.Gauge
	TransactionTotal      *prometheus.CounterVec
	ResponseTimeHistogram *prometheus.HistogramVec
	Buckets               []float64
}

func InitPrometheusHttpMetric(prefix string, buckets []float64) *PrometheusHttpMetric {
	phm := PrometheusHttpMetric{
		Prefix: prefix,
		ClientConnected: promauto.NewGauge(prometheus.GaugeOpts{
			Name: prefix + "_client_connected",
			Help: "Number of active client connections",
		}),
		TransactionTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: prefix + "_requests_total",
			Help: "total HTTP requests processed",
		}, []string{"code", "method"},
		),
		ResponseTimeHistogram: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    prefix + "_response_time",
			Help:    "Histogram of response time for handler",
			Buckets: buckets,
		}, []string{"handler", "method"}),
	}

	return &phm
}

func Tokenize(w http.ResponseWriter, r *http.Request) {
	var tokenSlice []string
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	doc, err := prose.NewDocument(string(body))
	for _, tok := range doc.Tokens() {
		tokenSlice = append(tokenSlice, tok.Text)
	}

	unique.Sort(unique.StringSlice{&tokenSlice})
	jArr, err := json.Marshal(tokenSlice)
	if err != nil {
		http.Error(w, "could not parse request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jArr)
}

func (phm *PrometheusHttpMetric) WrapHandler(handlerLabel string, handlerFunc http.HandlerFunc) http.Handler {
	handle := http.HandlerFunc(handlerFunc)
	wrappedHandler := promhttp.InstrumentHandlerInFlight(phm.ClientConnected,
		promhttp.InstrumentHandlerCounter(phm.TransactionTotal,
			promhttp.InstrumentHandlerDuration(phm.ResponseTimeHistogram.MustCurryWith(prometheus.Labels{"handler": handlerLabel}),
				handle),
		),
	)
	return wrappedHandler
}
