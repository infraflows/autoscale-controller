package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	ReconcileTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "infraflow_autoscaler_reconcile_total",
			Help: "Total number of reconciliations",
		},
		[]string{"kind"},
	)
)

func Init() {
	metrics.Registry.MustRegister(ReconcileTotal)
}
