package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var pageCounter = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "page_counter",
	Help: "requests per page",
}, []string{"path"})
