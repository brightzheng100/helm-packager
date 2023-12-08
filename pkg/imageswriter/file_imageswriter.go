// Copyright © 2023 Bright Zheng <bright.zheng@outlook.com>
// SPDX-License-Identifier: Apache-2.0

package imageswriter

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/brightzheng100/helm-packager/pkg/api"
	"github.com/brightzheng100/helm-packager/pkg/utils"
	"github.com/google/go-containerregistry/pkg/crane"
)

type fileimageswriter struct {
	toDir string
}

func NewFileImagesWriter(toDir string) *fileimageswriter {
	return &fileimageswriter{
		toDir: toDir,
	}
}

func (iw *fileimageswriter) Write(ctx context.Context, chart *api.Chart, config api.Config) error {
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

func (iw *fileimageswriter) writeImages(ctx context.Context, chartName string, images []string, config api.Config) error {
	utils.AddChartImages(config.TreeRoot, chartName, images)

	imgDir := fmt.Sprintf("%s/%s/%s/", iw.toDir, chartName, "images")
	if err := os.MkdirAll(imgDir, 0755); err != nil {
		return fmt.Errorf("failed to mkdir %s: ", imgDir)
	}

	// docker.io/bitnami/apache:2.4.58-debian-11-r1
	for _, imgref := range images {
		image, err := crane.Pull(imgref)
		if err != nil {
			return err
		}

		nametag := strings.Split(filepath.Base(imgref), ":")

		err = crane.Save(image, imgref, fmt.Sprintf("%s/%s-%s.tar", imgDir, nametag[0], nametag[1]))
		if err != nil {
			return err
		}
	}

	return nil
}

func (iw *fileimageswriter) Finish(ctx context.Context, config api.Config) error {
	return nil
}
