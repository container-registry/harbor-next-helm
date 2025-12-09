{{/*
Expand the name of the chart.
*/}}
{{- define "harbor.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "harbor.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "harbor.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "harbor.labels" -}}
helm.sh/chart: {{ include "harbor.chart" . }}
{{ include "harbor.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "harbor.selectorLabels" -}}
app.kubernetes.io/name: {{ include "harbor.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Component labels - adds component name to common labels
Usage: {{ include "harbor.componentLabels" (dict "root" . "component" "core") }}
*/}}
{{- define "harbor.componentLabels" -}}
{{ include "harbor.labels" .root }}
app.kubernetes.io/component: {{ .component }}
{{- end }}

{{/*
Component selector labels
Usage: {{ include "harbor.componentSelectorLabels" (dict "root" . "component" "core") }}
*/}}
{{- define "harbor.componentSelectorLabels" -}}
{{ include "harbor.selectorLabels" .root }}
app.kubernetes.io/component: {{ .component }}
{{- end }}

{{/*
=============================================================================
toEnvVars - Convert nested map to flat environment variables
=============================================================================
This is the core helper that enables future-proof configuration.
Any Harbor config option can be set in values.yaml without chart changes.

Usage in ConfigMap:
  {{- include "harbor.toEnvVars" (dict "values" .Values.core.config "prefix" "" "isSecret" false) | nindent 2 }}

Usage in Secret:
  {{- include "harbor.toEnvVars" (dict "values" .Values.core.secret "prefix" "" "isSecret" true) | nindent 2 }}

Example input:
  config:
    storage:
      type: s3
      s3:
        bucket: my-bucket
        region: us-east-1

Example output (ConfigMap):
  STORAGE_TYPE: "s3"
  STORAGE_S3_BUCKET: "my-bucket"
  STORAGE_S3_REGION: "us-east-1"

Example output (Secret):
  STORAGE_TYPE: "czM="
  STORAGE_S3_BUCKET: "bXktYnVja2V0"
  STORAGE_S3_REGION: "dXMtZWFzdC0x"
*/}}
{{- define "harbor.toEnvVars" -}}
{{- $prefix := "" }}
{{- if .prefix }}{{- $prefix = printf "%s_" (.prefix | upper) }}{{- end }}
{{- range $key, $value := .values }}
{{- if kindIs "map" $value }}
{{- /* Recursively process nested maps */ -}}
{{- include "harbor.toEnvVars" (dict "values" $value "prefix" (printf "%s%s" $prefix ($key | upper)) "isSecret" $.isSecret) }}
{{- else if kindIs "slice" $value }}
{{- /* Join arrays with comma */ -}}
{{- if $.isSecret }}
{{ $prefix }}{{ $key | upper }}: {{ $value | join "," | b64enc | quote }}
{{- else }}
{{ $prefix }}{{ $key | upper }}: {{ $value | join "," | quote }}
{{- end }}
{{- else if kindIs "bool" $value }}
{{- /* Handle booleans */ -}}
{{- if $.isSecret }}
{{ $prefix }}{{ $key | upper }}: {{ $value | toString | b64enc | quote }}
{{- else }}
{{ $prefix }}{{ $key | upper }}: {{ $value | toString | quote }}
{{- end }}
{{- else if not (kindIs "invalid" $value) }}
{{- /* Handle strings and numbers, skip nil/empty */ -}}
{{- if $.isSecret }}
{{ $prefix }}{{ $key | upper }}: {{ $value | toString | b64enc | quote }}
{{- else }}
{{ $prefix }}{{ $key | upper }}: {{ $value | toString | quote }}
{{- end }}
{{- end }}
{{- end }}
{{- end }}

{{/*
=============================================================================
Image helpers
=============================================================================
*/}}

{{/*
Return the proper image name
Usage: {{ include "harbor.image" (dict "imageRoot" .Values.core.image "global" .Values.image "chart" .Chart) }}
*/}}
{{- define "harbor.image" -}}
{{- $tag := .imageRoot.tag | default .chart.AppVersion -}}
{{- printf "%s:%s" .imageRoot.repository $tag -}}
{{- end }}

{{/*
Return image pull policy
*/}}
{{- define "harbor.imagePullPolicy" -}}
{{- .Values.image.pullPolicy | default "IfNotPresent" -}}
{{- end }}

{{/*
Return image pull secrets
*/}}
{{- define "harbor.imagePullSecrets" -}}
{{- if .Values.imagePullSecrets }}
imagePullSecrets:
{{- range .Values.imagePullSecrets }}
  - name: {{ . }}
{{- end }}
{{- end }}
{{- end }}

{{/*
=============================================================================
Database helpers
=============================================================================
*/}}

{{/*
Return the database host
*/}}
{{- define "harbor.database.host" -}}
{{- .Values.database.host -}}
{{- end }}

{{/*
Return the database port
*/}}
{{- define "harbor.database.port" -}}
{{- .Values.database.port | default 5432 -}}
{{- end }}

{{/*
Return the database name
*/}}
{{- define "harbor.database.database" -}}
{{- .Values.database.database | default "registry" -}}
{{- end }}

{{/*
Return the database username
*/}}
{{- define "harbor.database.username" -}}
{{- .Values.database.username | default "postgres" -}}
{{- end }}

{{/*
Return the database password secret name
*/}}
{{- define "harbor.database.secretName" -}}
{{- if .Values.database.existingSecret }}
{{- .Values.database.existingSecret }}
{{- else }}
{{- include "harbor.fullname" . }}-database
{{- end }}
{{- end }}

{{/*
Return the database sslmode
*/}}
{{- define "harbor.database.sslmode" -}}
{{- .Values.database.sslmode | default "disable" -}}
{{- end }}

{{/*
Return the full database URL for Core
*/}}
{{- define "harbor.database.coreUrl" -}}
postgresql://{{ include "harbor.database.username" . }}:$(POSTGRESQL_PASSWORD)@{{ include "harbor.database.host" . }}:{{ include "harbor.database.port" . }}/{{ include "harbor.database.database" . }}?sslmode={{ include "harbor.database.sslmode" . }}
{{- end }}

{{/*
=============================================================================
Redis helpers
=============================================================================
*/}}

{{/*
Return the Redis host
*/}}
{{- define "harbor.redis.host" -}}
{{- if .Values.valkey.enabled }}
{{- include "harbor.fullname" . }}-valkey-master
{{- else }}
{{- .Values.externalRedis.host }}
{{- end }}
{{- end }}

{{/*
Return the Redis port
*/}}
{{- define "harbor.redis.port" -}}
{{- if .Values.valkey.enabled }}
{{- 6379 }}
{{- else }}
{{- .Values.externalRedis.port | default 6379 }}
{{- end }}
{{- end }}

{{/*
Return the Redis password secret name
*/}}
{{- define "harbor.redis.secretName" -}}
{{- if .Values.valkey.enabled }}
{{- include "harbor.fullname" . }}-valkey
{{- else if .Values.externalRedis.existingSecret }}
{{- .Values.externalRedis.existingSecret }}
{{- else }}
{{- include "harbor.fullname" . }}-redis
{{- end }}
{{- end }}

{{/*
Return the Redis password key in the secret
*/}}
{{- define "harbor.redis.secretKey" -}}
{{- if .Values.valkey.enabled -}}
valkey-password
{{- else -}}
REDIS_PASSWORD
{{- end -}}
{{- end }}

{{/*
Return the Redis URL for Harbor components
*/}}
{{- define "harbor.redis.url" -}}
redis://:$(REDIS_PASSWORD)@{{ include "harbor.redis.host" . }}:{{ include "harbor.redis.port" . }}/0
{{- end }}

{{/*
=============================================================================
Internal URL helpers
=============================================================================
*/}}

{{/*
Return the Core internal URL
*/}}
{{- define "harbor.core.url" -}}
http://{{ include "harbor.fullname" . }}-core:80
{{- end }}

{{/*
Return the Portal internal URL
*/}}
{{- define "harbor.portal.url" -}}
http://{{ include "harbor.fullname" . }}-portal:80
{{- end }}

{{/*
Return the Registry internal URL
*/}}
{{- define "harbor.registry.url" -}}
http://{{ include "harbor.fullname" . }}-registry:5000
{{- end }}

{{/*
Return the Registry controller internal URL
*/}}
{{- define "harbor.registryctl.url" -}}
http://{{ include "harbor.fullname" . }}-registry:8080
{{- end }}

{{/*
Return the Jobservice internal URL
*/}}
{{- define "harbor.jobservice.url" -}}
http://{{ include "harbor.fullname" . }}-jobservice:80
{{- end }}

{{/*
Return the Trivy adapter URL (if enabled)
*/}}
{{- define "harbor.trivy.url" -}}
{{- if .Values.trivy.enabled }}
http://{{ include "harbor.fullname" . }}-harbor-scanner-trivy:8080
{{- end }}
{{- end }}

{{/*
=============================================================================
Service Account helpers
=============================================================================
*/}}

{{/*
Create the name of the service account to use for a component
Usage: {{ include "harbor.serviceAccountName" (dict "root" . "component" "core" "serviceAccount" .Values.core.serviceAccount) }}
*/}}
{{- define "harbor.serviceAccountName" -}}
{{- if .serviceAccount.create }}
{{- default (printf "%s-%s" (include "harbor.fullname" .root) .component) .serviceAccount.name }}
{{- else }}
{{- default "default" .serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
=============================================================================
Secret Key helpers
=============================================================================
*/}}

{{/*
Return the secret key for encryption
*/}}
{{- define "harbor.secretKey" -}}
{{- if .Values.secretKey }}
{{- .Values.secretKey }}
{{- else }}
{{- /* Generate a deterministic key based on release name */ -}}
{{- $key := printf "%s-harbor-secret-key" .Release.Name | sha256sum | trunc 16 }}
{{- $key }}
{{- end }}
{{- end }}

{{/*
=============================================================================
External URL helpers
=============================================================================
*/}}

{{/*
Return the external URL
*/}}
{{- define "harbor.externalURL" -}}
{{- .Values.externalURL }}
{{- end }}

{{/*
Return the core external URL (same as externalURL for now)
*/}}
{{- define "harbor.coreURL" -}}
{{- include "harbor.externalURL" . }}
{{- end }}

{{/*
=============================================================================
TLS helpers
=============================================================================
*/}}

{{/*
Check if internal TLS is enabled
*/}}
{{- define "harbor.internalTLS.enabled" -}}
{{- false }}
{{- end }}

{{/*
Return the TLS secret name for a component
Usage: {{ include "harbor.tlsSecretName" (dict "root" . "component" "core") }}
*/}}
{{- define "harbor.tlsSecretName" -}}
{{- if eq .component "core" }}
{{- if .root.Values.tls.customSecrets.core }}
{{- .root.Values.tls.customSecrets.core }}
{{- else }}
{{- include "harbor.fullname" .root }}-core-tls
{{- end }}
{{- else if eq .component "registry" }}
{{- if .root.Values.tls.customSecrets.registry }}
{{- .root.Values.tls.customSecrets.registry }}
{{- else }}
{{- include "harbor.fullname" .root }}-registry-tls
{{- end }}
{{- end }}
{{- end }}

{{/*
=============================================================================
Validation helpers
=============================================================================
*/}}

{{/*
Validate required values
*/}}
{{- define "harbor.validateValues" -}}
{{- if not .Values.externalURL }}
{{- fail "externalURL is required. Please set externalURL in your values." }}
{{- end }}
{{- if not .Values.database.host }}
{{- fail "database.host is required. Please set database.host in your values." }}
{{- end }}
{{- end }}
