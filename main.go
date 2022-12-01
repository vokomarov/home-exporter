package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/vokomarov/home-exporter/config"
	"github.com/vokomarov/home-exporter/metrics"
	"github.com/vokomarov/home-exporter/telegram"
)

func main() {
	var err error

	if err = config.Load(); err != nil {
		panic(err)
	}

	if telegram.Bot, err = telegram.NewBot(); err != nil {
		panic(err)
	}

	metrics.Listen()

	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":2112", nil); err != nil {
		panic(err)
	}
}
