package consts

// Private prefixes for annotations.
const (
	hpaPrefix = "hpa.infraflow.co/"
	vpaPrefix = "vpa.infraflow.co/"

	cpuHPAPrefix    = "cpu.hpa.infraflow.co/"
	memoryHPAPrefix = "memory.hpa.infraflow.co/"

	cpuVPAPrefix    = "cpu.vpa.infraflow.co/"
	memoryVPAPrefix = "memory.vpa.infraflow.co/"
)

// HPAAnnotationMinReplicas defines the minimum number of replicas for the workload.
// Value: string. Example: "2".
const HPAAnnotationMinReplicas = hpaPrefix + "minReplicas"

// HPAAnnotationMaxReplicas defines the maximum number of replicas for the workload.
// Value: string. Example: "10".
const HPAAnnotationMaxReplicas = hpaPrefix + "maxReplicas"

// CPUHPAAnnotationTargetAverageUtilization defines the target average CPU utilization (percentage) for HPA scaling.
// Value: string (percentage). Example: "80".
const CPUHPAAnnotationTargetAverageUtilization = cpuHPAPrefix + "targetAverageUtilization"

// CPUHPAAnnotationTargetAverageValue defines the target average CPU consumption (cores) for HPA scaling.
// Value: string (CPU quantity). Example: "500m" (= 0.5 cores).
const CPUHPAAnnotationTargetAverageValue = cpuHPAPrefix + "targetAverageValue"

// MemoryHPAAnnotationTargetAverageUtilization defines the target average memory utilization (percentage) for HPA scaling.
// Value: string (percentage). Example: "75".
const MemoryHPAAnnotationTargetAverageUtilization = memoryHPAPrefix + "targetAverageUtilization"

// MemoryHPAAnnotationTargetAverageValue defines the target average memory consumption (bytes) for HPA scaling.
// Value: string (memory size). Example: "512Mi".
const MemoryHPAAnnotationTargetAverageValue = memoryHPAPrefix + "targetAverageValue"

// VPAAnnotationMinAllowedCPU defines the minimum allowed CPU (cores) for a container in VPA recommendations.
// Value: string (CPU quantity). Example: "200m".
const VPAAnnotationMinAllowedCPU = cpuVPAPrefix + "minAllowed"

// VPAAnnotationMaxAllowedCPU defines the maximum allowed CPU (cores) for a container in VPA recommendations.
// Value: string (CPU quantity). Example: "2".
const VPAAnnotationMaxAllowedCPU = cpuVPAPrefix + "maxAllowed"

// VPAAnnotationMinAllowedMemory defines the minimum allowed memory (bytes) for a container in VPA recommendations.
// Value: string (memory size). Example: "256Mi".
const VPAAnnotationMinAllowedMemory = memoryVPAPrefix + "minAllowed"

// VPAAnnotationMaxAllowedMemory defines the maximum allowed memory (bytes) for a container in VPA recommendations.
// Value: string (memory size). Example: "4Gi".
const VPAAnnotationMaxAllowedMemory = memoryVPAPrefix + "maxAllowed"

// VPAAnnotationUpdateMode defines the update mode for VPA (e.g., Auto, Off, Initial).
// Value: string. Allowed values: "Auto", "Off", "Initial".
const VPAAnnotationUpdateMode = vpaPrefix + "updateMode"

// VPAAnnotationResourcePolicy defines the VPA resource policy configuration for a container or workload.
// Value: string (JSON-encoded resource policy definition).
const VPAAnnotationResourcePolicy = vpaPrefix + "resourcePolicy"

// VPAAnnotationContainerPolicy defines container-specific resource policies under VPA configuration.
// Value: string (JSON-encoded container policies).
const VPAAnnotationContainerPolicy = vpaPrefix + "containerPolicies"

const AutoScaleFinalizer = "finalizers.infraflow.co/autoscale"
