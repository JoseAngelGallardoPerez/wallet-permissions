{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "wallet-permissions.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "wallet-permissions.fullname" -}}
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
{{- define "wallet-permissions.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "wallet-permissions.labels" -}}
helm.sh/chart: {{ include "wallet-permissions.chart" . }}
{{ include "wallet-permissions.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "wallet-permissions.selectorLabels" -}}
app.kubernetes.io/name: {{ include "wallet-permissions.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "wallet-permissions.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "wallet-permissions.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create tag name of the image
*/}}
{{- define "wallet-permissions.imageTag" -}}
{{ .Values.image.tag | default .Chart.AppVersion }}
{{- end }}

{{/*
Create the name of the image repository
*/}}
{{- define "wallet-permissions.imageRepository" -}}
{{ .Values.image.repository | default (printf "velmie/%s" .Chart.Name) }}
{{- end }}

{{/*
Create full image repository name including tag
*/}}
{{- define "wallet-permissions.imageRepositoryWithTag" -}}
{{ include "wallet-permissions.imageRepository" . }}:{{ include "wallet-permissions.imageTag" . }}
{{- end }}

{{/*
Create full database migration image repository name
*/}}
{{- define "wallet-permissions.dbMigrationImageRepositoryWithTag" -}}
{{ include "wallet-permissions.imageRepository" . }}-db-migration:{{ include "wallet-permissions.imageTag" . }}
{{- end }}