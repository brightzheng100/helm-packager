// Copyright © 2023 Bright Zheng <bright.zheng@outlook.com>
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"os"

	"github.com/brightzheng100/helm-packager/pkg/chartloader"
	"github.com/brightzheng100/helm-packager/pkg/chartwriter"
	"github.com/brightzheng100/helm-packager/pkg/imageswriter"
	"github.com/brightzheng100/helm-packager/pkg/pipeline"
)

//go:embed charts
var embeddedCharts embed.FS

var chartName = flag.String("chart", "robotshop", "Chart name.")

func main() {
	flag.Parse()

	chartFS, err := fs.Sub(embeddedCharts, fmt.Sprintf("charts/%s", *chartName))
	if err != nil {
		panic(err)
	}

	// stdout
	stdout(chartFS)

	// file
	file(chartFS)
}

func stdout(chartFS fs.FS) {
	cl := chartloader.NewEmbedChartLoader(chartFS)
	cw := chartwriter.NewStdoutChartWriter()
	iw := imageswriter.NewStdoutImagesWriter()

	ctx := context.Background()

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

func file(chartFS fs.FS) {

	toDir := "./_charts"

	// create folder if needed
	if _, err := os.Stat(toDir); os.IsNotExist(err) {
		err := os.Mkdir(toDir, 0766)
		if err != nil {
			panic(err)
		}
	}

	cl := chartloader.NewEmbedChartLoader(chartFS)
	cw := chartwriter.NewFileChartWriter(toDir)
	iw := imageswriter.NewFileImagesWriter(toDir)

	ctx := context.Background()

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
