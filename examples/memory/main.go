// Copyright © 2023 Bright Zheng <bright.zheng@outlook.com>
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"testing/fstest"

	"github.com/brightzheng100/helm-packager/pkg/chartloader"
	"github.com/brightzheng100/helm-packager/pkg/chartwriter"
	"github.com/brightzheng100/helm-packager/pkg/imageswriter"
	"github.com/brightzheng100/helm-packager/pkg/pipeline"
)

// Chart in memory
const (
	chart = `
apiVersion: v2
name: memory-chart
version: 0.1.0
`
	configmap = `
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ printf "%s-%s" .Chart.Name .Release.Name | trunc 63 | trimSuffix "-" }}
  namespace: {{ .Release.Namespace }}
  labels:
    {{- with .Values.labels -}}
    {{ toYaml . | nindent 4 }}
    {{- end }}
data:
  key: value
`
	deployment = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ printf "%s-%s" .Chart.Name .Release.Name | trunc 63 | trimSuffix "-" }}
  namespace: {{ .Release.Namespace }}
  labels:
    {{- with .Values.labels -}}
    {{ toYaml . | nindent 4 }}
    {{- end }}
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        {{- with .Values.labels -}}
        {{ toYaml . | nindent 8 }}
        {{- end }}
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
`
)

func main() {
	ctx := context.Background()

	// Create chart in memory.
	chartFS := make(fstest.MapFS)
	chartFS["Chart.yaml"] = &fstest.MapFile{Data: []byte(chart)}
	chartFS["templates/configmap.yaml"] = &fstest.MapFile{Data: []byte(configmap)}
	chartFS["templates/deployment.yaml"] = &fstest.MapFile{Data: []byte(deployment)}

	cl := chartloader.NewEmbedChartLoader(chartFS)
	cw := chartwriter.NewStdoutChartWriter()
	iw := imageswriter.NewStdoutImagesWriter()

	cp := pipeline.NewBuilder(ctx).
		WithChartLoader(cl).
		WithChartWriter(cw).
		WithImagesWriter(iw).
		ConfigureChartFilesIncluded(true).
		Complete()

	err := cp.Process()
	if err != nil {
		panic(err)
	}
}
