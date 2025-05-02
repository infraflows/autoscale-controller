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
	"time"

	"github.com/infraflows/autoscale-controller/pkg/consts"
	"github.com/infraflows/autoscale-controller/pkg/kube"
	"github.com/infraflows/autoscale-controller/pkg/metrics"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	// v1 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
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
	//logger.V(1).Info("Reconciling workload", "namespace", req.Namespace, "name", req.Name)

	workload, kind, err := r.getWorkload(ctx, req)
	if err != nil {
		logger.Error(err, "Failed to get workload")
		return ctrl.Result{}, err
	}

	if workload == nil {
		logger.V(1).Info("There are no matching workloads.")
		return ctrl.Result{RequeueAfter: 3 * time.Second}, nil
	}

	// 处理 finalizer
	cleanupFn := func(ctx context.Context, obj client.Object) error {
		return r.deleteHPA(ctx, obj)
	}
	if err := kube.HandleFinalizerWithCleanup(ctx, r.Client, workload, consts.AutoScaleFinalizer, logger, cleanupFn); err != nil {
		logger.Error(err, "Failed to handle finalizer")
		return ctrl.Result{}, err
	}

	annotations := workload.GetAnnotations()
	if r.shouldManageHPA(annotations) {
		logger.Info("Found HPA annotations, starting async reconciliation",
			"namespace", req.Namespace,
			"name", req.Name,
			"kind", kind)

		// 异步处理 HPA 创建/更新
		go func() {
			asyncCtx := context.Background()
			if err := r.reconcileHPA(asyncCtx, workload, kind); err != nil {
				logger.Error(err, "Failed to reconcile HPA in async process")
			}
		}()
	} else {
		// 如果没有 HPA 注解，静默删除可能存在的 HPA
		if err := r.deleteHPA(ctx, workload); err != nil {
			logger.V(1).Info("Failed to delete HPA", "error", err)
		}
	}

	// if r.shouldManageVPA(annotations) {
	// 	if err := r.reconcileVPA(ctx, workload, kind); err != nil {
	// 		logger.Error(err, "Failed to reconcile VPA")
	// 		return ctrl.Result{}, err
	// 	}
	// } else {
	// 	if err := r.deleteVPA(ctx, workload); err != nil {
	// 		logger.Error(err, "Failed to delete VPA")
	// 		return ctrl.Result{}, err
	// 	}
	// }

	return ctrl.Result{RequeueAfter: 3 * time.Second}, nil
}

// getWorkload 获取工作负载
func (r *AutoScaleReconciler) getWorkload(ctx context.Context, req ctrl.Request) (client.Object, string, error) {
	candidates := []client.Object{
		&appsv1.Deployment{},
		&appsv1.StatefulSet{},
		&appsv1.DaemonSet{},
	}

	for _, c := range candidates {
		err := r.Get(ctx, req.NamespacedName, c)
		if err == nil {
			return c, c.GetObjectKind().GroupVersionKind().Kind, nil
		}
		if !errors.IsNotFound(err) {
			return nil, "", err
		}
	}

	return nil, "", nil
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
// func (r *AutoScaleReconciler) reconcileVPA(ctx context.Context, workload client.Object, kind string) error {
// 	desired, err := kube.BuildDesiredVPA(workload, kind)
// 	if err != nil {
// 		return err
// 	}
// 	current := &v1.VerticalPodAutoscaler{}
// 	err = r.Get(ctx, client.ObjectKeyFromObject(desired), current)
// 	if errors.IsNotFound(err) {
// 		controllerutil.SetControllerReference(workload, desired, r.Scheme)
// 		return r.Create(ctx, desired)
// 	} else if err != nil {
// 		return err
// 	}
// 	controllerutil.SetControllerReference(workload, current, r.Scheme)
// 	if !kube.EqualVPA(current, desired) {
// 		current.Spec = desired.Spec
// 		return r.Update(ctx, current)
// 	}
// 	return nil
// }

// deleteVPA 删除与工作负载关联的VPA
// func (r *AutoScaleReconciler) deleteVPA(ctx context.Context, workload client.Object) error {
// 	vpa := &v1.VerticalPodAutoscaler{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      workload.GetName(),
// 			Namespace: workload.GetNamespace(),
// 		},
// 	}
// 	return client.IgnoreNotFound(r.Delete(ctx, vpa))
// }

// shouldManageHPA 检查工作负载的注解是否包含HPA相关的配置
// 支持的注解前缀：
// - hpa.infraflow.co/
func (r *AutoScaleReconciler) shouldManageHPA(annotations map[string]string) bool {
	for key := range annotations {
		if strings.HasPrefix(key, "hpa.infraflow.co/") {
			return true
		}
	}
	return false
}

// shouldManageVPA 检查工作负载的注解是否包含VPA相关的配置
// 支持的注解前缀：
// - vpa.infraflow.co/
// func (r *AutoScaleReconciler) shouldManageVPA(annotations map[string]string) bool {
// 	for key := range annotations {
// 		if strings.HasPrefix(key, "vpa.infraflow.co/") {
// 			return true
// 		}
// 	}
// 	return false
// }

func (r *AutoScaleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		Watches(&appsv1.StatefulSet{}, handler.EnqueueRequestsFromMapFunc(r.findObjectsForStatefulSet)).
		Watches(&appsv1.DaemonSet{}, handler.EnqueueRequestsFromMapFunc(r.findObjectsForDaemonSet)).
		Owns(&autoscalingv2.HorizontalPodAutoscaler{}).
		//Owns(&v1.VerticalPodAutoscaler{}).
		Complete(r)
}

// findObjectsForStatefulSet 为 StatefulSet 查找关联的对象
func (r *AutoScaleReconciler) findObjectsForStatefulSet(ctx context.Context, obj client.Object) []reconcile.Request {
	sts := obj.(*appsv1.StatefulSet)
	return []reconcile.Request{
		{
			NamespacedName: types.NamespacedName{
				Name:      sts.Name,
				Namespace: sts.Namespace,
			},
		},
	}
}

// findObjectsForDaemonSet 为 DaemonSet 查找关联的对象
func (r *AutoScaleReconciler) findObjectsForDaemonSet(ctx context.Context, obj client.Object) []reconcile.Request {
	ds := obj.(*appsv1.DaemonSet)
	return []reconcile.Request{
		{
			NamespacedName: types.NamespacedName{
				Name:      ds.Name,
				Namespace: ds.Namespace,
			},
		},
	}
}
