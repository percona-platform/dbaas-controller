// dbaas-controller
// Copyright (C) 2020 Percona LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

// Package v1 contains API Schema definitions for the pxc v1 API group
package v1

import (
	cmmeta "github.com/jetstack/cert-manager/pkg/apis/meta/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/percona-platform/dbaas-controller/k8_api/common"
)

// PerconaXtraDBClusterSpec defines the desired state of PerconaXtraDBCluster.
type PerconaXtraDBClusterSpec struct { //nolint:maligned
	Platform              string                               `json:"platform,omitempty"`
	CRVersion             string                               `json:"crVersion,omitempty"`
	Pause                 bool                                 `json:"pause,omitempty"`
	SecretsName           string                               `json:"secretsName,omitempty"`
	VaultSecretName       string                               `json:"vaultSecretName,omitempty"`
	SSLSecretName         string                               `json:"sslSecretName,omitempty"`
	SSLInternalSecretName string                               `json:"sslInternalSecretName,omitempty"`
	TLS                   *TLSSpec                             `json:"tls,omitempty"`
	PXC                   *PodSpec                             `json:"pxc,omitempty"`
	ProxySQL              *PodSpec                             `json:"proxysql,omitempty"`
	HAProxy               *PodSpec                             `json:"haproxy,omitempty"`
	PMM                   *PMMSpec                             `json:"pmm,omitempty"`
	Backup                *PXCScheduledBackup                  `json:"backup,omitempty"`
	UpdateStrategy        appsv1.StatefulSetUpdateStrategyType `json:"updateStrategy,omitempty"`
	UpgradeOptions        UpgradeOptions                       `json:"upgradeOptions,omitempty"`
	AllowUnsafeConfig     bool                                 `json:"allowUnsafeConfigurations,omitempty"`
}

// TLSSpec holds cluster's TLS specs.
type TLSSpec struct {
	SANs       []string                `json:"SANs,omitempty"`
	IssuerConf *cmmeta.ObjectReference `json:"issuerConf,omitempty"`
}

// UpgradeOptions holds configuration options to handle automatic upgrades.
type UpgradeOptions struct {
	VersionServiceEndpoint string `json:"versionServiceEndpoint,omitempty"`
	Apply                  string `json:"apply,omitempty"`
	Schedule               string `json:"schedule,omitempty"`
}

