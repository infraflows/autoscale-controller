# Autoscale Controller

Autoscale Controller 是一个用于自动管理 Kubernetes 工作负载扩缩容生命周期的 Controller。
通过为 Deployment、StatefulSet、DaemonSet 等工作负载资源添加标准化的 Annotations，自动创建、更新和清理对应的 Horizontal Pod Autoscale (HPA) 和 Vertical Pod Autoscale (VPA) 资源，实现智能化的扩缩容管理。

Autoscale Controller 简化了扩缩容策略的配置流程，并通过自动维护 HPA 和 VPA 的生命周期，提升了集群资源利用率和应用稳定性。


## 功能特性

- 支持自动创建和管理HPA和VPA
- 支持通过注解配置HPA和VPA参数
- 支持多种工作负载类型（Deployment、StatefulSet、DaemonSet）
- 支持CPU和内存的自动扩缩容
- 支持多种更新模式（Auto、Initial、Off）
- 支持资源策略（json格式）

## 使用方法

### 安装

```bash
kubectl apply -f https://raw.githubusercontent.com/infraflows/autoscaler-controller/main/dist/install.yaml
```

### 配置示例

#### HPA配置

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: example
  annotations:
    hpa.infraflow.co/min-replicas: "1"
    hpa.infraflow.co/max-replicas: "10"
    cpu.hpa.infraflow.co/target-average-utilization: "80"
    memory.hpa.infraflow.co/target-average-utilization: "70"
spec:
  # ... 其他配置 ...
```

#### VPA配置

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: example
  annotations:
    vpa.infraflow.co/update-mode: "Auto"
    vpa.infraflow.co/resource-policy: |
      {
        "containerPolicies": [
          {
            "containerName": "*",
            "minAllowed": {
              "cpu": "100m",
              "memory": "100Mi"
            },
            "maxAllowed": {
              "cpu": "1",
              "memory": "1Gi"
            }
          }
        ]
      }
spec:
  # ... 其他配置 ...
```
更多配置示例请参考[示例配置](config/samples/)

## 支持的注解

[Annotations](docs/annotations.md)

## 许可证

Apache License 2.0

