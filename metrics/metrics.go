package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/vokomarov/home-exporter/config"
)

type Metric interface {
	Run()
	Stop()
}

var (
	metrics            map[int]Metric
	homeInternetStatus = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "home_internet_status",
		Help: "Verifies internet connectivity of home router WAN IP",
	}, []string{"home", "host", "port", "method"})
)

func Listen() {
	Stop()

	metrics = make(map[int]Metric, len(config.Global.Homes))

	for i, home := range config.Global.Homes {
		metric := NewInternetStatusMetric(home)
		if metric == nil {
			continue
		}

		metrics[i] = metric

		go metrics[i].Run()
	}
}

func Stop() {
	for i := range metrics {
		if metrics[i] == nil {
			continue
		}

		metrics[i].Stop()
	}
}
