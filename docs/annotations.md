# Infraflow Autoscale Annotations 说明文档

Infraflow Autoscale Operator 通过读取资源对象（Deployment、StatefulSet、DaemonSet）的 Annotations，动态管理 HPA（HorizontalPodAutoscaler）和 VPA（VerticalPodAutoscaler）配置。只需在资源对象的 Metadata 中添加特定 Annotation，即可启用或定制自动扩缩容策略。

## HPA（水平自动扩缩容）相关 Annotations

| Annotation Key | 类型 | 示例值 | 描述 |
|----------------|------|--------|------|
| `hpa.infraflow.co/minReplicas` | string | "2" | 最小副本数 |
| `hpa.infraflow.co/maxReplicas` | string | "10" | 最大副本数 |
| `cpu.hpa.infraflow.co/targetAverageUtilization` | string | "70" | CPU 使用率目标（百分比 %） |
| `cpu.hpa.infraflow.co/targetAverageValue` | string | "500m" | CPU 使用量目标（核数） |
| `memory.hpa.infraflow.co/targetAverageUtilization` | string | "75" | 内存使用率目标（百分比 %） |
| `memory.hpa.infraflow.co/targetAverageValue` | string | "512Mi" | 内存使用量目标（字节数） |

## External Metrics（Prometheus 自定义指标）相关 Annotations

| Annotation Key | 类型 | 示例值 | 描述 |
|----------------|------|--------|------|
| `prometheus.hpa.infraflow.co/metricName` | string | "http_requests_total" | Prometheus 指标名称 |
| `prometheus.hpa.infraflow.co/targetAverageValue` | string | "100" | 目标指标值（一般是每副本指标期望值） |

> 说明：使用 External Metrics 时，需搭配 Prometheus Adapter，并确保相关 Metric 已注册到 Kubernetes Metrics API。

## VPA（垂直自动扩缩容）相关 Annotations

| Annotation Key | 类型 | 示例值 | 描述 |
|----------------|------|--------|------|
| `vpa.infraflow.co/updateMode` | string | "Auto" , "Initial" , "Off" | VPA 更新模式。Auto 表示自动调整，Initial 表示仅初始化时设置，Off 禁用更新 |
| `cpu.vpa.infraflow.co/minAllowed` | string | "200m" | 容器允许的最小 CPU 资源限制 |
| `cpu.vpa.infraflow.co/maxAllowed` | string | "2" | 容器允许的最大 CPU 资源限制 |
| `memory.vpa.infraflow.co/minAllowed` | string | "256Mi" | 容器允许的最小内存资源限制 |
| `memory.vpa.infraflow.co/maxAllowed` | string | "4Gi" | 容器允许的最大内存资源限制 |
| `vpa.infraflow.co/resourcePolicy`	| map[string]string | `{ "containerPolicies": [...] }`|	PodResourcePolicy 配置，详细控制各容器的扩缩规则|
| `vpa.infraflow.co/containerPolicies` |	map[string]string |	`[{ "containerName": "app", "minAllowed": {"cpu": "200m"} }]` | ContainerResourcePolicy 列表，独立配置单个容器的资源策略|

>说明：
>
>`resourcePolicy` 是完整的 PodResourcePolicy JSON
>
>`containerPolicies` 是只指定 ContainerResourcePolicy 列表（内部合并到 resourcePolicy.containerPolicies 字段）。

## Finalizer

Infraflow Autoscaler Operator 自动为管理的 Workload 增加以下 Finalizer：

| Finalizer | 描述 |
|-----------|------|
| `finalizers.infraflow.co/autoscale` | 确保 Workload 删除时自动清理对应的 HPA 和 VPA，避免资源悬挂 | 

## 补充说明
 + 所有 Annotation 的值都必须是字符串格式。

 + CPU 单位：millicore（m），如 500m = 0.5 核。

 + 内存单位：如 Mi, Gi，例如 512Mi。

 + VPA 的 resourcePolicy / containerPolicies 需要提供合法 JSON，且符合 Kubernetes VPA API 结构。

 + External Metrics 需要正确部署 Prometheus Adapter。

