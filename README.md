# harbor

![Version: 3.0.0](https://img.shields.io/badge/Version-3.0.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 2.12.0](https://img.shields.io/badge/AppVersion-2.12.0-informational?style=flat-square)

A modern, simplified Helm chart for Harbor container registry

## TL;DR

```bash
helm repo add harbor-next https://example.com/charts
helm install my-harbor harbor-next/harbor \
  --set externalURL=https://harbor.example.com \
  --set database.host=my-postgres.example.com \
  --set database.password=secret
```

## Introduction

This is a modern, simplified Helm chart for [Harbor](https://goharbor.io/) - the cloud native registry for Kubernetes.

### Key Features

- **Future-proof configuration**: Uses `toEnvVars` pattern - any Harbor config option works without chart updates
- **Simplified architecture**: ~20 templates vs 48 in legacy chart
- **Production-ready**: External database only (use CloudNativePG, AWS RDS, etc.)
- **Flexible ingress**: Standard Ingress, Gateway API, or custom via `extraManifests`
- **Optional components**: Valkey (Redis) subchart, Trivy scanner subchart

### What's Different from Legacy Chart

| Legacy Chart | This Chart |
|-------------|-----------|
| Built-in PostgreSQL StatefulSet | External database only |
| Built-in Redis Deployment | Optional Valkey subchart or external |
| nginx reverse proxy | Direct Ingress/Gateway API |
| 5 expose types | Ingress, Gateway, or `extraManifests` |
| ~1,115 lines values.yaml | ~400 lines |
| 48 template files | ~20 templates |

## Prerequisites

- Kubernetes 1.26+
- Helm 3.x
- **External PostgreSQL database** (required)
- PV provisioner (for registry persistence)

## Installing the Chart

```bash
helm install my-harbor harbor-next/harbor \
  --namespace harbor \
  --create-namespace \
  --set externalURL=https://harbor.example.com \
  --set database.host=postgres.example.com \
  --set database.password=your-password \
  --set ingress.hosts[0].host=harbor.example.com
```

## Uninstalling the Chart

```bash
helm uninstall my-harbor --namespace harbor
```

## Requirements

Kubernetes: `>=1.26.0-0`

| Repository | Name | Version |
|------------|------|---------|
| https://aquasecurity.github.io/helm-charts/ | harbor-scanner-trivy | 0.x.x |
| https://valkey.io/valkey-helm/ | valkey | 0.9.x |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| core.affinity | object | `{}` | Affinity rules for Core pods |
| core.config | object | {} | Harbor Core application config (converted to env vars in ConfigMap) Any Harbor Core config can be set here without chart changes |
| core.extraEnv | list | [] | Extra environment variables with valueFrom support |
| core.image | object | `{"repository":"goharbor/harbor-core","tag":""}` | Core image settings |
| core.image.repository | string | `"goharbor/harbor-core"` | Core image repository |
| core.image.tag | string | `""` | Core image tag (defaults to appVersion) |
| core.nodeSelector | object | `{}` | Node selector for Core pods |
| core.podAnnotations | object | `{}` | Additional pod annotations for Core |
| core.podLabels | object | `{}` | Additional pod labels for Core |
| core.podSecurityContext | object | `{"fsGroup":10000}` | Pod security context for Core |
| core.replicas | int | `1` | Number of Core replicas |
| core.resources | object | `{}` | Core resource requests and limits |
| core.secret | object | {} | Sensitive config for Core (converted to env vars in Secret) |
| core.securityContext | object | `{"runAsGroup":10000,"runAsNonRoot":true,"runAsUser":10000}` | Security context for Core container |
| core.serviceAccount | object | `{"annotations":{},"create":true,"name":""}` | Service account settings for Core |
| core.serviceAccount.annotations | object | `{}` | Service account annotations |
| core.serviceAccount.create | bool | `true` | Create a service account for Core |
| core.serviceAccount.name | string | `""` | Service account name (auto-generated if empty) |
| core.tolerations | list | `[]` | Tolerations for Core pods |
| database.database | string | `"registry"` | Database name |
| database.existingSecret | string | `""` | Existing secret containing database credentials Must have key: POSTGRESQL_PASSWORD |
| database.host | string | `""` | Database host (required) |
| database.maxIdleConns | int | `100` | Maximum idle connections |
| database.maxOpenConns | int | `900` | Maximum open connections |
| database.password | string | `""` | Database password (ignored if existingSecret is set) |
| database.port | int | `5432` | Database port |
| database.sslmode | string | `"disable"` | SSL mode for database connection |
| database.username | string | `"postgres"` | Database username |
| exporter.affinity | object | `{}` | Affinity rules for Exporter pods |
| exporter.config | object | {} | Exporter application config (converted to env vars in ConfigMap) |
| exporter.enabled | bool | `true` | Enable Harbor exporter for Prometheus metrics |
| exporter.extraEnv | list | [] | Extra environment variables with valueFrom support |
| exporter.image | object | `{"repository":"goharbor/harbor-exporter","tag":""}` | Exporter image settings |
| exporter.image.repository | string | `"goharbor/harbor-exporter"` | Exporter image repository |
| exporter.image.tag | string | `""` | Exporter image tag (defaults to appVersion) |
| exporter.nodeSelector | object | `{}` | Node selector for Exporter pods |
| exporter.podAnnotations | object | `{}` | Additional pod annotations for Exporter |
| exporter.podLabels | object | `{}` | Additional pod labels for Exporter |
| exporter.podSecurityContext | object | `{"fsGroup":10000}` | Pod security context for Exporter |
| exporter.replicas | int | `1` | Number of Exporter replicas |
| exporter.resources | object | `{}` | Exporter resource requests and limits |
| exporter.secret | object | {} | Sensitive config for Exporter (converted to env vars in Secret) |
| exporter.securityContext | object | `{"runAsGroup":10000,"runAsNonRoot":true,"runAsUser":10000}` | Security context for Exporter container |
| exporter.serviceAccount | object | `{"annotations":{},"create":true,"name":""}` | Service account settings for Exporter |
| exporter.tolerations | list | `[]` | Tolerations for Exporter pods |
| externalRedis.existingSecret | string | `""` | Existing secret containing Redis password Must have key: REDIS_PASSWORD |
| externalRedis.host | string | `""` | External Redis host |
| externalRedis.password | string | `""` | External Redis password |
| externalRedis.port | int | `6379` | External Redis port |
| externalRedis.sentinelMasterSet | string | `""` | Sentinel master set name (for Redis Sentinel) |
| externalURL | string | "" | External URL for Harbor (REQUIRED) This is the URL users will use to access Harbor (e.g., https://harbor.example.com) |
| extraManifests | list | [] | Extra static manifests to deploy These are merged with chart labels and deployed as-is |
| extraTemplateManifests | list | [] | Extra templated manifests to deploy These can use .Values, .Release, and other template functions |
| fullnameOverride | string | `""` | Override the full name |
| gateway | object | `{"enabled":false,"hostnames":[],"parentRefs":[]}` | Gateway API configuration (alternative to ingress) |
| gateway.enabled | bool | `false` | Enable Gateway API HTTPRoute |
| gateway.hostnames | list | `[]` | Hostnames for the HTTPRoute |
| gateway.parentRefs | list | `[]` | Gateway parent references |
| harborAdminPassword | string | "Harbor12345" | Harbor admin password (initial setup) |
| image | object | `{"pullPolicy":"IfNotPresent"}` | Global image settings |
| image.pullPolicy | string | `"IfNotPresent"` | Image pull policy for all Harbor components |
| imagePullSecrets | list | [] | List of image pull secrets |
| ingress | object | `{"annotations":{},"className":"","enabled":true,"hosts":[{"host":"harbor.example.com","paths":[{"path":"/","pathType":"Prefix"}]}],"tls":[]}` | Ingress configuration |
| ingress.annotations | object | `{}` | Ingress annotations |
| ingress.className | string | `""` | Ingress class name |
| ingress.enabled | bool | `true` | Enable ingress |
| ingress.hosts | list | `[{"host":"harbor.example.com","paths":[{"path":"/","pathType":"Prefix"}]}]` | Ingress hosts configuration |
| ingress.tls | list | `[]` | Ingress TLS configuration |
| jobservice.affinity | object | `{}` | Affinity rules for Jobservice pods |
| jobservice.config | object | {} | Jobservice application config (converted to env vars in ConfigMap) |
| jobservice.extraEnv | list | [] | Extra environment variables with valueFrom support |
| jobservice.image | object | `{"repository":"goharbor/harbor-jobservice","tag":""}` | Jobservice image settings |
| jobservice.image.repository | string | `"goharbor/harbor-jobservice"` | Jobservice image repository |
| jobservice.image.tag | string | `""` | Jobservice image tag (defaults to appVersion) |
| jobservice.nodeSelector | object | `{}` | Node selector for Jobservice pods |
| jobservice.podAnnotations | object | `{}` | Additional pod annotations for Jobservice |
| jobservice.podLabels | object | `{}` | Additional pod labels for Jobservice |
| jobservice.podSecurityContext | object | `{"fsGroup":10000}` | Pod security context for Jobservice |
| jobservice.replicas | int | `1` | Number of Jobservice replicas |
| jobservice.resources | object | `{}` | Jobservice resource requests and limits |
| jobservice.secret | object | {} | Sensitive config for Jobservice (converted to env vars in Secret) |
| jobservice.securityContext | object | `{"runAsGroup":10000,"runAsNonRoot":true,"runAsUser":10000}` | Security context for Jobservice container |
| jobservice.serviceAccount | object | `{"annotations":{},"create":true,"name":""}` | Service account settings for Jobservice |
| jobservice.tolerations | list | `[]` | Tolerations for Jobservice pods |
| logLevel | string | `"info"` | Log level for all components (debug, info, warning, error, fatal) |
| metrics.serviceMonitor | object | `{"enabled":false,"honorLabels":true,"interval":"30s","labels":{},"namespace":"","scrapeTimeout":"10s"}` | Enable Prometheus ServiceMonitor |
| metrics.serviceMonitor.enabled | bool | `false` | Create ServiceMonitor resource |
| metrics.serviceMonitor.honorLabels | bool | `true` | Honor labels |
| metrics.serviceMonitor.interval | string | `"30s"` | Scrape interval |
| metrics.serviceMonitor.labels | object | `{}` | Additional labels for ServiceMonitor |
| metrics.serviceMonitor.namespace | string | `""` | ServiceMonitor namespace (defaults to release namespace) |
| metrics.serviceMonitor.scrapeTimeout | string | `"10s"` | Scrape timeout |
| nameOverride | string | `""` | Override the chart name |
| portal.affinity | object | `{}` | Affinity rules for Portal pods |
| portal.config | object | {} | Portal application config (converted to env vars in ConfigMap) |
| portal.extraEnv | list | [] | Extra environment variables with valueFrom support |
| portal.image | object | `{"repository":"goharbor/harbor-portal","tag":""}` | Portal image settings |
| portal.image.repository | string | `"goharbor/harbor-portal"` | Portal image repository |
| portal.image.tag | string | `""` | Portal image tag (defaults to appVersion) |
| portal.nodeSelector | object | `{}` | Node selector for Portal pods |
| portal.podAnnotations | object | `{}` | Additional pod annotations for Portal |
| portal.podLabels | object | `{}` | Additional pod labels for Portal |
| portal.podSecurityContext | object | `{"fsGroup":10000}` | Pod security context for Portal |
| portal.replicas | int | `1` | Number of Portal replicas |
| portal.resources | object | `{}` | Portal resource requests and limits |
| portal.secret | object | {} | Sensitive config for Portal (converted to env vars in Secret) |
| portal.securityContext | object | `{"runAsGroup":10000,"runAsNonRoot":true,"runAsUser":10000}` | Security context for Portal container |
| portal.serviceAccount | object | `{"annotations":{},"create":true,"name":""}` | Service account settings for Portal |
| portal.tolerations | list | `[]` | Tolerations for Portal pods |
| registry.affinity | object | `{}` | Affinity rules for Registry pods |
| registry.config | object | {} | Registry application config (converted to env vars in ConfigMap) |
| registry.controller | object | `{"image":{"repository":"goharbor/harbor-registryctl","tag":""}}` | Registryctl image settings |
| registry.controller.image.repository | string | `"goharbor/harbor-registryctl"` | Registryctl image repository |
| registry.controller.image.tag | string | `""` | Registryctl image tag (defaults to appVersion) |
| registry.extraEnv | list | [] | Extra environment variables with valueFrom support |
| registry.image | object | `{"repository":"goharbor/registry-photon","tag":""}` | Registry image settings |
| registry.image.repository | string | `"goharbor/registry-photon"` | Registry image repository |
| registry.image.tag | string | `""` | Registry image tag (defaults to appVersion) |
| registry.nodeSelector | object | `{}` | Node selector for Registry pods |
| registry.persistence | object | `{"accessModes":["ReadWriteOnce"],"annotations":{},"enabled":true,"existingClaim":"","size":"50Gi","storageClass":""}` | Registry persistence settings |
| registry.persistence.accessModes | list | `["ReadWriteOnce"]` | PVC access modes |
| registry.persistence.annotations | object | `{}` | Annotations for PVC |
| registry.persistence.enabled | bool | `true` | Enable persistence for registry |
| registry.persistence.existingClaim | string | `""` | Existing PVC name (disables dynamic provisioning) |
| registry.persistence.size | string | `"50Gi"` | PVC size |
| registry.persistence.storageClass | string | `""` | Storage class for PVC |
| registry.podAnnotations | object | `{}` | Additional pod annotations for Registry |
| registry.podLabels | object | `{}` | Additional pod labels for Registry |
| registry.podSecurityContext | object | `{"fsGroup":10000}` | Pod security context for Registry |
| registry.replicas | int | `1` | Number of Registry replicas |
| registry.resources | object | `{}` | Registry resource requests and limits |
| registry.secret | object | {} | Sensitive config for Registry (converted to env vars in Secret) |
| registry.securityContext | object | `{"runAsGroup":10000,"runAsNonRoot":true,"runAsUser":10000}` | Security context for Registry container |
| registry.serviceAccount | object | `{"annotations":{},"create":true,"name":""}` | Service account settings for Registry |
| registry.storage | object | `{"azure":{},"filesystem":{"rootdirectory":"/storage"},"gcs":{},"oss":{},"s3":{},"type":"filesystem"}` | Registry storage configuration |
| registry.storage.azure | object | `{}` | Azure Blob storage settings |
| registry.storage.filesystem | object | `{"rootdirectory":"/storage"}` | Filesystem storage settings |
| registry.storage.gcs | object | `{}` | Google Cloud Storage settings |
| registry.storage.oss | object | `{}` | Alibaba Cloud OSS settings |
| registry.storage.s3 | object | `{}` | S3 storage settings |
| registry.storage.type | string | `"filesystem"` | Storage type: filesystem, s3, azure, gcs, oss |
| registry.tolerations | list | `[]` | Tolerations for Registry pods |
| secretKey | string | auto-generated | Secret key for encryption (16 characters) Used for encrypting credentials stored in the database |
| tls.certManager | object | `{"duration":"2160h","enabled":false,"issuerRef":{},"renewBefore":"360h"}` | cert-manager integration |
| tls.certManager.duration | string | `"2160h"` | Certificate duration |
| tls.certManager.enabled | bool | `false` | Enable cert-manager for TLS certificates |
| tls.certManager.issuerRef | object | `{}` | cert-manager issuer reference |
| tls.certManager.renewBefore | string | `"360h"` | Certificate renewal before expiry |
| tls.customSecrets | object | `{"core":"","registry":""}` | Custom TLS secrets (alternative to cert-manager) |
| tls.customSecrets.core | string | `""` | TLS secret for core/portal |
| tls.customSecrets.registry | string | `""` | TLS secret for registry |
| trivy.enabled | bool | `false` | Enable Trivy scanner subchart |
| valkey.architecture | string | `"standalone"` | Valkey architecture: standalone or replication |
| valkey.auth | object | `{"enabled":true,"password":""}` | Valkey authentication settings |
| valkey.enabled | bool | `true` | Enable Valkey subchart |
| valkey.master | object | `{"persistence":{"enabled":false}}` | Valkey master configuration |

## Configuration Examples

### Minimal Production Setup

```yaml
externalURL: https://harbor.example.com

database:
  host: postgres.example.com
  password: your-db-password

ingress:
  enabled: true
  className: nginx
  hosts:
    - host: harbor.example.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: harbor-tls
      hosts:
        - harbor.example.com

valkey:
  enabled: true
```

### S3 Storage Backend

```yaml
registry:
  storage:
    type: s3
    s3:
      region: us-east-1
      bucket: my-harbor-bucket
      accesskey: AKIAIOSFODNN7EXAMPLE
      secretkey: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
  persistence:
    enabled: false
```

### Gateway API

```yaml
ingress:
  enabled: false

gateway:
  enabled: true
  parentRefs:
    - name: my-gateway
      namespace: default
  hostnames:
    - harbor.example.com
```

### External Redis

```yaml
valkey:
  enabled: false

externalRedis:
  host: redis.example.com
  port: 6379
  password: redis-password
```

### CloudNativePG via extraManifests

```yaml
extraManifests:
  - apiVersion: postgresql.cnpg.io/v1
    kind: Cluster
    metadata:
      name: harbor-db
    spec:
      instances: 3
      storage:
        size: 10Gi
```

## Upgrading

### From Legacy Harbor Chart

This chart is not a direct upgrade path from the legacy Harbor chart. You should:

1. Backup your Harbor data
2. Deploy this chart as a new installation
3. Migrate data using Harbor's replication feature

----------------------------------------------
Autogenerated from chart metadata using [helm-docs v1.14.2](https://github.com/norwoodj/helm-docs/releases/v1.14.2)
