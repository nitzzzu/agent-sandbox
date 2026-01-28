/*
 * Copyright 2025 The https://github.com/agent-sandbox/agent-sandbox Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package api

import (
	"encoding/json"
	"time"
)

// Defines values for NodeStatus.
const (
	NodeStatusConnecting NodeStatus = "connecting"
	NodeStatusDraining   NodeStatus = "draining"
	NodeStatusReady      NodeStatus = "ready"
	NodeStatusUnhealthy  NodeStatus = "unhealthy"
)

// Defines values for SandboxState.
const (
	Paused   SandboxState = "paused"
	Running  SandboxState = "running"
	Creating SandboxState = "creating"
	Ready    SandboxState = "ready"
	Unready  SandboxState = "unready"
)

// Defines values for TemplateBuildStatus.
const (
	TemplateBuildStatusBuilding TemplateBuildStatus = "building"
	TemplateBuildStatusError    TemplateBuildStatus = "error"
	TemplateBuildStatusReady    TemplateBuildStatus = "ready"
	TemplateBuildStatusWaiting  TemplateBuildStatus = "waiting"
)

// Defines values for GetTeamsTeamIDMetricsMaxParamsMetric.
const (
	ConcurrentSandboxes GetTeamsTeamIDMetricsMaxParamsMetric = "concurrent_sandboxes"
	SandboxStartRate    GetTeamsTeamIDMetricsMaxParamsMetric = "sandbox_start_rate"
)

// AWSRegistry defines model for AWSRegistry.
type AWSRegistry struct {
	// AwsAccessKeyId AWS Access Key ID for ECR authentication
	AwsAccessKeyId string `json:"awsAccessKeyId"`

	// AwsRegion AWS Region where the ECR registry is located
	AwsRegion string `json:"awsRegion"`

	// AwsSecretAccessKey AWS Secret Access Key for ECR authentication
	AwsSecretAccessKey string `json:"awsSecretAccessKey"`

	// Type Type of registry authentication
	Type AWSRegistryType `json:"type"`
}

// AWSRegistryType Type of registry authentication
type AWSRegistryType string

// AdminSandboxKillResult defines model for AdminSandboxKillResult.
type AdminSandboxKillResult struct {
	// FailedCount Number of sandboxes that failed to kill
	FailedCount int `json:"failedCount"`

	// KilledCount Number of sandboxes successfully killed
	KilledCount int `json:"killedCount"`
}

// AssignTemplateTagRequest defines model for AssignTemplateTagRequest.
type AssignTemplateTagRequest struct {
	// Names Names of the template
	Names []string `json:"names"`

	// Target Target template name in "alias:tag" format
	Target string `json:"target"`
}

// BuildLogEntry defines model for BuildLogEntry.
type BuildLogEntry struct {
	// Level State of the sandbox
	Level LogLevel `json:"level"`

	// Message Log message content
	Message string `json:"message"`

	// Step Step in the build process related to the log entry
	Step *string `json:"step,omitempty"`

	// Timestamp Timestamp of the log entry
	Timestamp time.Time `json:"timestamp"`
}

// BuildStatusReason defines model for BuildStatusReason.
type BuildStatusReason struct {
	// LogEntries Log entries related to the status reason
	LogEntries *[]BuildLogEntry `json:"logEntries,omitempty"`

	// Message Message with the status reason, currently reporting only for error status
	Message string `json:"message"`

	// Step Step that failed
	Step *string `json:"step,omitempty"`
}

// CPUCount CPU cores for the sandbox
type CPUCount = int32

// ConnectSandbox defines model for ConnectSandbox.
type ConnectSandbox struct {
	// Timeout Timeout in seconds from the current time after which the sandbox should expire
	Timeout int32 `json:"timeout"`
}

// DiskMetrics defines model for DiskMetrics.
type DiskMetrics struct {
	// Device Device name
	Device string `json:"device"`

	// FilesystemType Filesystem type (e.g., ext4, xfs)
	FilesystemType string `json:"filesystemType"`

	// MountPoint Mount point of the disk
	MountPoint string `json:"mountPoint"`

	// TotalBytes Total space in bytes
	TotalBytes uint64 `json:"totalBytes"`

	// UsedBytes Used space in bytes
	UsedBytes uint64 `json:"usedBytes"`
}

// DiskSizeMB Disk size for the sandbox in MiB
type DiskSizeMB = int32

// EnvVars defines model for EnvVars.
type EnvVars map[string]string

// EnvdVersion Version of the envd running in the sandbox
type EnvdVersion = string

// Error defines model for Error.
type Error struct {
	// Code Error code
	Code int32 `json:"code"`

	// Message Error
	Message string `json:"message"`
}

// FromImageRegistry defines model for FromImageRegistry.
type FromImageRegistry struct {
	union json.RawMessage
}

// GCPRegistry defines model for GCPRegistry.
type GCPRegistry struct {
	// ServiceAccountJson Service Account JSON for GCP authentication
	ServiceAccountJson string `json:"serviceAccountJson"`

	// Type Type of registry authentication
	Type GCPRegistryType `json:"type"`
}

// GCPRegistryType Type of registry authentication
type GCPRegistryType string

// GeneralRegistry defines model for GeneralRegistry.
type GeneralRegistry struct {
	// Password Password to use for the registry
	Password string `json:"password"`

	// Type Type of registry authentication
	Type GeneralRegistryType `json:"type"`

	// Username Username to use for the registry
	Username string `json:"username"`
}

// GeneralRegistryType Type of registry authentication
type GeneralRegistryType string

// IdentifierMaskingDetails defines model for IdentifierMaskingDetails.
type IdentifierMaskingDetails struct {
	// MaskedValuePrefix Prefix used in masked version of the token or key
	MaskedValuePrefix string `json:"maskedValuePrefix"`

	// MaskedValueSuffix Suffix used in masked version of the token or key
	MaskedValueSuffix string `json:"maskedValueSuffix"`

	// Prefix Prefix that identifies the token or key type
	Prefix string `json:"prefix"`

	// ValueLength Length of the token or key
	ValueLength int `json:"valueLength"`
}

// ListedSandbox defines model for ListedSandbox.
type ListedSandbox struct {
	// Alias Alias of the template
	Alias *string `json:"alias,omitempty"`

	// ClientID Identifier of the client
	// Deprecated:
	ClientID string `json:"clientID"`

	// CpuCount CPU cores for the sandbox
	CpuCount CPUCount `json:"cpuCount"`

	// DiskSizeMB Disk size for the sandbox in MiB
	DiskSizeMB DiskSizeMB `json:"diskSizeMB"`

	// EndAt Time when the sandbox will expire
	EndAt time.Time `json:"endAt"`

	// EnvdVersion Version of the envd running in the sandbox
	EnvdVersion EnvdVersion `json:"envdVersion"`

	// MemoryMB Memory for the sandbox in MiB
	MemoryMB MemoryMB         `json:"memoryMB"`
	Metadata *SandboxMetadata `json:"metadata,omitempty"`

	// SandboxID Identifier of the sandbox
	SandboxID string `json:"sandboxID"`

	// StartedAt Time when the sandbox was started
	StartedAt time.Time `json:"startedAt"`

	// State State of the sandbox
	State SandboxState `json:"state"`

	// TemplateID Identifier of the template from which is the sandbox created
	TemplateID string `json:"templateID"`
}

// LogLevel State of the sandbox
type LogLevel string

// LogsDirection Direction of the logs that should be returned
type LogsDirection string

// LogsSource Source of the logs that should be returned
type LogsSource string

// MachineInfo defines model for MachineInfo.
type MachineInfo struct {
	// CpuArchitecture CPU architecture of the node
	CpuArchitecture string `json:"cpuArchitecture"`

	// CpuFamily CPU family of the node
	CpuFamily string `json:"cpuFamily"`

	// CpuModel CPU model of the node
	CpuModel string `json:"cpuModel"`

	// CpuModelName CPU model name of the node
	CpuModelName string `json:"cpuModelName"`
}

// MaxTeamMetric Team metric with timestamp
type MaxTeamMetric struct {
	// Timestamp Timestamp of the metric entry
	// Deprecated:
	Timestamp time.Time `json:"timestamp"`

	// TimestampUnix Timestamp of the metric entry in Unix time (seconds since epoch)
	TimestampUnix int64 `json:"timestampUnix"`

	// Value The maximum value of the requested metric in the given interval
	Value float32 `json:"value"`
}

// Mcp MCP configuration for the sandbox
type Mcp map[string]interface{}

// MemoryMB Memory for the sandbox in MiB
type MemoryMB = int32

// NewAccessToken defines model for NewAccessToken.
type NewAccessToken struct {
	// Name Name of the access token
	Name string `json:"name"`
}

// NewSandbox defines model for NewSandbox.
type NewSandbox struct {
	// AllowInternetAccess Allow sandbox to access the internet. When set to false, it behaves the same as specifying denyOut to 0.0.0.0/0 in the network config.
	AllowInternetAccess *bool `json:"allow_internet_access,omitempty"`

	// AutoPause Automatically pauses the sandbox after the timeout
	AutoPause *bool             `json:"autoPause,omitempty"`
	EnvVars   map[string]string `json:"envVars,omitempty"`

	// Mcp MCP configuration for the sandbox
	Mcp      *Mcp                  `json:"mcp"`
	Metadata map[string]string     `json:"metadata,omitempty"`
	Network  *SandboxNetworkConfig `json:"network,omitempty"`

	// Secure Secure all system communication with sandbox
	Secure *bool `json:"secure,omitempty"`

	// TemplateID Identifier of the required template
	TemplateID string `json:"templateID"`

	// Timeout Time to live for the sandbox in seconds.
	Timeout int `json:"timeout,omitempty"`
}

// NewTeamAPIKey defines model for NewTeamAPIKey.
type NewTeamAPIKey struct {
	// Name Name of the API key
	Name string `json:"name"`
}

// Node defines model for Node.
type Node struct {
	// ClusterID Identifier of the cluster
	ClusterID string `json:"clusterID"`

	// Commit Commit of the orchestrator
	Commit string `json:"commit"`

	// CreateFails Number of sandbox create fails
	CreateFails uint64 `json:"createFails"`

	// CreateSuccesses Number of sandbox create successes
	CreateSuccesses uint64 `json:"createSuccesses"`

	// Id Identifier of the node
	Id          string      `json:"id"`
	MachineInfo MachineInfo `json:"machineInfo"`

	// Metrics Node metrics
	Metrics NodeMetrics `json:"metrics"`

	// NodeID Identifier of the nomad node
	// Deprecated:
	NodeID string `json:"nodeID"`

	// SandboxCount Number of sandboxes running on the node
	SandboxCount uint32 `json:"sandboxCount"`

	// SandboxStartingCount Number of starting Sandboxes
	SandboxStartingCount int `json:"sandboxStartingCount"`

	// ServiceInstanceID Service instance identifier of the node
	ServiceInstanceID string `json:"serviceInstanceID"`

	// Status Status of the node
	Status NodeStatus `json:"status"`

	// Version Version of the orchestrator
	Version string `json:"version"`
}

// NodeDetail defines model for NodeDetail.
type NodeDetail struct {
	// CachedBuilds List of cached builds id on the node
	CachedBuilds []string `json:"cachedBuilds"`

	// ClusterID Identifier of the cluster
	ClusterID string `json:"clusterID"`

	// Commit Commit of the orchestrator
	Commit string `json:"commit"`

	// CreateFails Number of sandbox create fails
	CreateFails uint64 `json:"createFails"`

	// CreateSuccesses Number of sandbox create successes
	CreateSuccesses uint64 `json:"createSuccesses"`

	// Id Identifier of the node
	Id          string      `json:"id"`
	MachineInfo MachineInfo `json:"machineInfo"`

	// Metrics Node metrics
	Metrics NodeMetrics `json:"metrics"`

	// NodeID Identifier of the nomad node
	// Deprecated:
	NodeID string `json:"nodeID"`

	// Sandboxes List of sandboxes running on the node
	Sandboxes []ListedSandbox `json:"sandboxes"`

	// ServiceInstanceID Service instance identifier of the node
	ServiceInstanceID string `json:"serviceInstanceID"`

	// Status Status of the node
	Status NodeStatus `json:"status"`

	// Version Version of the orchestrator
	Version string `json:"version"`
}

// NodeMetrics Node metrics
type NodeMetrics struct {
	// AllocatedCPU Number of allocated CPU cores
	AllocatedCPU uint32 `json:"allocatedCPU"`

	// AllocatedMemoryBytes Amount of allocated memory in bytes
	AllocatedMemoryBytes uint64 `json:"allocatedMemoryBytes"`

	// CpuCount Total number of CPU cores on the node
	CpuCount uint32 `json:"cpuCount"`

	// CpuPercent Node CPU usage percentage
	CpuPercent uint32 `json:"cpuPercent"`

	// Disks Detailed metrics for each disk/mount point
	Disks []DiskMetrics `json:"disks"`

	// MemoryTotalBytes Total node memory in bytes
	MemoryTotalBytes uint64 `json:"memoryTotalBytes"`

	// MemoryUsedBytes Node memory used in bytes
	MemoryUsedBytes uint64 `json:"memoryUsedBytes"`
}

// NodeStatus Status of the node
type NodeStatus string

// ResumedSandbox defines model for ResumedSandbox.
type ResumedSandbox struct {
	// AutoPause Automatically pauses the sandbox after the timeout
	// Deprecated:
	AutoPause *bool `json:"autoPause,omitempty"`

	// Timeout Time to live for the sandbox in seconds.
	Timeout *int32 `json:"timeout,omitempty"`
}

// Sandbox defines model for Sandbox.
type Sandbox struct {
	// Alias Alias of the template
	Alias *string `json:"alias,omitempty"`

	// ClientID Identifier of the client
	// Deprecated:
	ClientID string `json:"clientID"`

	// Domain Base domain where the sandbox traffic is accessible
	Domain string `json:"domain,omitempty"`

	// EnvdAccessToken Access token used for envd communication
	EnvdAccessToken string `json:"envdAccessToken,omitempty"`

	// EnvdVersion Version of the envd running in the sandbox
	EnvdVersion EnvdVersion `json:"envdVersion,omitempty"`

	// SandboxID Identifier of the sandbox
	SandboxID string `json:"sandboxID"`

	// TemplateID Identifier of the template from which is the sandbox created
	TemplateID string `json:"templateID,omitempty"`

	// TrafficAccessToken Token required for accessing sandbox via proxy.
	TrafficAccessToken string `json:"trafficAccessToken,omitempty"`

	Metadata map[string]string `json:"metadata,omitempty"`

	// CpuCount CPU cores for the sandbox
	CpuCount int64 `json:"cpuCount"`

	// DiskSizeMB Disk size for the sandbox in MiB
	DiskSizeMB int64 `json:"diskSizeMB"`

	// EndAt Time when the sandbox will expire
	EndAt time.Time `json:"endAt"`

	// MemoryMB Memory for the sandbox in MiB
	MemoryMB int64 `json:"memoryMB"`

	// StartedAt Time when the sandbox was started
	StartedAt time.Time `json:"startedAt"`

	// State State of the sandbox
	State SandboxState `json:"state"`
}

// SandboxLog Log entry with timestamp and line
type SandboxLog struct {
	// Line Log line content
	Line string `json:"line"`

	// Timestamp Timestamp of the log entry
	Timestamp time.Time `json:"timestamp"`
}

// SandboxLogEntry defines model for SandboxLogEntry.
type SandboxLogEntry struct {
	Fields map[string]string `json:"fields"`

	// Level State of the sandbox
	Level LogLevel `json:"level"`

	// Message Log message content
	Message string `json:"message"`

	// Timestamp Timestamp of the log entry
	Timestamp time.Time `json:"timestamp"`
}

// SandboxLogs defines model for SandboxLogs.
type SandboxLogs struct {
	// LogEntries Structured logs of the sandbox
	LogEntries []SandboxLogEntry `json:"logEntries"`

	// Logs Logs of the sandbox
	Logs []SandboxLog `json:"logs"`
}

// SandboxMetadata defines model for SandboxMetadata.
type SandboxMetadata map[string]string

// SandboxMetric Metric entry with timestamp and line
type SandboxMetric struct {
	// CpuCount Number of CPU cores
	CpuCount int32 `json:"cpuCount"`

	// CpuUsedPct CPU usage percentage
	CpuUsedPct float32 `json:"cpuUsedPct"`

	// DiskTotal Total disk space in bytes
	DiskTotal int64 `json:"diskTotal"`

	// DiskUsed Disk used in bytes
	DiskUsed int64 `json:"diskUsed"`

	// MemTotal Total memory in bytes
	MemTotal int64 `json:"memTotal"`

	// MemUsed Memory used in bytes
	MemUsed int64 `json:"memUsed"`

	// Timestamp Timestamp of the metric entry
	// Deprecated:
	Timestamp time.Time `json:"timestamp"`

	// TimestampUnix Timestamp of the metric entry in Unix time (seconds since epoch)
	TimestampUnix int64 `json:"timestampUnix"`
}

// SandboxNetworkConfig defines model for SandboxNetworkConfig.
type SandboxNetworkConfig struct {
	// AllowOut List of allowed CIDR blocks or IP addresses for egress traffic. Allowed addresses always take precedence over blocked addresses.
	AllowOut *[]string `json:"allowOut,omitempty"`

	// AllowPublicTraffic Specify if the sandbox URLs should be accessible only with authentication.
	AllowPublicTraffic *bool `json:"allowPublicTraffic,omitempty"`

	// DenyOut List of denied CIDR blocks or IP addresses for egress traffic
	DenyOut *[]string `json:"denyOut,omitempty"`

	// MaskRequestHost Specify host mask which will be used for all sandbox requests
	MaskRequestHost *string `json:"maskRequestHost,omitempty"`
}

// SandboxState State of the sandbox
type SandboxState string

// SandboxesWithMetrics defines model for SandboxesWithMetrics.
type SandboxesWithMetrics struct {
	Sandboxes map[string]SandboxMetric `json:"sandboxes"`
}

// Team defines model for Team.
type Team struct {
	// ApiKey API key for the team
	ApiKey string `json:"apiKey"`

	// IsDefault Whether the team is the default team
	IsDefault bool `json:"isDefault"`

	// Name Name of the team
	Name string `json:"name"`

	// TeamID Identifier of the team
	TeamID string `json:"teamID"`
}

// TeamMetric Team metric with timestamp
type TeamMetric struct {
	// ConcurrentSandboxes The number of concurrent sandboxes for the team
	ConcurrentSandboxes int32 `json:"concurrentSandboxes"`

	// SandboxStartRate Number of sandboxes started per second
	SandboxStartRate float32 `json:"sandboxStartRate"`

	// Timestamp Timestamp of the metric entry
	// Deprecated:
	Timestamp time.Time `json:"timestamp"`

	// TimestampUnix Timestamp of the metric entry in Unix time (seconds since epoch)
	TimestampUnix int64 `json:"timestampUnix"`
}

// Template defines model for Template.
type Template struct {
	// Aliases Aliases of the template
	Aliases []string `json:"aliases"`

	// BuildCount Number of times the template was built
	BuildCount int32 `json:"buildCount"`

	// BuildID Identifier of the last successful build for given template
	BuildID string `json:"buildID"`

	// BuildStatus Status of the template build
	BuildStatus TemplateBuildStatus `json:"buildStatus"`

	// CpuCount CPU cores for the sandbox
	CpuCount CPUCount `json:"cpuCount"`

	// CreatedAt Time when the template was created
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy *string   `json:"createdBy"`

	// DiskSizeMB Disk size for the sandbox in MiB
	DiskSizeMB DiskSizeMB `json:"diskSizeMB"`

	// EnvdVersion Version of the envd running in the sandbox
	EnvdVersion EnvdVersion `json:"envdVersion"`

	// LastSpawnedAt Time when the template was last used
	LastSpawnedAt *time.Time `json:"lastSpawnedAt"`

	// MemoryMB Memory for the sandbox in MiB
	MemoryMB MemoryMB `json:"memoryMB"`

	// Public Whether the template is public or only accessible by the team
	Public bool `json:"public"`

	// SpawnCount Number of times the template was used
	SpawnCount int64 `json:"spawnCount"`

	// TemplateID Identifier of the template
	TemplateID string `json:"templateID"`

	// UpdatedAt Time when the template was last updated
	UpdatedAt time.Time `json:"updatedAt"`
}

// TemplateAliasResponse defines model for TemplateAliasResponse.
type TemplateAliasResponse struct {
	// Public Whether the template is public or only accessible by the team
	Public bool `json:"public"`

	// TemplateID Identifier of the template
	TemplateID string `json:"templateID"`
}

// TemplateBuild defines model for TemplateBuild.
type TemplateBuild struct {
	// BuildID Identifier of the build
	BuildID string `json:"buildID"`

	// CpuCount CPU cores for the sandbox
	CpuCount CPUCount `json:"cpuCount"`

	// CreatedAt Time when the build was created
	CreatedAt time.Time `json:"createdAt"`

	// DiskSizeMB Disk size for the sandbox in MiB
	DiskSizeMB *DiskSizeMB `json:"diskSizeMB,omitempty"`

	// EnvdVersion Version of the envd running in the sandbox
	EnvdVersion *EnvdVersion `json:"envdVersion,omitempty"`

	// FinishedAt Time when the build was finished
	FinishedAt *time.Time `json:"finishedAt,omitempty"`

	// MemoryMB Memory for the sandbox in MiB
	MemoryMB MemoryMB `json:"memoryMB"`

	// Status Status of the template build
	Status TemplateBuildStatus `json:"status"`

	// UpdatedAt Time when the build was last updated
	UpdatedAt time.Time `json:"updatedAt"`
}

// TemplateBuildFileUpload defines model for TemplateBuildFileUpload.
type TemplateBuildFileUpload struct {
	// Present Whether the file is already present in the cache
	Present bool `json:"present"`

	// Url Url where the file should be uploaded to
	Url *string `json:"url,omitempty"`
}

// TemplateBuildInfo defines model for TemplateBuildInfo.
type TemplateBuildInfo struct {
	// BuildID Identifier of the build
	BuildID string `json:"buildID"`

	// LogEntries Build logs structured
	LogEntries []BuildLogEntry `json:"logEntries"`

	// Logs Build logs
	Logs   []string           `json:"logs"`
	Reason *BuildStatusReason `json:"reason,omitempty"`

	// Status Status of the template build
	Status TemplateBuildStatus `json:"status"`

	// TemplateID Identifier of the template
	TemplateID string `json:"templateID"`
}

// TemplateBuildLogsResponse defines model for TemplateBuildLogsResponse.
type TemplateBuildLogsResponse struct {
	// Logs Build logs structured
	Logs []BuildLogEntry `json:"logs"`
}

// TemplateBuildRequest defines model for TemplateBuildRequest.
type TemplateBuildRequest struct {
	// Alias Alias of the template
	Alias *string `json:"alias,omitempty"`

	// CpuCount CPU cores for the sandbox
	CpuCount *CPUCount `json:"cpuCount,omitempty"`

	// Dockerfile Dockerfile for the template
	Dockerfile string `json:"dockerfile"`

	// MemoryMB Memory for the sandbox in MiB
	MemoryMB *MemoryMB `json:"memoryMB,omitempty"`

	// ReadyCmd Ready check command to execute in the template after the build
	ReadyCmd *string `json:"readyCmd,omitempty"`

	// StartCmd Start command to execute in the template after the build
	StartCmd *string `json:"startCmd,omitempty"`

	// TeamID Identifier of the team
	TeamID *string `json:"teamID,omitempty"`
}

// TemplateBuildRequestV2 defines model for TemplateBuildRequestV2.
type TemplateBuildRequestV2 struct {
	// Alias Alias of the template
	Alias string `json:"alias"`

	// CpuCount CPU cores for the sandbox
	CpuCount *CPUCount `json:"cpuCount,omitempty"`

	// MemoryMB Memory for the sandbox in MiB
	MemoryMB *MemoryMB `json:"memoryMB,omitempty"`

	// TeamID Identifier of the team
	// Deprecated:
	TeamID *string `json:"teamID,omitempty"`
}

// TemplateBuildRequestV3 defines model for TemplateBuildRequestV3.
type TemplateBuildRequestV3 struct {
	// Alias Alias of the template. Deprecated, use names instead.
	// Deprecated:
	Alias *string `json:"alias,omitempty"`

	// CpuCount CPU cores for the sandbox
	CpuCount *CPUCount `json:"cpuCount,omitempty"`

	// MemoryMB Memory for the sandbox in MiB
	MemoryMB *MemoryMB `json:"memoryMB,omitempty"`

	// Names Names of the template
	Names *[]string `json:"names,omitempty"`

	// TeamID Identifier of the team
	// Deprecated:
	TeamID *string `json:"teamID,omitempty"`
}

// TemplateBuildStartV2 defines model for TemplateBuildStartV2.
type TemplateBuildStartV2 struct {
	// Force Whether the whole build should be forced to run regardless of the cache
	Force *bool `json:"force,omitempty"`

	// FromImage Image to use as a base for the template build
	FromImage         *string            `json:"fromImage,omitempty"`
	FromImageRegistry *FromImageRegistry `json:"fromImageRegistry,omitempty"`

	// FromTemplate Template to use as a base for the template build
	FromTemplate *string `json:"fromTemplate,omitempty"`

	// ReadyCmd Ready check command to execute in the template after the build
	ReadyCmd *string `json:"readyCmd,omitempty"`

	// StartCmd Start command to execute in the template after the build
	StartCmd *string `json:"startCmd,omitempty"`

	// Steps List of steps to execute in the template build
	Steps *[]TemplateStep `json:"steps,omitempty"`
}

// TemplateBuildStatus Status of the template build
type TemplateBuildStatus string

// TemplateLegacy defines model for TemplateLegacy.
type TemplateLegacy struct {
	// Aliases Aliases of the template
	Aliases []string `json:"aliases"`

	// BuildCount Number of times the template was built
	BuildCount int32 `json:"buildCount"`

	// BuildID Identifier of the last successful build for given template
	BuildID string `json:"buildID"`

	// CpuCount CPU cores for the sandbox
	CpuCount CPUCount `json:"cpuCount"`

	// CreatedAt Time when the template was created
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy *string   `json:"createdBy"`

	// DiskSizeMB Disk size for the sandbox in MiB
	DiskSizeMB DiskSizeMB `json:"diskSizeMB"`

	// EnvdVersion Version of the envd running in the sandbox
	EnvdVersion EnvdVersion `json:"envdVersion"`

	// LastSpawnedAt Time when the template was last used
	LastSpawnedAt *time.Time `json:"lastSpawnedAt"`

	// MemoryMB Memory for the sandbox in MiB
	MemoryMB MemoryMB `json:"memoryMB"`

	// Public Whether the template is public or only accessible by the team
	Public bool `json:"public"`

	// SpawnCount Number of times the template was used
	SpawnCount int64 `json:"spawnCount"`

	// TemplateID Identifier of the template
	TemplateID string `json:"templateID"`

	// UpdatedAt Time when the template was last updated
	UpdatedAt time.Time `json:"updatedAt"`
}

// TemplateRequestResponseV3 defines model for TemplateRequestResponseV3.
type TemplateRequestResponseV3 struct {
	// Aliases Aliases of the template
	Aliases []string `json:"aliases"`

	// BuildID Identifier of the last successful build for given template
	BuildID string `json:"buildID"`

	// Public Whether the template is public or only accessible by the team
	Public bool `json:"public"`

	// TemplateID Identifier of the template
	TemplateID string `json:"templateID"`
}

// TemplateStep Step in the template build process
type TemplateStep struct {
	// Args Arguments for the step
	Args *[]string `json:"args,omitempty"`

	// FilesHash Hash of the files used in the step
	FilesHash *string `json:"filesHash,omitempty"`

	// Force Whether the step should be forced to run regardless of the cache
	Force *bool `json:"force,omitempty"`

	// Type Type of the step
	Type string `json:"type"`
}

// TemplateTag defines model for TemplateTag.
type TemplateTag struct {
	// BuildID Identifier of the build associated with this tag
	BuildID string `json:"buildID"`

	// Names Assigned names of the template
	Names []string `json:"names"`
}

// TemplateUpdateRequest defines model for TemplateUpdateRequest.
type TemplateUpdateRequest struct {
	// Public Whether the template is public or only accessible by the team
	Public *bool `json:"public,omitempty"`
}

// TemplateWithBuilds defines model for TemplateWithBuilds.
type TemplateWithBuilds struct {
	// Aliases Aliases of the template
	Aliases []string `json:"aliases"`

	// Builds List of builds for the template
	Builds []TemplateBuild `json:"builds"`

	// CreatedAt Time when the template was created
	CreatedAt time.Time `json:"createdAt"`

	// LastSpawnedAt Time when the template was last used
	LastSpawnedAt *time.Time `json:"lastSpawnedAt"`

	// Public Whether the template is public or only accessible by the team
	Public bool `json:"public"`

	// SpawnCount Number of times the template was used
	SpawnCount int64 `json:"spawnCount"`

	// TemplateID Identifier of the template
	TemplateID string `json:"templateID"`

	// UpdatedAt Time when the template was last updated
	UpdatedAt time.Time `json:"updatedAt"`
}

// UpdateTeamAPIKey defines model for UpdateTeamAPIKey.
type UpdateTeamAPIKey struct {
	// Name New name for the API key
	Name string `json:"name"`
}

// AccessTokenID defines model for accessTokenID.
type AccessTokenID = string

// ApiKeyID defines model for apiKeyID.
type ApiKeyID = string

// BuildID defines model for buildID.
type BuildID = string

// NodeID defines model for nodeID.
type NodeID = string

// PaginationLimit defines model for paginationLimit.
type PaginationLimit = int32

// PaginationNextToken defines model for paginationNextToken.
type PaginationNextToken = string

// SandboxID defines model for sandboxID.
type SandboxID = string

// TeamID defines model for teamID.
type TeamID = string

// TemplateID defines model for templateID.
type TemplateID = string

// N400 defines model for 400.
type N400 = Error

// N401 defines model for 401.
type N401 = Error

// N403 defines model for 403.
type N403 = Error

// N404 defines model for 404.
type N404 = Error

// N409 defines model for 409.
type N409 = Error

// N500 defines model for 500.
type N500 = Error

// GetNodesNodeIDParams defines parameters for GetNodesNodeID.
type GetNodesNodeIDParams struct {
	// ClusterID Identifier of the cluster
	ClusterID *string `form:"clusterID,omitempty" json:"clusterID,omitempty"`
}

// GetSandboxesParams defines parameters for GetSandboxes.
type GetSandboxesParams struct {
	// Metadata Metadata query used to filter the sandboxes (e.g. "user=abc&app=prod"). Each key and values must be URL encoded.
	Metadata *string `form:"metadata,omitempty" json:"metadata,omitempty"`
}

// GetSandboxesMetricsParams defines parameters for GetSandboxesMetrics.
type GetSandboxesMetricsParams struct {
	// SandboxIds Comma-separated list of sandbox IDs to get metrics for
	SandboxIds []string `form:"sandbox_ids" json:"sandbox_ids"`
}

// GetSandboxesSandboxIDLogsParams defines parameters for GetSandboxesSandboxIDLogs.
type GetSandboxesSandboxIDLogsParams struct {
	// Start Starting timestamp of the logs that should be returned in milliseconds
	Start *int64 `form:"start,omitempty" json:"start,omitempty"`

	// Limit Maximum number of logs that should be returned
	Limit *int32 `form:"limit,omitempty" json:"limit,omitempty"`
}

// GetSandboxesSandboxIDMetricsParams defines parameters for GetSandboxesSandboxIDMetrics.
type GetSandboxesSandboxIDMetricsParams struct {
	// Start Unix timestamp for the start of the interval, in seconds, for which the metrics
	Start *int64 `form:"start,omitempty" json:"start,omitempty"`
	End   *int64 `form:"end,omitempty" json:"end,omitempty"`
}

// PostSandboxesSandboxIDRefreshesJSONBody defines parameters for PostSandboxesSandboxIDRefreshes.
type PostSandboxesSandboxIDRefreshesJSONBody struct {
	// Duration Duration for which the sandbox should be kept alive in seconds
	Duration *int `json:"duration,omitempty"`
}

// PostSandboxesSandboxIDTimeoutJSONBody defines parameters for PostSandboxesSandboxIDTimeout.
type PostSandboxesSandboxIDTimeoutJSONBody struct {
	// Timeout Timeout in seconds from the current time after which the sandbox should expire
	Timeout int32 `json:"timeout"`
}

// GetTeamsTeamIDMetricsParams defines parameters for GetTeamsTeamIDMetrics.
type GetTeamsTeamIDMetricsParams struct {
	// Start Unix timestamp for the start of the interval, in seconds, for which the metrics
	Start *int64 `form:"start,omitempty" json:"start,omitempty"`
	End   *int64 `form:"end,omitempty" json:"end,omitempty"`
}

// GetTeamsTeamIDMetricsMaxParams defines parameters for GetTeamsTeamIDMetricsMax.
type GetTeamsTeamIDMetricsMaxParams struct {
	// Start Unix timestamp for the start of the interval, in seconds, for which the metrics
	Start *int64 `form:"start,omitempty" json:"start,omitempty"`
	End   *int64 `form:"end,omitempty" json:"end,omitempty"`

	// Metric Metric to retrieve the maximum value for
	Metric GetTeamsTeamIDMetricsMaxParamsMetric `form:"metric" json:"metric"`
}

// GetTeamsTeamIDMetricsMaxParamsMetric defines parameters for GetTeamsTeamIDMetricsMax.
type GetTeamsTeamIDMetricsMaxParamsMetric string

// GetTemplatesParams defines parameters for GetTemplates.
type GetTemplatesParams struct {
	TeamID *string `form:"teamID,omitempty" json:"teamID,omitempty"`
}

// GetTemplatesTemplateIDParams defines parameters for GetTemplatesTemplateID.
type GetTemplatesTemplateIDParams struct {
	// NextToken Cursor to start the list from
	NextToken *PaginationNextToken `form:"nextToken,omitempty" json:"nextToken,omitempty"`

	// Limit Maximum number of items to return per page
	Limit *PaginationLimit `form:"limit,omitempty" json:"limit,omitempty"`
}

// GetTemplatesTemplateIDBuildsBuildIDLogsParams defines parameters for GetTemplatesTemplateIDBuildsBuildIDLogs.
type GetTemplatesTemplateIDBuildsBuildIDLogsParams struct {
	// Cursor Starting timestamp of the logs that should be returned in milliseconds
	Cursor *int64 `form:"cursor,omitempty" json:"cursor,omitempty"`

	// Limit Maximum number of logs that should be returned
	Limit     *int32         `form:"limit,omitempty" json:"limit,omitempty"`
	Direction *LogsDirection `form:"direction,omitempty" json:"direction,omitempty"`
	Level     *LogLevel      `form:"level,omitempty" json:"level,omitempty"`

	// Source Source of the logs that should be returned from
	Source *LogsSource `form:"source,omitempty" json:"source,omitempty"`
}

// GetTemplatesTemplateIDBuildsBuildIDStatusParams defines parameters for GetTemplatesTemplateIDBuildsBuildIDStatus.
type GetTemplatesTemplateIDBuildsBuildIDStatusParams struct {
	// LogsOffset Index of the starting build log that should be returned with the template
	LogsOffset *int32 `form:"logsOffset,omitempty" json:"logsOffset,omitempty"`

	// Limit Maximum number of logs that should be returned
	Limit *int32    `form:"limit,omitempty" json:"limit,omitempty"`
	Level *LogLevel `form:"level,omitempty" json:"level,omitempty"`
}

// GetV2SandboxesParams defines parameters for GetV2Sandboxes.
type GetV2SandboxesParams struct {
	// Metadata Metadata query used to filter the sandboxes (e.g. "user=abc&app=prod"). Each key and values must be URL encoded.
	Metadata *string `form:"metadata,omitempty" json:"metadata,omitempty"`

	// State Filter sandboxes by one or more states
	State *[]SandboxState `form:"state,omitempty" json:"state,omitempty"`

	// NextToken Cursor to start the list from
	NextToken *PaginationNextToken `form:"nextToken,omitempty" json:"nextToken,omitempty"`

	// Limit Maximum number of items to return per page
	Limit *PaginationLimit `form:"limit,omitempty" json:"limit,omitempty"`
}

// PostAccessTokensJSONRequestBody defines body for PostAccessTokens for application/json ContentType.
type PostAccessTokensJSONRequestBody = NewAccessToken

// PostApiKeysJSONRequestBody defines body for PostApiKeys for application/json ContentType.
type PostApiKeysJSONRequestBody = NewTeamAPIKey

// PatchApiKeysApiKeyIDJSONRequestBody defines body for PatchApiKeysApiKeyID for application/json ContentType.
type PatchApiKeysApiKeyIDJSONRequestBody = UpdateTeamAPIKey

// PostSandboxesJSONRequestBody defines body for PostSandboxes for application/json ContentType.
type PostSandboxesJSONRequestBody = NewSandbox

// PostSandboxesSandboxIDConnectJSONRequestBody defines body for PostSandboxesSandboxIDConnect for application/json ContentType.
type PostSandboxesSandboxIDConnectJSONRequestBody = ConnectSandbox

// PostSandboxesSandboxIDRefreshesJSONRequestBody defines body for PostSandboxesSandboxIDRefreshes for application/json ContentType.
type PostSandboxesSandboxIDRefreshesJSONRequestBody PostSandboxesSandboxIDRefreshesJSONBody

// PostSandboxesSandboxIDResumeJSONRequestBody defines body for PostSandboxesSandboxIDResume for application/json ContentType.
type PostSandboxesSandboxIDResumeJSONRequestBody = ResumedSandbox

// PostSandboxesSandboxIDTimeoutJSONRequestBody defines body for PostSandboxesSandboxIDTimeout for application/json ContentType.
type PostSandboxesSandboxIDTimeoutJSONRequestBody PostSandboxesSandboxIDTimeoutJSONBody

// PostTemplatesJSONRequestBody defines body for PostTemplates for application/json ContentType.
type PostTemplatesJSONRequestBody = TemplateBuildRequest

// PostTemplatesTagsJSONRequestBody defines body for PostTemplatesTags for application/json ContentType.
type PostTemplatesTagsJSONRequestBody = AssignTemplateTagRequest

// PatchTemplatesTemplateIDJSONRequestBody defines body for PatchTemplatesTemplateID for application/json ContentType.
type PatchTemplatesTemplateIDJSONRequestBody = TemplateUpdateRequest

// PostTemplatesTemplateIDJSONRequestBody defines body for PostTemplatesTemplateID for application/json ContentType.
type PostTemplatesTemplateIDJSONRequestBody = TemplateBuildRequest

// PostV2TemplatesJSONRequestBody defines body for PostV2Templates for application/json ContentType.
type PostV2TemplatesJSONRequestBody = TemplateBuildRequestV2

// PostV2TemplatesTemplateIDBuildsBuildIDJSONRequestBody defines body for PostV2TemplatesTemplateIDBuildsBuildID for application/json ContentType.
type PostV2TemplatesTemplateIDBuildsBuildIDJSONRequestBody = TemplateBuildStartV2

// PostV3TemplatesJSONRequestBody defines body for PostV3Templates for application/json ContentType.
type PostV3TemplatesJSONRequestBody = TemplateBuildRequestV3

// AsAWSRegistry returns the union data inside the FromImageRegistry as a AWSRegistry
func (t FromImageRegistry) AsAWSRegistry() (AWSRegistry, error) {
	var body AWSRegistry
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromAWSRegistry overwrites any union data inside the FromImageRegistry as the provided AWSRegistry
func (t *FromImageRegistry) FromAWSRegistry(v AWSRegistry) error {
	v.Type = "aws"
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// AsGCPRegistry returns the union data inside the FromImageRegistry as a GCPRegistry
func (t FromImageRegistry) AsGCPRegistry() (GCPRegistry, error) {
	var body GCPRegistry
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromGCPRegistry overwrites any union data inside the FromImageRegistry as the provided GCPRegistry
func (t *FromImageRegistry) FromGCPRegistry(v GCPRegistry) error {
	v.Type = "gcp"
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// AsGeneralRegistry returns the union data inside the FromImageRegistry as a GeneralRegistry
func (t FromImageRegistry) AsGeneralRegistry() (GeneralRegistry, error) {
	var body GeneralRegistry
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromGeneralRegistry overwrites any union data inside the FromImageRegistry as the provided GeneralRegistry
func (t *FromImageRegistry) FromGeneralRegistry(v GeneralRegistry) error {
	v.Type = "registry"
	b, err := json.Marshal(v)
	t.union = b
	return err
}

func (t FromImageRegistry) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *FromImageRegistry) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}
