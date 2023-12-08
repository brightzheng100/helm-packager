// Copyright © 2023 Bright Zheng <bright.zheng@outlook.com>
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"os"

	"github.com/brightzheng100/helm-packager/pkg/api"
	"github.com/brightzheng100/helm-packager/pkg/chartloader"
	"github.com/brightzheng100/helm-packager/pkg/chartwriter"
	"github.com/brightzheng100/helm-packager/pkg/imageswriter"
	"github.com/brightzheng100/helm-packager/pkg/pipeline"
	"github.com/spf13/cobra"
)

var pullCmdLongDesc = `  Pull command pulls the remote Helm charts and their images to local.

  Usage:

  helm-packager pull
    --from-chart-repo <REMOTE_REPOSITORY_URL>
    --from-charts <CHART_NAME>[:<CHART_VERSION>][,<CHART_NAME>[:<CHART_VERSION>]]
   [--to-dir <CHARTS_DIR>]
   [--char-files-included true/false]

  Examples:

  # Pull Helm chart "nginx" from Bitnami repository and its images to specific local folder

  helm-packager pull \
    --from-chart-repo oci://registry-1.docker.io/bitnamicharts \
    --from-charts nginx \
    --to-dir ./charts

  # Pull to print info about Helm charts "apache" (with specified version) and "nginx" and their images
  
  helm-packager pull \
    --from-chart-repo oci://registry-1.docker.io/bitnamicharts \
    --from-charts apache:10.2.3,nginx
`

var p = &pull{}

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull command pulls the remote Helm charts and their images to local",
	Long:  pullCmdLongDesc,
	Run: func(cmd *cobra.Command, args []string) {
		runPull(p, args)
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)

	pullCmd.Flags().StringVar(&p.fromChartRepo, "from-chart-repo", "", "Helm repository URL, e.g. https://charts.bitnami.com/bitnami")
	pullCmd.Flags().StringSliceVar(&p.fromCharts, "from-charts", []string{}, "Helm chart(s) with optional version tag, separated by commar, e.g. apache:10.2.3,nginx")
	pullCmd.Flags().StringVar(&p.toDir, "to-dir", "", "Optional, the directory for pulled Helm charts and their images' tarball files. When not specified, the command will only print out the structure")
	pullCmd.Flags().BoolVar(&p.chartFilesIncluded, "char-files-included", false, "Optional, the flag to indicate whether the chart files should be included and pulled")

	pullCmd.MarkFlagRequired("from-chart-repo")
	pullCmd.MarkFlagRequired("from-charts")
}

type pull struct {
	fromChartRepo      string
	fromCharts         []string
	toDir              string
	chartFilesIncluded bool
}

func runPull(pull *pull, args []string) {
	ctx := context.Background()

	var cl api.ChartLoader
	var cw api.ChartWriter
	var iw api.ImagesWriter

	if pull.toDir == "" {
		cl = chartloader.NewRemoteChartLoader(pull.fromChartRepo, pull.fromCharts, ".")
		cw = chartwriter.NewStdoutChartWriter()
		iw = imageswriter.NewStdoutImagesWriter()
	} else {
		// create folder if needed
		if _, err := os.Stat(pull.toDir); os.IsNotExist(err) {
			err := os.Mkdir(pull.toDir, 0766)
			if err != nil {
				panic(err)
			}
		}

		cl = chartloader.NewRemoteChartLoader(pull.fromChartRepo, pull.fromCharts, pull.toDir)
		cw = chartwriter.NewFileChartWriter(pull.toDir)
		iw = imageswriter.NewFileImagesWriter(pull.toDir)
	}

	cp := pipeline.NewBuilder(ctx).
		WithChartLoader(cl).
		WithChartWriter(cw).
		WithImagesWriter(iw).
		ConfigureChartFilesIncluded(pull.chartFilesIncluded).
		Complete()

	err := cp.Process()
	if err != nil {
		panic(err)
	}
}
