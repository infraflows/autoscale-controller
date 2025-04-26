/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"strings"

	"github.com/infraflows/autoscale-controller/pkg/consts"
	"github.com/infraflows/autoscale-controller/pkg/kube"
	"github.com/infraflows/autoscale-controller/pkg/metrics"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	autoscalingv1beta2 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta2"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type AutoScaleReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Event  record.EventRecorder
}

func init() {
	metrics.Init()
}

func (r *AutoScaleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconciling workload", "namespace", req.Namespace, "name", req.Name)

	var workload client.Object
	var kind string
	candidates := []client.Object{
		&appsv1.Deployment{},
		&appsv1.StatefulSet{},
		&appsv1.DaemonSet{},
	}
	for _, c := range candidates {
		err := r.Get(ctx, req.NamespacedName, c)
		if err == nil {
			workload = c
			kind = c.GetObjectKind().GroupVersionKind().Kind
			break
		}
		if !errors.IsNotFound(err) {
			logger.Error(err, "Failed to get workload")
			return ctrl.Result{}, err
		}
	}
	if workload == nil {
		logger.Info("Workload not found, skipping")
		return ctrl.Result{}, nil
	}

	if workload.GetDeletionTimestamp() != nil {
		// 如果删除中，执行清理逻辑
		if controllerutil.ContainsFinalizer(workload, consts.AutoScaleFinalizer) {
			_ = r.deleteHPA(ctx, workload)
			_ = r.deleteVPA(ctx, workload)
			controllerutil.RemoveFinalizer(workload, consts.AutoScaleFinalizer)
			return ctrl.Result{}, r.Update(ctx, workload)
		}
		return ctrl.Result{}, nil
	}
	if !controllerutil.ContainsFinalizer(workload, consts.AutoScaleFinalizer) {
		controllerutil.AddFinalizer(workload, consts.AutoScaleFinalizer)
		_ = r.Update(ctx, workload)
	}

	annotations := workload.GetAnnotations()
	if r.shouldManageHPA(annotations) {
		if err := r.reconcileHPA(ctx, workload, kind); err != nil {
			logger.Error(err, "Failed to reconcile HPA")
			return ctrl.Result{}, err
		}
	} else {
		if err := r.deleteHPA(ctx, workload); err != nil {
			logger.Error(err, "Failed to delete HPA")
			return ctrl.Result{}, err
		}
	}

	if r.shouldManageVPA(annotations) {
		if err := r.reconcileVPA(ctx, workload, kind); err != nil {
			logger.Error(err, "Failed to reconcile VPA")
			return ctrl.Result{}, err
		}
	} else {
		if err := r.deleteVPA(ctx, workload); err != nil {
			logger.Error(err, "Failed to delete VPA")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// reconcileHPA 协调Horizontal Pod Autoscale
// 1. 构建期望的HPA配置
// 2. 检查现有HPA是否存在
// 3. 创建新的HPA或更新现有的HPA
func (r *AutoScaleReconciler) reconcileHPA(ctx context.Context, workload client.Object, kind string) error {
	desired := kube.BuildDesiredHPA(workload, kind)
	current := &autoscalingv2.HorizontalPodAutoscaler{}
	err := r.Get(ctx, client.ObjectKeyFromObject(desired), current)
	if errors.IsNotFound(err) {
		controllerutil.SetControllerReference(workload, desired, r.Scheme)
		return r.Create(ctx, desired)
	} else if err != nil {
		return err
	}
	controllerutil.SetControllerReference(workload, current, r.Scheme)
	if !kube.EqualHPA(current, desired) {
		current.Spec = desired.Spec
		return r.Update(ctx, current)
	}
	return nil
}

// deleteHPA 删除与工作负载关联的HPA
func (r *AutoScaleReconciler) deleteHPA(ctx context.Context, workload client.Object) error {
	hpa := &autoscalingv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      workload.GetName(),
			Namespace: workload.GetNamespace(),
		},
	}
	return client.IgnoreNotFound(r.Delete(ctx, hpa))
}

// reconcileVPA 协调Vertical Pod Autoscaler
// 1. 构建期望的VPA配置
// 2. 检查现有VPA是否存在
// 3. 创建新的VPA或更新现有的VPA
func (r *AutoScaleReconciler) reconcileVPA(ctx context.Context, workload client.Object, kind string) error {
	desired, err := kube.BuildDesiredVPA(workload, kind)
	if err != nil {
		return err
	}
	current := &autoscalingv1beta2.VerticalPodAutoscaler{}
	err = r.Get(ctx, client.ObjectKeyFromObject(desired), current)
	if errors.IsNotFound(err) {
		controllerutil.SetControllerReference(workload, desired, r.Scheme)
		return r.Create(ctx, desired)
	} else if err != nil {
		return err
	}
	controllerutil.SetControllerReference(workload, current, r.Scheme)
	if !kube.EqualVPA(current, desired) {
		current.Spec = desired.Spec
		return r.Update(ctx, current)
	}
	return nil
}

// deleteVPA 删除与工作负载关联的VPA
func (r *AutoScaleReconciler) deleteVPA(ctx context.Context, workload client.Object) error {
	vpa := &autoscalingv1beta2.VerticalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      workload.GetName(),
			Namespace: workload.GetNamespace(),
		},
	}
	return client.IgnoreNotFound(r.Delete(ctx, vpa))
}

// shouldManageHPA 检查工作负载的注解是否包含HPA相关的配置
// 支持的注解前缀：
// - cpu.hpa.infraflow.co/
// - memory.hpa.infraflow.co/
func (r *AutoScaleReconciler) shouldManageHPA(annotations map[string]string) bool {
	for key := range annotations {
		if strings.HasPrefix(key, "cpu.hpa.infraflow.co/") || strings.HasPrefix(key, "memory.hpa.infraflow.co/") {
			return true
		}
	}
	return false
}

// shouldManageVPA 检查工作负载的注解是否包含VPA相关的配置
// 支持的注解前缀：
// - cpu.vpa.infraflow.co/
// - memory.vpa.infraflow.co/
// - vpa.infraflow.co/
func (r *AutoScaleReconciler) shouldManageVPA(annotations map[string]string) bool {
	for key := range annotations {
		if strings.HasPrefix(key, "cpu.vpa.infraflow.co/") || strings.HasPrefix(key, "memory.vpa.infraflow.co/") || strings.HasPrefix(key, "vpa.infraflow.co/") {
			return true
		}
	}
	return false
}

func (r *AutoScaleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		For(&appsv1.StatefulSet{}).
		For(&appsv1.DaemonSet{}).
		Owns(&autoscalingv2.HorizontalPodAutoscaler{}).
		Owns(&autoscalingv1beta2.VerticalPodAutoscaler{}).
		Complete(r)
}
