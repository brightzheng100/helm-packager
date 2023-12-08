// Copyright © 2023 Bright Zheng <bright.zheng@outlook.com>
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"os"

	"github.com/brightzheng100/helm-packager/pkg/chartloader"
	"github.com/brightzheng100/helm-packager/pkg/chartwriter"
	"github.com/brightzheng100/helm-packager/pkg/imageswriter"
	"github.com/brightzheng100/helm-packager/pkg/pipeline"
)

var (
	fromRepo   = "oci://registry-1.docker.io/bitnamicharts"
	fromCharts = []string{"apache:10.2.3", "nginx"}
	toDir      = "./_charts"
)

func main() {
	ctx := context.Background()

	// create folder if needed
	if _, err := os.Stat(toDir); os.IsNotExist(err) {
		err := os.Mkdir(toDir, 0766)
		if err != nil {
			panic(err)
		}
	}

	cl := chartloader.NewRemoteChartLoader(fromRepo, fromCharts, toDir)
	cw := chartwriter.NewStdoutChartWriter()
	iw := imageswriter.NewFileImagesWriter(toDir)

	cp := pipeline.NewBuilder(ctx).
		WithChartLoader(cl).
		WithChartWriter(cw).
		WithImagesWriter(iw).
		ConfigureChartFilesIncluded(false).
		Complete()

	err := cp.Process()
	if err != nil {
		panic(err)
	}
}