// PXCScheduledBackup holds the config for cluster scheduled backups.
type PXCScheduledBackup struct {
	Image              string                        `json:"image,omitempty"`
	ImagePullSecrets   []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`
	Schedule           []PXCScheduledBackupSchedule  `json:"schedule,omitempty"`
	Storages           map[string]*BackupStorageSpec `json:"storages,omitempty"`
	ServiceAccountName string                        `json:"serviceAccountName,omitempty"`
	Annotations        map[string]string             `json:"annotations,omitempty"`
}

// PXCScheduledBackupSchedule holds the backup schedule.
type PXCScheduledBackupSchedule struct {
	Name        string `json:"name,omitempty"`
	Schedule    string `json:"schedule,omitempty"`
	Keep        int    `json:"keep,omitempty"`
	StorageName string `json:"storageName,omitempty"`
}

// AppState Application state.
type AppState string

const (
	// AppStateUnknown application state is unknown.
	AppStateUnknown AppState = "unknown"
	// AppStateInit initializing state.
	AppStateInit AppState = "initializing"
	// AppStateReady application is ready state.
	AppStateReady AppState = "ready"
	// AppStateError error state.
	AppStateError AppState = "error"
)

// PerconaXtraDBClusterStatus defines the observed state of PerconaXtraDBCluster.
type PerconaXtraDBClusterStatus struct {
	PXC                AppStatus          `json:"pxc,omitempty"`
	ProxySQL           AppStatus          `json:"proxysql,omitempty"`
	HAProxy            AppStatus          `json:"haproxy,omitempty"`
	Backup             AppStatus          `json:"backup,omitempty"`
	PMM                AppStatus          `json:"pmm,omitempty"`
	Host               string             `json:"host,omitempty"`
	Messages           []string           `json:"message,omitempty"`
	Status             AppState           `json:"state,omitempty"`
	Conditions         []ClusterCondition `json:"conditions,omitempty"`
	ObservedGeneration int64              `json:"observedGeneration,omitempty"`
}

// ConditionStatus tells if the cluster condition can be determined.
type ConditionStatus string

// ClusterConditionType is the current condition state string.
type ClusterConditionType string

// ClusterCondition holds exported fields with the cluster condition.
type ClusterCondition struct {
	Status             ConditionStatus      `json:"status,omitempty"`
	Type               ClusterConditionType `json:"type,omitempty"`
	LastTransitionTime metav1.Time          `json:"lastTransitionTime,omitempty"`
	Reason             string               `json:"reason,omitempty"`
	Message            string               `json:"message,omitempty"`
}

// AppStatus holds exported fields representing the application status information.
type AppStatus struct {
	Size    int32    `json:"size,omitempty"`
	Ready   int32    `json:"ready,omitempty"`
	Status  AppState `json:"status,omitempty"`
	Message string   `json:"message,omitempty"`
	Version string   `json:"version,omitempty"`
	Image   string   `json:"image,omitempty"`
}

// PerconaXtraDBCluster is the Schema for the perconaxtradbclusters API.
type PerconaXtraDBCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PerconaXtraDBClusterSpec   `json:"spec,omitempty"`
	Status PerconaXtraDBClusterStatus `json:"status,omitempty"`
}

// PerconaXtraDBClusterList contains a list of PerconaXtraDBCluster.
type PerconaXtraDBClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PerconaXtraDBCluster `json:"items"`
}

// PodSpec hold pod's exported fields representing the pod configuration.
type PodSpec struct { //nolint:maligned
	Enabled                       bool                                    `json:"enabled,omitempty"`
	Size                          int32                                   `json:"size,omitempty"`
	Image                         string                                  `json:"image,omitempty"`
	Resources                     *common.PodResources                    `json:"resources,omitempty"`
	SidecarResources              *common.PodResources                    `json:"sidecarResources,omitempty"`
	VolumeSpec                    *common.VolumeSpec                      `json:"volumeSpec,omitempty"`
	Affinity                      *PodAffinity                            `json:"affinity,omitempty"`
	NodeSelector                  map[string]string                       `json:"nodeSelector,omitempty"`
	Tolerations                   []corev1.Toleration                     `json:"tolerations,omitempty"`
	PriorityClassName             string                                  `json:"priorityClassName,omitempty"`
	Annotations                   map[string]string                       `json:"annotations,omitempty"`
	Labels                        map[string]string                       `json:"labels,omitempty"`
	ImagePullSecrets              []corev1.LocalObjectReference           `json:"imagePullSecrets,omitempty"`
	Configuration                 string                                  `json:"configuration,omitempty"`
	PodDisruptionBudget           *PodDisruptionBudgetSpec                `json:"podDisruptionBudget,omitempty"`
	VaultSecretName               string                                  `json:"vaultSecretName,omitempty"`
	SSLSecretName                 string                                  `json:"sslSecretName,omitempty"`
	SSLInternalSecretName         string                                  `json:"sslInternalSecretName,omitempty"`
	TerminationGracePeriodSeconds *int64                                  `json:"gracePeriod,omitempty"`
	ForceUnsafeBootstrap          bool                                    `json:"forceUnsafeBootstrap,omitempty"`
	ServiceType                   corev1.ServiceType                      `json:"serviceType,omitempty"`
	ExternalTrafficPolicy         corev1.ServiceExternalTrafficPolicyType `json:"externalTrafficPolicy,omitempty"`
	LoadBalancerSourceRanges      []string                                `json:"loadBalancerSourceRanges,omitempty"`
	ServiceAnnotations            map[string]string                       `json:"serviceAnnotations,omitempty"`
	SchedulerName                 string                                  `json:"schedulerName,omitempty"`
	ReadinessInitialDelaySeconds  *int32                                  `json:"readinessDelaySec,omitempty"`
	LivenessInitialDelaySeconds   *int32                                  `json:"livenessDelaySec,omitempty"`
	PodSecurityContext            *corev1.PodSecurityContext              `json:"podSecurityContext,omitempty"`
	ContainerSecurityContext      *corev1.SecurityContext                 `json:"containerSecurityContext,omitempty"`
	ServiceAccountName            string                                  `json:"serviceAccountName,omitempty"`
	ImagePullPolicy               corev1.PullPolicy                       `json:"imagePullPolicy,omitempty"`
}

// PodDisruptionBudgetSpec POD disruption budget specs.
type PodDisruptionBudgetSpec struct {
	MinAvailable   *intstr.IntOrString `json:"minAvailable,omitempty"`
	MaxUnavailable *intstr.IntOrString `json:"maxUnavailable,omitempty"`
}

// PodAffinity POD's affinity.
type PodAffinity struct {
	TopologyKey *string          `json:"antiAffinityTopologyKey,omitempty"`
	Advanced    *corev1.Affinity `json:"advanced,omitempty"`
}

// PMMSpec hold exported fields representing PMM specs.
type PMMSpec struct {
	Enabled                  bool                    `json:"enabled,omitempty"`
	ServerHost               string                  `json:"serverHost,omitempty"`
	Image                    string                  `json:"image,omitempty"`
	ServerUser               string                  `json:"serverUser,omitempty"`
	Resources                *common.PodResources    `json:"resources,omitempty"`
	ContainerSecurityContext *corev1.SecurityContext `json:"containerSecurityContext,omitempty"`
	ImagePullPolicy          corev1.PullPolicy       `json:"imagePullPolicy,omitempty"`
}

// BackupStorageSpec holds backup's storage specs.
type BackupStorageSpec struct {
	Type                     BackupStorageType          `json:"type"`
	S3                       BackupStorageS3Spec        `json:"s3,omitempty"`
	Volume                   *common.VolumeSpec         `json:"volume,omitempty"`
	NodeSelector             map[string]string          `json:"nodeSelector,omitempty"`
	Resources                *common.PodResources       `json:"resources,omitempty"`
	Affinity                 *corev1.Affinity           `json:"affinity,omitempty"`
	Tolerations              []corev1.Toleration        `json:"tolerations,omitempty"`
	Annotations              map[string]string          `json:"annotations,omitempty"`
	Labels                   map[string]string          `json:"labels,omitempty"`
	SchedulerName            string                     `json:"schedulerName,omitempty"`
	PriorityClassName        string                     `json:"priorityClassName,omitempty"`
	PodSecurityContext       *corev1.PodSecurityContext `json:"podSecurityContext,omitempty"`
	ContainerSecurityContext *corev1.SecurityContext    `json:"containerSecurityContext,omitempty"`
}

// BackupStorageType backup storage type.
type BackupStorageType string

// nolint:unused
const (
	// BackupStorageFilesystem use local filesystem for storage.
	BackupStorageFilesystem BackupStorageType = "filesystem"
	// BackupStorageS3 use S3 for storage.
	BackupStorageS3 BackupStorageType = "s3"
)

// BackupStorageS3Spec holds the S3 configuration.
type BackupStorageS3Spec struct {
	Bucket            string `json:"bucket"`
	CredentialsSecret string `json:"credentialsSecret"`
	Region            string `json:"region,omitempty"`
	EndpointURL       string `json:"endpointUrl,omitempty"`
}

// Volume represents a backup volume.
type Volume struct { //nolint:unused
	PVCs    []corev1.PersistentVolumeClaim
	Volumes []corev1.Volume
}

// AffinityTopologyKeyOff Affinity Topology Key Off.
const AffinityTopologyKeyOff = "none"
