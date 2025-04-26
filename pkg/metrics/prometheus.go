package metrics

import (
	"k8s.io/apimachinery/pkg/api/resource"

	autoscalingv2 "k8s.io/api/autoscaling/v2"
)

func PrometheusExternalMetric(name string, quantity resource.Quantity) autoscalingv2.MetricSpec {
	return autoscalingv2.MetricSpec{
		Type: autoscalingv2.ExternalMetricSourceType,
		External: &autoscalingv2.ExternalMetricSource{
			Metric: autoscalingv2.MetricIdentifier{
				Name: name,
			},
			Target: autoscalingv2.MetricTarget{
				Type:         autoscalingv2.AverageValueMetricType,
				AverageValue: &quantity,
			},
		},
	}
}
