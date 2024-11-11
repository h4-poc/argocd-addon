/*
Copyright 2024.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ApplicationTemplateSpec defines the desired state of ApplicationTemplate
type ApplicationTemplateSpec struct {
	// Name is the name of the ApplicationTemplate
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// RepoURL is the source repository URL
	// +kubebuilder:validation:Required
	RepoURL string `json:"repoURL"`

	// TargetRevision is the git revision to use (branch, tag, commit)
	// +optional
	// +kubebuilder:default:=main
	TargetRevision string `json:"targetRevision,omitempty"`

	// Helm defines the Helm-specific template configuration
	// +optional
	Helm *HelmConfig `json:"helm,omitempty"`

	// Kustomize defines the Kustomize-specific template configuration
	// +optional
	Kustomize *KustomizeConfig `json:"kustomize,omitempty"`
}

// HelmConfig defines Helm-specific configuration
type HelmConfig struct {
	// Chart is the name of the Helm chart
	// +kubebuilder:validation:Required
	Chart string `json:"chart"`

	// Version is the version of the Helm chart
	// +kubebuilder:validation:Required
	Version string `json:"version"`

	// Repository is the Helm repository URL
	// +optional
	Repository string `json:"repository,omitempty"`

	// DefaultValuesPath is the path to the default values file
	// +optional
	DefaultValuesPath string `json:"defaultValuesPath,omitempty"`

	// RenderTargets defines the targets for rendering with different values
	// +optional
	RenderTargets []HelmRenderTarget `json:"renderTargets,omitempty"`
}

// HelmRenderTarget defines a target for rendering Helm charts
type HelmRenderTarget struct {
	// ValuesPath is the path to the values file relative to source repo
	// +kubebuilder:validation:Required
	ValuesPath string `json:"valuesPath"`

	// DestinationCluster defines the cluster selector for deployment
	// +kubebuilder:validation:Required
	DestinationCluster ClusterSelector `json:"destinationCluster"`
}

// KustomizeConfig defines Kustomize-specific configuration
type KustomizeConfig struct {
	// RenderTargets defines the targets for Kustomize overlays
	// +optional
	RenderTargets []KustomizeRenderTarget `json:"renderTargets,omitempty"`
}

// KustomizeRenderTarget defines a target for rendering Kustomize overlays
type KustomizeRenderTarget struct {
	// Path is the overlay directory path
	// +kubebuilder:validation:Required
	Path string `json:"path"`

	// DestinationCluster defines the cluster selector for deployment
	// +kubebuilder:validation:Required
	DestinationCluster ClusterSelector `json:"destinationCluster"`
}

// ClusterSelector defines how to select destination clusters
type ClusterSelector struct {
	// Name is a direct reference to an ArgoCD cluster
	// +optional
	Name string `json:"name,omitempty"`

	// MatchLabels is a map of labels to match clusters
	// +optional
	MatchLabels map[string]string `json:"matchLabels,omitempty"`
}

// ApplicationTemplateStatus defines the observed state of ApplicationTemplate
type ApplicationTemplateStatus struct {
	// Phase represents the overall status of the ApplicationTemplate
	// +optional
	Phase string `json:"phase,omitempty"`

	// MatchedClusters lists the clusters that matched the selectors
	// +optional
	MatchedClusters []MatchedCluster `json:"matchedClusters,omitempty"`

	// RenderedFiles lists the rendered manifest files
	// +optional
	RenderedFiles []RenderedFile `json:"renderedFiles,omitempty"`

	// Conditions represent the latest available observations of the ApplicationTemplate's state
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// MatchedCluster represents a cluster that matched the selectors
type MatchedCluster struct {
	// Name is the name of the matched cluster
	Name string `json:"name"`

	// MatchedBy indicates how the cluster was matched (name or labels)
	MatchedBy string `json:"matchedBy"`

	// MatchedLabels lists the labels that matched (if matched by labels)
	// +optional
	MatchedLabels map[string]string `json:"matchedLabels,omitempty"`

	// HelmValuesPath is the path to the Helm values file used for this cluster
	// +optional
	HelmValuesPath string `json:"helmValuesPath,omitempty"`

	// KustomizePath is the path to the Kustomize overlay used for this cluster
	// +optional
	KustomizePath string `json:"kustomizePath,omitempty"`

	// Rendered indicates if manifests were rendered for this cluster
	Rendered bool `json:"rendered"`
}

// RenderedFile represents a rendered manifest file
type RenderedFile struct {
	// Path is the path to the rendered file
	Path string `json:"path"`

	// Cluster is the name of the target cluster
	Cluster string `json:"cluster"`

	// Type indicates how the file was rendered (helm, kustomize, or both)
	Type string `json:"type"`

	// Timestamp indicates when the file was rendered
	Timestamp string `json:"timestamp"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=apptemp

// ApplicationTemplate is the Schema for the applicationtemplates API
type ApplicationTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApplicationTemplateSpec   `json:"spec,omitempty"`
	Status ApplicationTemplateStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ApplicationTemplateList contains a list of ApplicationTemplate
type ApplicationTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApplicationTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ApplicationTemplate{}, &ApplicationTemplateList{})
}
