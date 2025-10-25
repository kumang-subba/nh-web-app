package main

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	s := &http.Server{
		Addr: ":3000",
	}

	var (
		hwReg *prometheus.Registry
		ph    prometheus.Histogram
	)

	hwReg = prometheus.NewRegistry()

	ph = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "req_time_hello_world",
		Help:    "How long each hello world request took in milliseconds",
		Buckets: []float64{10, 25, 50, 100, 250, 500, 1000, 1500, 2000},
	})

	hwReg.MustRegister(ph)

	gcReg := prometheus.NewRegistry()

	gcReg.MustRegister(collectors.NewGoCollector(
		collectors.WithGoCollectorRuntimeMetrics(collectors.MetricsGC),
	))

	hw := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		t := time.Now()
		w.Write([]byte("Hello world"))
		ph.Observe(float64(time.Since(t).Milliseconds()))
	}

	http.Handle("/metrics", promhttp.HandlerFor(hwReg, promhttp.HandlerOpts{Registry: hwReg}))
	http.Handle("/metricsgc", promhttp.HandlerFor(gcReg, promhttp.HandlerOpts{Registry: gcReg}))
	http.HandleFunc("/", hw)
	s.ListenAndServe()
}
