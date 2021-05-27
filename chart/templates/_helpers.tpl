{{- define "labels.standard" -}}
app: {{ .Values.name }}
chart: "{{ .Chart.Name }}-{{ .Chart.Version }}"
release: "{{ .Release.Name }}"
heritage: "{{ .Release.Service }}"
{{- end -}}

{{- define "final.image.version" -}}
{{ if .Values.image.version }}
{{ .Values.image.version }}
{{ else }}
{{ .Chart.Version }}
{{ end }}
{{- end -}}