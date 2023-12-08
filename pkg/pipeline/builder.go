// Copyright © 2023 Bright Zheng <bright.zheng@outlook.com>
// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"context"

	"github.com/brightzheng100/helm-packager/pkg/api"
	"github.com/brightzheng100/helm-packager/pkg/utils"
)

type Builder struct {
	cp *packager
}

func NewBuilder(ctx context.Context) *Builder {
	return &Builder{
		cp: &packager{
			ctx: ctx,
			Config: api.Config{
				TreeRoot: utils.NewRootTree(),
			},
		},
	}
}

func (pb *Builder) ConfigureChartFilesIncluded(included bool) *Builder {
	pb.cp.ChartFilesIncluded = included
	return pb
}

func (pb *Builder) ConfigureDryrun(dryrun bool) *Builder {
	pb.cp.Dryrun = dryrun
	return pb
}

func (pb *Builder) WithChartLoader(cl api.ChartLoader) *Builder {
	pb.cp.cl = cl
	return pb
}

func (pb *Builder) WithChartWriter(cw api.ChartWriter) *Builder {
	pb.cp.cw = cw
	return pb
}

func (pb *Builder) WithImagesWriter(iw api.ImagesWriter) *Builder {
	pb.cp.iw = iw
	return pb
}

func (pb *Builder) Complete() *packager {
	return pb.cp
}
