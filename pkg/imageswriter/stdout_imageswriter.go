// Copyright © 2023 Bright Zheng <bright.zheng@outlook.com>
// SPDX-License-Identifier: Apache-2.0

package imageswriter

import (
	"context"
	"fmt"
	"regexp"

	"github.com/brightzheng100/helm-packager/pkg/api"
	"github.com/brightzheng100/helm-packager/pkg/utils"
)

var expression = "..|.image? | select(.)"
var r, _ = regexp.Compile("^[A-Za-z]")

type stdoutimageswriter struct {
}

func NewStdoutImagesWriter() *stdoutimageswriter {
	return &stdoutimageswriter{}
}

func (iw *stdoutimageswriter) Write(ctx context.Context, chart *api.Chart, config api.Config) error {
	// Templatize the chart
	manifest, err := Templatize(ctx, chart)
	if err != nil {
		return fmt.Errorf("could not templatize Helm chart: %w", err)
	}

	images, err := ExtractImages(ctx, manifest)
	if err != nil {
		return fmt.Errorf("could not extract images from Helm chart: %w", err)
	}

	if err = iw.writeImages(ctx, chart.C.Metadata.Name, images, config); err != nil {
		return fmt.Errorf("could not write images from Helm chart: %w", err)
	}

	return nil
}

func (iw *stdoutimageswriter) writeImages(ctx context.Context, chartName string, images []string, config api.Config) error {
	utils.AddChartImages(config.TreeRoot, chartName, images)
	return nil
}

func (iw *stdoutimageswriter) Finish(ctx context.Context, config api.Config) error {
	return nil
}
