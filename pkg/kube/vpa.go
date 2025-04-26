package kube

import (
	"encoding/json"
	"fmt"

	"github.com/infraflows/autoscale-controller/pkg/consts"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	autoscalingv1beta2 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// UpdateModeAuto 表示自动调整模式
	UpdateModeAuto = "Auto"
	// UpdateModeInitial 表示仅初始化时设置模式
	UpdateModeInitial = "Initial"
	// UpdateModeOff 表示禁用更新模式
	UpdateModeOff = "Off"
)

// ValidateUpdateMode 验证更新模式是否有效
func ValidateUpdateMode(mode string) error {
	switch mode {
	case UpdateModeAuto, UpdateModeInitial, UpdateModeOff:
		return nil
	default:
		return fmt.Errorf("invalid update mode: %s, must be one of: %s, %s, %s",
			mode, UpdateModeAuto, UpdateModeInitial, UpdateModeOff)
	}
}

// BuildDesiredVPA 根据工作负载的注解构建期望的Vertical Pod Autoscale配置
// 支持的注解：
// - vpa.infraflow.co/update-mode: 更新模式（Auto/Initial/Off）
// - vpa.infraflow.co/resource-policy: 资源策略（JSON格式）
// 如果没有指定更新模式，默认使用Auto模式
func BuildDesiredVPA(workload client.Object, kind string) (*autoscalingv1beta2.VerticalPodAutoscaler, error) {
	annotations := workload.GetAnnotations()
	vpa := &autoscalingv1beta2.VerticalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      workload.GetName(),
			Namespace: workload.GetNamespace(),
		},
		Spec: autoscalingv1beta2.VerticalPodAutoscalerSpec{
			TargetRef: &autoscalingv1.CrossVersionObjectReference{
				APIVersion: "apps/v1",
				Kind:       kind,
				Name:       workload.GetName(),
			},
		},
	}

	if val, ok := annotations[consts.VPAAnnotationUpdateMode]; ok {
		if err := ValidateUpdateMode(val); err != nil {
			return nil, err
		}
		mode := autoscalingv1beta2.UpdateMode(val)
		vpa.Spec.UpdatePolicy = &autoscalingv1beta2.PodUpdatePolicy{
			UpdateMode: &mode,
		}
	} else {
		// 如果没有指定更新模式，默认使用 Auto 模式
		mode := autoscalingv1beta2.UpdateMode(UpdateModeAuto)
		vpa.Spec.UpdatePolicy = &autoscalingv1beta2.PodUpdatePolicy{
			UpdateMode: &mode,
		}
	}

	if val, ok := annotations[consts.VPAAnnotationResourcePolicy]; ok {
		var policy autoscalingv1beta2.PodResourcePolicy
		if err := json.Unmarshal([]byte(val), &policy); err == nil {
			vpa.Spec.ResourcePolicy = &policy
		}
	}
	return vpa, nil
}

// EqualVPA 比较两个VPA配置是否相等
// 比较内容包括：
// - 目标引用名称
// - 更新模式
func EqualVPA(a, b *autoscalingv1beta2.VerticalPodAutoscaler) bool {
	return a.Spec.TargetRef != nil && b.Spec.TargetRef != nil &&
		a.Spec.TargetRef.Name == b.Spec.TargetRef.Name &&
		a.Spec.UpdatePolicy != nil && b.Spec.UpdatePolicy != nil &&
		*a.Spec.UpdatePolicy.UpdateMode == *b.Spec.UpdatePolicy.UpdateMode
}
