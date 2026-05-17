{{/*
Expand the name of the chart.
*/}}
{{- define "moneyrabbit.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "moneyrabbit.fullname" -}}
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

{{- define "moneyrabbit.frontend.fullname" -}}
{{- printf "%s-frontend" (include "moneyrabbit.fullname" .) | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "moneyrabbit.backend.fullname" -}}
{{- printf "%s-backend" (include "moneyrabbit.fullname" .) | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "moneyrabbit.mariadb.fullname" -}}
{{- printf "%s-mariadb" (include "moneyrabbit.fullname" .) | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "moneyrabbit.labels" -}}
helm.sh/chart: {{ printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
app.kubernetes.io/name: {{ include "moneyrabbit.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
{{- end }}

{{- define "moneyrabbit.frontend.selectorLabels" -}}
app.kubernetes.io/name: {{ include "moneyrabbit.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/component: frontend
{{- end }}

{{- define "moneyrabbit.backend.selectorLabels" -}}
app.kubernetes.io/name: {{ include "moneyrabbit.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/component: backend
{{- end }}
