package consts

// Private prefixes for annotations.
const (
	hpaPrefix = "hpa.infraflow.co/"
	vpaPrefix = "vpa.infraflow.co/"
)

// HPAAnnotationMinReplicas defines the minimum number of replicas for the workload.
// Value: string. Example: "2".
const HPACpuMinReplicas = hpaPrefix + "cpu.minReplicas"

// HPAAnnotationMaxReplicas defines the maximum number of replicas for the workload.
// Value: string. Example: "10".
const HPACpuMaxReplicas = hpaPrefix + "cpu.maxReplicas"

// CPUHPAAnnotationTargetAverageUtilization defines the target average CPU utilization (percentage) for HPA scaling.
// Value: string (percentage). Example: "80".
const HPACpuTargetAverageUtilization = hpaPrefix + "cpu.targetAverageUtilization"

// CPUHPAAnnotationTargetAverageValue defines the target average CPU consumption (cores) for HPA scaling.
// Value: string (CPU quantity). Example: "500m" (= 0.5 cores).
const HPACpuTargetAverageValue = hpaPrefix + "cpu.targetAverageValue"

// MemoryHPAAnnotationTargetAverageUtilization defines the target average memory utilization (percentage) for HPA scaling.
// Value: string (percentage). Example: "75".
const HPAMemoryTargetAverageUtilization = hpaPrefix + "memory.targetAverageUtilization"

// MemoryHPAAnnotationTargetAverageValue defines the target average memory consumption (bytes) for HPA scaling.
// Value: string (memory size). Example: "512Mi".
const HPAMemoryTargetAverageValue = hpaPrefix + "memory.targetAverageValue"

// VPAAnnotationMinAllowedCPU defines the minimum allowed CPU (cores) for a container in VPA recommendations.
// Value: string (CPU quantity). Example: "200m".
const VPACpuMinAllowed = vpaPrefix + "cpu.minAllowed"

// VPAAnnotationMaxAllowedCPU defines the maximum allowed CPU (cores) for a container in VPA recommendations.
// Value: string (CPU quantity). Example: "2".
const VPACpuMaxAllowed = vpaPrefix + "cpu.maxAllowed"

// VPAAnnotationMinAllowedMemory defines the minimum allowed memory (bytes) for a container in VPA recommendations.
// Value: string (memory size). Example: "256Mi".
const VPAMemoryMinAllowed = vpaPrefix + "memory.minAllowed"

// VPAAnnotationMaxAllowedMemory defines the maximum allowed memory (bytes) for a container in VPA recommendations.
// Value: string (memory size). Example: "4Gi".
const VPAMemoryMaxAllowed = vpaPrefix + "memory.maxAllowed"

// VPAAnnotationUpdateMode defines the update mode for VPA (e.g., Auto, Off, Initial).
// Value: string. Allowed values: "Auto", "Off", "Initial".
const VPAMode = vpaPrefix + "mode"

// VPAAnnotationResourcePolicy defines the VPA resource policy configuration for a container or workload.
// Value: string (JSON-encoded resource policy definition).
const VPAResourcePolicy = vpaPrefix + "resourcePolicy"

// VPAAnnotationContainerPolicy defines container-specific resource policies under VPA configuration.
// Value: string (JSON-encoded container policies).
const VPAContainerPolicy = vpaPrefix + "containerPolicies"

const AutoScaleFinalizer = "finalizers.infraflow.co/autoscale"
