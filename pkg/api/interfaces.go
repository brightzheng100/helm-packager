// Copyright © 2023 Bright Zheng <bright.zheng@outlook.com>
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
)

// ChartPackager defines the interfaces that handle a ReadOnly Helm chart.
type ChartPackager interface {
	GetPackager() *ChartPackager

	// Process is to process based on the built pipeline
	Process() error
}

// ChartLoader defines the interfaces for how to load a Helm chart
type ChartLoader interface {
	Load(ctx context.Context, config Config) ([]*Chart, error)
	Finish(ctx context.Context, config Config) error
}

// ChartWriter defines the interfaces for how to write a Helm chart
type ChartWriter interface {
	Write(ctx context.Context, chart *Chart, config Config) error
	Finish(ctx context.Context, config Config) error
}

// ImagesWriter defines the interfaces for how to write the images of a given Helm chart
type ImagesWriter interface {
	Write(ctx context.Context, chart *Chart, config Config) error
	Finish(ctx context.Context, config Config) error
}
