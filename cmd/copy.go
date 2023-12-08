// Copyright © 2023 Bright Zheng <bright.zheng@outlook.com>
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var copyCmdLongDesc = `  Copy command copies Helm charts and their images from source to target Helm chart repository / image registry.

  Usage:

  helm-packager copy \
    --from-chart-repo <REMOTE_REPOSITORY_URL> \
    --from-charts <CHART_NAME>[:<CHART_VERSION>][,<CHART_NAME>[:<CHART_VERSION>]] \
    --to-chart-repo <TARGETED HELM REPOSITORY TO PUSH CHARTS TO> \
    --to-image-registry <TARGETED IMAGE REGISTRY TO PUSH IMAGES TO>

  Examples:

  # Copy Helm charts "apache" with specific version "10.2.3" and another Helm chart "nginx" from Bitnami repository to the private Helm chart repository / image registry.

  helm-packager copy \
    --from-chart-repo oci://registry-1.docker.io/bitnamicharts \
    --from-charts apache:10.2.3,nginx \
    --to-chart-repo https://my.chart.repo \
    --to-image-resitry https://my.docker.registry
`

var c = &copy{}

// copyCmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy command copies Helm charts and their images from source to target Helm chart repository / image registry.",
	Long:  copyCmdLongDesc,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("copy called")
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)

	copyCmd.Flags().StringVar(&c.fromChartRepo, "from-dir", "", "Local directory that has exported Helm charts and images, e.g. ./charts")
	copyCmd.Flags().StringSliceVar(&c.fromCharts, "from-charts", []string{}, "Helm chart(s) with optional version tag, separated by commar, e.g. apache:10.2.3,nginx")
	copyCmd.Flags().StringVar(&c.ToChartRepo, "to-chart-repo", "", "The target Helm chart repository URL")
	copyCmd.Flags().StringVar(&c.ToImageRegistry, "to-image-registry", "", "The target image registry URL")

	copyCmd.MarkFlagRequired("from-dir")
	copyCmd.MarkFlagRequired("to-chart-repo")
	copyCmd.MarkFlagRequired("to-image-registry")
}

type copy struct {
	fromChartRepo   string
	fromCharts      []string
	ToChartRepo     string
	ToImageRegistry string
}

func runCopy(c *copy, args []string) {
	fmt.Println("...to be implemented...")
}
