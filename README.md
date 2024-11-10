# ArgoCD Addon Controller Manager

This project extends ArgoCD functionality by providing two Custom Resource Definitions (CRDs) that enhance multi-cluster application deployment and cluster management.

## Overview

### ApplicationTemplate CRD

The `ApplicationTemplate` CRD enables hybrid rendering support for both Helm charts and Kustomize configurations. Key features include:

- Hybrid rendering support (Helm + Kustomize)
- Multiple target cluster deployment
- Environment-specific value configurations
- Cluster selection via labels or direct reference

Example:

```yaml
apiVersion: argocd-addon.github.com/v1alpha1
kind: ApplicationTemplate
metadata:
  name: sample-application-template
  namespace: application-templates
spec:
  name: app-chart                # required: name of the ApplicationTemplate
  repoURL: https://example.com/redis  # required: source repository URL
  targetRevision: main # option: default is main
  helm:                         # required: Helm configuration
    chart: "my-app-chart"       # required: chart name
    version: "1.0.0"           # required: chart version
    repository: "oci://registry-1.docker.io/bitnamicharts/redis"  # optional: Helm repository URL
    defaultValuesPath: "helm/values.yaml"  # optional: default values file path
    renderTargets:             # targets for rendering with different values
      - valuesPath: "helm/values-sit.yaml"  # values file path relative to source repo
        destinationCluster:    # required: cluster selector for deployment
          matchLabels:
            environment: sit
            region: cn-hangzhou
      - valuesPath: "helm/values-uat.yaml"
        destinationCluster:
          matchLabels:
            environment: uat
            region: cn-beijing
      - valuesPath: "helm/values-sit0.yaml"
        destinationCluster:
          name: kubernetes.default.svc-3396314289  # direct reference to ArgoCD cluster
  kustomize:                   # optional: Kustomize configuration
    renderTargets:            # targets for Kustomize overlays
      - path: "overlays/sit"  # overlay directory path
        destinationCluster:
          matchLabels:
            environment: sit
            region: cn-hangzhou
      - path: "overlays/uat"
        destinationCluster:
          matchLabels:
            environment: uat
            region: cn-beijing
status:
  phase: Succeeded            # Overall status
  matchedClusters:           # List of matched destination clusters
    - name: "kubernetes.default.svc-3396314289"
      matchedBy: "name"      # Matched by direct name reference
      rendered: true
    - name: "cluster-47.242.186.46-1493148463"
      matchedBy: "labels"    # Matched by label selector
      matchedLabels:
        environment: sit
        region: cn-hangzhou
      rendered: true
    - name: "cluster-47.242.187.46-1493148464"
      matchedBy: "labels"
      matchedLabels:
        environment: uat
        region: cn-beijing
      rendered: true
  renderedFiles:             # List of rendered manifest files
    - path: "/tmp/redis-sit-cluster-kubernetes.default.svc-3396314289.yaml"
      cluster: "kubernetes.default.svc-3396314289"
      type: "helm+kustomize"  # Rendered with both Helm and Kustomize
      timestamp: "2024-11-08T16:56:16+08:00"
    - path: "/tmp/redis-uat-cluster-47.242.186.46-1493148463.yaml"
      cluster: "cluster-47.242.186.46-1493148463"
      type: "helm+kustomize"
      timestamp: "2024-11-08T16:56:17+08:00"
    - path: "/tmp/redis-kubernetes.default.svc-3396314289.yaml"
      cluster: "kubernetes.default.svc-3396314289"
      type: "helm"           # Rendered with Helm only
      timestamp: "2024-11-08T16:56:18+08:00"
  conditions:                # Status conditions
    - type: Ready
      status: "True"
      lastUpdateTime: "2024-11-08T16:56:18+08:00"
      reason: "RenderingSucceeded"
      message: "Successfully rendered manifests for all matched clusters"
    - type: ClusterMatching
      status: "True"
      lastUpdateTime: "2024-11-08T16:56:16+08:00"
      reason: "ClustersMatched"
      message: "Successfully matched 3 destination clusters"
```

Key Features:
- **Metadata Management**:
  - Environment classification (sit, uat, prod)
  - Regional information
  - Cloud provider tracking
- **Secret Management**:
  - Secure storage of cluster credentials
  - Namespace-scoped secret references
- **Label-based Operations**:
  - Use standard Kubernetes labels for cluster selection
  - Enable environment and region-based targeting
- **ArgoCD Integration**:
  - Direct mapping to ArgoCD cluster configurations
  - Automatic sync with ArgoCD cluster states

## License

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