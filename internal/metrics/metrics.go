package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	orderLabel = "order"
)

var (
	issuedOrders = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "issued_orders_total",
		Help: "total number issued orders",
	}, []string{
		orderLabel,
	})
)

func AddIssuedOrders(count int) {
	issuedOrders.With(prometheus.Labels{
		orderLabel: "issued",
	}).Add(float64(count))
}
