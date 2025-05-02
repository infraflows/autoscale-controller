package consts

// Private prefixes for annotations.
const (
	hpaPrefix = "hpa.infraflow.co/"
	vpaPrefix = "vpa.infraflow.co/"
)

// HPAMinReplicas defines the minimum number of replicas for the workload.
// Value: string. Example: "2".
const HPAMinReplicas = hpaPrefix + "minReplicas"

// HPAMaxReplicas defines the maximum number of replicas for the workload.
// Value: string. Example: "10".
const HPAMaxReplicas = hpaPrefix + "maxReplicas"

// HPACpuTargetAverageUtilization defines the target average CPU utilization (percentage) for HPA scaling.
// Value: string (percentage). Example: "80".
const HPACpuTargetAverageUtilization = hpaPrefix + "cpu.targetAverageUtilization"

// HPACpuTargetAverageValue defines the target average CPU consumption (cores) for HPA scaling.
// Value: string (CPU quantity). Example: "500m" (= 0.5 cores).
const HPACpuTargetAverageValue = hpaPrefix + "cpu.targetAverageValue"

// HPAMemoryTargetAverageUtilization defines the target average memory utilization (percentage) for HPA scaling.
// Value: string (percentage). Example: "75".
const HPAMemoryTargetAverageUtilization = hpaPrefix + "memory.targetAverageUtilization"

// HPAMemoryTargetAverageValue defines the target average memory consumption (bytes) for HPA scaling.
// Value: string (memory size). Example: "512Mi".
const HPAMemoryTargetAverageValue = hpaPrefix + "memory.targetAverageValue"

// VPACpuMinAllowed defines the minimum allowed CPU (cores) for a container in VPA recommendations.
// Value: string (CPU quantity). Example: "200m".
const VPACpuMinAllowed = vpaPrefix + "cpu.minAllowed"

// VPACpuMaxAllowed defines the maximum allowed CPU (cores) for a container in VPA recommendations.
// Value: string (CPU quantity). Example: "2".
const VPACpuMaxAllowed = vpaPrefix + "cpu.maxAllowed"

// VPAMemoryMinAllowed defines the minimum allowed memory (bytes) for a container in VPA recommendations.
// Value: string (memory size). Example: "256Mi".
const VPAMemoryMinAllowed = vpaPrefix + "memory.minAllowed"

// VPAMemoryMaxAllowed defines the maximum allowed memory (bytes) for a container in VPA recommendations.
// Value: string (memory size). Example: "4Gi".
const VPAMemoryMaxAllowed = vpaPrefix + "memory.maxAllowed"

// VPAUpdateMode defines the update mode for VPA (e.g., Auto, Off, Initial).
// Value: string. Allowed values: "Auto", "Off", "Initial".
const VPAUpdateMode = vpaPrefix + "updateMode"

// VPAResourcePolicy defines the VPA resource policy configuration for a container or workload.
// Value: string (JSON-encoded resource policy definition).
const VPAResourcePolicy = vpaPrefix + "resourcePolicy"

// VPAContainerPolicy defines container-specific resource policies under VPA configuration.
// Value: string (JSON-encoded container policies).
const VPAContainerPolicy = vpaPrefix + "containerPolicies"

const AutoScaleFinalizer = "finalizers.infraflow.co/autoscale"
