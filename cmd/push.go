// Copyright © 2023 Bright Zheng <bright.zheng@outlook.com>
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var pushCmdLongDesc = `  Push command pushs the exported local Helm charts and their images to remote Helm repository (e.g. ChartMuseum) and image registry.

  Usage:

  helm-packager push \
    --from-dir <EXPORTED DIR WITH CHARTS AND IMAGES> \
   [--from-charts <CHART_NAME>[,<CHART_NAME>] \
    --to-chart-repo <TARGETED HELM REPOSITORY TO PUSH CHARTS TO> \
    --to-image-registry <TARGETED IMAGE REGISTRY TO PUSH IMAGES TO>

  Examples:

  # Push all charts and their images witin a specified "./_charts" folder:

  helm-packager push \
    --from-dir ./_charts \
    --to-chart-repo https://my.chart.repo \
    --to-image-resitry https://my.docker.registry

  # Or to selectively push one or some charts and their images witin a specified "./_charts" folder:
  
  helm-packager push \
    --from-dir ./_charts \
    --from-charts apache,nginx \
    --to-chart-repo https://my.chart.repo \
    --to-image-resitry https://my.docker.registry
`

var s = &push{}

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push command pushs the exported local Helm charts and their images to remote Helm repository (e.g. ChartMuseum) and image registry.",
	Long:  pushCmdLongDesc,
	Run: func(cmd *cobra.Command, args []string) {
		runPush(s, args)
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)

	pushCmd.Flags().StringVar(&s.fromDir, "from-dir", "", "Local directory that has exported Helm charts and images, e.g. ./charts")
	pushCmd.Flags().StringSliceVar(&s.fromCharts, "from-charts", []string{}, "Helm chart(s) with optional version tag, separated by commar, e.g. apache:10.2.3,nginx")
	pushCmd.Flags().StringVar(&s.toChartRepo, "to-chart-repo", "", "The target Helm chart repository URL")
	pushCmd.Flags().StringVar(&s.toImageRegistry, "to-image-registry", "", "The target image registry URL")

	pushCmd.MarkFlagRequired("from-dir")
	pushCmd.MarkFlagRequired("to-chart-repo")
	pushCmd.MarkFlagRequired("to-image-registry")
}

type push struct {
	fromDir         string
	fromCharts      []string
	toChartRepo     string
	toImageRegistry string
}

func runPush(push *push, args []string) {
	fmt.Println("...to be implemented...")
}
