package kube

import (
	"strconv"

	"github.com/infraflows/autoscale-controller/pkg/consts"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// BuildDesiredHPA 根据工作负载的注解构建期望的Horizontal Pod autoscale配置
// 支持的注解：
// - hpa.infraflow.co/min-replicas: 最小副本数
// - hpa.infraflow.co/max-replicas: 最大副本数
// - cpu.hpa.infraflow.co/target-average-utilization: CPU利用率目标
// - cpu.hpa.infraflow.co/target-average-value: CPU使用量目标
// - memory.hpa.infraflow.co/target-average-utilization: 内存利用率目标
// - memory.hpa.infraflow.co/target-average-value: 内存使用量目标
func BuildDesiredHPA(workload client.Object, kind string) *autoscalingv2.HorizontalPodAutoscaler {
	annotations := workload.GetAnnotations()
	hpa := &autoscalingv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      workload.GetName(),
			Namespace: workload.GetNamespace(),
		},
		Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
				APIVersion: "apps/v1",
				Kind:       kind,
				Name:       workload.GetName(),
			},
		},
	}

	if val, ok := annotations[consts.HPAMinReplicas]; ok {
		if v, err := strconv.Atoi(val); err == nil {
			min := int32(v)
			hpa.Spec.MinReplicas = &min
		}
	}
	if val, ok := annotations[consts.HPAMaxReplicas]; ok {
		if v, err := strconv.Atoi(val); err == nil {
			hpa.Spec.MaxReplicas = int32(v)
		}
	}

	metrics := []autoscalingv2.MetricSpec{}
	if val, ok := annotations[consts.HPACpuTargetAverageUtilization]; ok {
		if target, err := strconv.Atoi(val); err == nil {
			t := int32(target)
			metrics = append(metrics, CPUUtilizationMetric(t))
		}
	}
	if val, ok := annotations[consts.HPACpuTargetAverageValue]; ok {
		if quantity, err := resource.ParseQuantity(val); err == nil {
			metrics = append(metrics, CPUValueMetric(quantity))
		}
	}
	if val, ok := annotations[consts.HPAMemoryTargetAverageUtilization]; ok {
		if target, err := strconv.Atoi(val); err == nil {
			t := int32(target)
			metrics = append(metrics, MemoryUtilizationMetric(t))
		}
	}
	if val, ok := annotations[consts.HPAMemoryTargetAverageValue]; ok {
		if quantity, err := resource.ParseQuantity(val); err == nil {
			metrics = append(metrics, MemoryValueMetric(quantity))
		}
	}

	hpa.Spec.Metrics = metrics
	return hpa
}

// CPUUtilizationMetric 基于CPU利用率的HPA指标配置
// target: 目标CPU利用率百分比
func CPUUtilizationMetric(target int32) autoscalingv2.MetricSpec {
	return autoscalingv2.MetricSpec{
		Type: autoscalingv2.ResourceMetricSourceType,
		Resource: &autoscalingv2.ResourceMetricSource{
			Name: corev1.ResourceCPU,
			Target: autoscalingv2.MetricTarget{
				Type:               autoscalingv2.UtilizationMetricType,
				AverageUtilization: &target,
			},
		},
	}
}

// CPUValueMetric 基于CPU使用量的HPA指标配置
// quantity: 目标CPU使用量
func CPUValueMetric(quantity resource.Quantity) autoscalingv2.MetricSpec {
	return autoscalingv2.MetricSpec{
		Type: autoscalingv2.ResourceMetricSourceType,
		Resource: &autoscalingv2.ResourceMetricSource{
			Name: corev1.ResourceCPU,
			Target: autoscalingv2.MetricTarget{
				Type:         autoscalingv2.AverageValueMetricType,
				AverageValue: &quantity,
			},
		},
	}
}

// MemoryUtilizationMetric 基于内存利用率的HPA指标配置
// target: 目标内存利用率百分比
func MemoryUtilizationMetric(target int32) autoscalingv2.MetricSpec {
	return autoscalingv2.MetricSpec{
		Type: autoscalingv2.ResourceMetricSourceType,
		Resource: &autoscalingv2.ResourceMetricSource{
			Name: corev1.ResourceMemory,
			Target: autoscalingv2.MetricTarget{
				Type:               autoscalingv2.UtilizationMetricType,
				AverageUtilization: &target,
			},
		},
	}
}

// MemoryValueMetric 基于内存使用量的HPA指标配置
// quantity: 目标内存使用量
func MemoryValueMetric(quantity resource.Quantity) autoscalingv2.MetricSpec {
	return autoscalingv2.MetricSpec{
		Type: autoscalingv2.ResourceMetricSourceType,
		Resource: &autoscalingv2.ResourceMetricSource{
			Name: corev1.ResourceMemory,
			Target: autoscalingv2.MetricTarget{
				Type:         autoscalingv2.AverageValueMetricType,
				AverageValue: &quantity,
			},
		},
	}
}

// EqualHPA 比较两个HPA配置是否相等
// 比较内容包括：
// - 目标引用
// - 指标配置
// - 最小副本数
// - 最大副本数
func EqualHPA(a, b *autoscalingv2.HorizontalPodAutoscaler) bool {
	return a.Spec.ScaleTargetRef == b.Spec.ScaleTargetRef &&
		EqualMetrics(a.Spec.Metrics, b.Spec.Metrics) &&
		EqualInt32Ptr(a.Spec.MinReplicas, b.Spec.MinReplicas) &&
		a.Spec.MaxReplicas == b.Spec.MaxReplicas
}

// EqualInt32Ptr 比较两个int32指针是否相等
func EqualInt32Ptr(a, b *int32) bool {
	if a == nil && b == nil {
		return true
	}
	if a != nil && b != nil && *a == *b {
		return true
	}
	return false
}

// EqualMetrics 比较两个指标配置数组是否相等
// 比较内容包括：
// - 指标类型
// - 资源名称
func EqualMetrics(a, b []autoscalingv2.MetricSpec) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Type != b[i].Type || a[i].Resource.Name != b[i].Resource.Name {
			return false
		}
	}
	return true
}
