// Copyright © 2023 Bright Zheng <bright.zheng@outlook.com>
// SPDX-License-Identifier: Apache-2.0

package chartwriter

import (
	"context"
	"fmt"

	"github.com/brightzheng100/helm-packager/pkg/api"
	"github.com/brightzheng100/helm-packager/pkg/utils"
)

type stdoutchartwriter struct {
}

func NewStdoutChartWriter() *stdoutchartwriter {
	return &stdoutchartwriter{}
}

func (cw *stdoutchartwriter) Write(ctx context.Context, chart *api.Chart, config api.Config) error {
	//cw.writeMetadata(ctx, chart)

	// write chart
	cw.writeChart(ctx, chart, config)

	// write chart files if enabled
	if config.ChartFilesIncluded {
		cw.writeChartFiles(ctx, chart, config)
	}

	return nil
}

// func (cw *stdoutchartwriter) writeMetadata(ctx context.Context, chart *api.Chart) {
// 	fmt.Println("----------------------------")
// 	fmt.Printf("---- Name: %s \n", chart.C.Metadata.Name)
// 	fmt.Printf("---- Version: %s \n", chart.C.Metadata.Version)
// 	fmt.Println("----------------------------")
// }

func (cw *stdoutchartwriter) writeChart(ctx context.Context, chart *api.Chart, config api.Config) {
	fileName := fmt.Sprintf("%s-%s.tgz", chart.C.Metadata.Name, chart.C.Metadata.Version)

	utils.AddChart(config.TreeRoot, chart.C.Metadata.Name, fileName)
}

func (cw *stdoutchartwriter) writeChartFiles(ctx context.Context, chart *api.Chart, config api.Config) {
	utils.AddChartFiles(config.TreeRoot, chart.C.Metadata.Name, chart.C.Raw)
}

func (cw *stdoutchartwriter) Finish(ctx context.Context, config api.Config) error {
	return nil
}
