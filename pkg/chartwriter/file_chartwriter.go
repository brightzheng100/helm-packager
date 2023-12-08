// Copyright © 2023 Bright Zheng <bright.zheng@outlook.com>
// SPDX-License-Identifier: Apache-2.0

package chartwriter

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	securejoin "github.com/cyphar/filepath-securejoin"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"

	"github.com/brightzheng100/helm-packager/pkg/api"
	"github.com/brightzheng100/helm-packager/pkg/utils"
)

type filechartwriter struct {
	toDir string
}

func NewFileChartWriter(toDir string) *filechartwriter {
	return &filechartwriter{
		toDir: toDir,
	}
}

// Write writes the chart files from fs.FS first
// and then archive it as tarball
func (cw *filechartwriter) Write(ctx context.Context, chart *api.Chart, config api.Config) error {
	// write chart files
	err := cw.writeChartFiles(ctx, chart, config)
	if err != nil {
		return fmt.Errorf("failed to write chart files: %w", err)
	}

	// write to archive the tarball
	err = cw.writeChart(ctx, chart, config)
	if err != nil {
		return fmt.Errorf("failed to write chart tarball: %w", err)
	}

	return nil
}

func (cw *filechartwriter) writeChartFiles(ctx context.Context, chart *api.Chart, config api.Config) error {
	// build tree if enabled
	if config.ChartFilesIncluded {
		utils.AddChartFiles(config.TreeRoot, chart.C.Metadata.Name, chart.C.Raw)
	}

	chartFilesRoot := fmt.Sprintf("%s/%s/%s", cw.toDir, chart.C.Metadata.Name, "chart")

	for _, file := range chart.C.Raw {
		outpath, err := securejoin.SecureJoin(chartFilesRoot, file.Name)
		if err != nil {
			return err
		}

		// Make sure the necessary subdirs get created.
		basedir := filepath.Dir(outpath)
		if err := os.MkdirAll(basedir, 0755); err != nil {
			return err
		}

		if err := os.WriteFile(outpath, file.Data, 0644); err != nil {
			return err
		}
	}

	return nil
}

func (cw *filechartwriter) writeChart(ctx context.Context, chart *api.Chart, config api.Config) error {
	chartFolder := fmt.Sprintf("%s/%s", cw.toDir, chart.C.Metadata.Name)
	chartFilesRoot := fmt.Sprintf("%s/%s/%s", cw.toDir, chart.C.Metadata.Name, "chart")

	ch, err := loader.LoadDir(chartFilesRoot)
	if err != nil {
		return err
	}

	_, err = chartutil.Save(ch, chartFolder)
	if err != nil {
		return fmt.Errorf("failed to write chart %s: %w", chart.C.Metadata.Name, err)
	}

	fileName := fmt.Sprintf("%s-%s.tgz", chart.C.Metadata.Name, chart.C.Metadata.Version)
	utils.AddChart(config.TreeRoot, chart.C.Metadata.Name, fileName)

	return nil
}

func (cw *filechartwriter) Finish(ctx context.Context, config api.Config) error {
	// we need to clean up if chart files are not included as they will be downloaded by default
	// if !config.ChartFilesIncluded {
	// 	for i := 0; i < len(cw.fromCharts); i++ {
	// 		chartVer := strings.Split(cl.fromCharts[i], ":")
	// 		chartName := chartVer[0]

	// 		// remote ${toDIr}/{chartName}/chart/*
	// 		os.RemoveAll(fmt.Sprintf("%s/%s/%s", cl.toDir, chartName, "chart"))
	// 	}
	// }
	return nil
}
