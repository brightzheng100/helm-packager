// Copyright © 2023 Bright Zheng <bright.zheng@outlook.com>
// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"context"
	"fmt"
	"sort"

	"github.com/brightzheng100/helm-packager/pkg/api"
)

type packager struct {
	ctx context.Context

	api.Config

	chart *api.Chart
	tree  *api.Tree

	cl api.ChartLoader
	cw api.ChartWriter
	iw api.ImagesWriter
}

func (cp *packager) GetPackager() *packager {
	return cp
}

func (cp *packager) Process() error {

	// load charts
	charts, err := cp.cl.Load(cp.ctx, cp.Config)
	if err != nil {
		return fmt.Errorf("could not load Helm charts: %w", err)
	}

	// sort charts by name
	sort.SliceStable(charts, func(i, j int) bool {
		return charts[i].C.Metadata.Name < charts[j].C.Metadata.Name
	})

	for _, chart := range charts {
		// write chart
		err = cp.cw.Write(cp.ctx, chart, cp.Config)
		if err != nil {
			return fmt.Errorf("could not write Helm chart: %w", err)
		}

		// write images
		err = cp.iw.Write(cp.ctx, chart, cp.Config)
		if err != nil {
			return fmt.Errorf("could not write images: %w", err)
		}
	}

	// clean up
	cp.cl.Finish(cp.ctx, cp.Config)
	cp.cw.Finish(cp.ctx, cp.Config)
	cp.iw.Finish(cp.ctx, cp.Config)

	// output
	fmt.Println(cp.Config.TreeRoot.T.String())

	return nil
}
