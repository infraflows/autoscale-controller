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
	"time"

	"github.com/infraflows/autoscale-controller/pkg/consts"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("AutoScale Controller", func() {
	const (
		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When reconciling a Deployment with HPA annotations", func() {
		var (
			ctx            context.Context
			deployment     *appsv1.Deployment
			namespace      string
			deploymentName string
		)

		BeforeEach(func() {
			ctx = context.Background()
			namespace = "test-namespace"
			deploymentName = "test-deployment"

			// 创建测试命名空间
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			Expect(k8sClient.Create(ctx, ns)).Should(Succeed())

			// 创建测试 Deployment
			deployment = &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      deploymentName,
					Namespace: namespace,
					Annotations: map[string]string{
						consts.HPAMinReplicas:                 "2",
						consts.HPAMaxReplicas:                 "10",
						consts.HPACpuTargetAverageUtilization: "80",
					},
				},
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app": deploymentName,
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app": deploymentName,
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "nginx",
									Image: "nginx:1.14.2",
								},
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, deployment)).Should(Succeed())
		})

		AfterEach(func() {
			// 清理测试资源
			Expect(k8sClient.Delete(ctx, deployment)).Should(Succeed())
			Expect(k8sClient.Delete(ctx, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}})).Should(Succeed())
		})

		It("Should create HPA when deployment has HPA annotations", func() {
			// 等待 HPA 创建
			Eventually(func() error {
				hpa := &autoscalingv2.HorizontalPodAutoscaler{}
				return k8sClient.Get(ctx, types.NamespacedName{
					Name:      deploymentName,
					Namespace: namespace,
				}, hpa)
			}, timeout, interval).Should(Succeed())

			// 验证 HPA 配置
			hpa := &autoscalingv2.HorizontalPodAutoscaler{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{
				Name:      deploymentName,
				Namespace: namespace,
			}, hpa)).Should(Succeed())

			Expect(hpa.Spec.MinReplicas).Should(HaveValue(Equal(int32(2))))
			Expect(hpa.Spec.MaxReplicas).Should(Equal(int32(10)))
			Expect(hpa.Spec.Metrics).Should(HaveLen(1))
			Expect(hpa.Spec.Metrics[0].Resource.Name).Should(Equal(corev1.ResourceCPU))
			Expect(hpa.Spec.Metrics[0].Resource.Target.AverageUtilization).Should(HaveValue(Equal(int32(80))))
		})

		It("Should update HPA when deployment annotations change", func() {
			// 更新 Deployment 注解
			Eventually(func() error {
				if err := k8sClient.Get(ctx, types.NamespacedName{
					Name:      deploymentName,
					Namespace: namespace,
				}, deployment); err != nil {
					return err
				}
				deployment.Annotations[consts.HPAMaxReplicas] = "15"
				return k8sClient.Update(ctx, deployment)
			}, timeout, interval).Should(Succeed())

			// 验证 HPA 更新
			Eventually(func() int32 {
				hpa := &autoscalingv2.HorizontalPodAutoscaler{}
				if err := k8sClient.Get(ctx, types.NamespacedName{
					Name:      deploymentName,
					Namespace: namespace,
				}, hpa); err != nil {
					return 0
				}
				return hpa.Spec.MaxReplicas
			}, timeout, interval).Should(Equal(int32(15)))
		})

		It("Should delete HPA when deployment annotations are removed", func() {
			// 移除 HPA 注解
			Eventually(func() error {
				if err := k8sClient.Get(ctx, types.NamespacedName{
					Name:      deploymentName,
					Namespace: namespace,
				}, deployment); err != nil {
					return err
				}
				delete(deployment.Annotations, consts.HPAMinReplicas)
				delete(deployment.Annotations, consts.HPAMaxReplicas)
				delete(deployment.Annotations, consts.HPACpuTargetAverageUtilization)
				return k8sClient.Update(ctx, deployment)
			}, timeout, interval).Should(Succeed())

			// 验证 HPA 被删除
			Eventually(func() bool {
				hpa := &autoscalingv2.HorizontalPodAutoscaler{}
				err := k8sClient.Get(ctx, types.NamespacedName{
					Name:      deploymentName,
					Namespace: namespace,
				}, hpa)
				return err != nil
			}, timeout, interval).Should(BeTrue())
		})
	})
})
