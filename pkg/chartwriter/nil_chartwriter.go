// Copyright © 2023 Bright Zheng <bright.zheng@outlook.com>
// SPDX-License-Identifier: Apache-2.0

package chartwriter

import (
	"context"

	"github.com/brightzheng100/helm-packager/pkg/api"
)

type nilchartwriter struct {
}

func NewNilChartWriter() *nilchartwriter {
	return &nilchartwriter{}
}

func (cw *nilchartwriter) Write(ctx context.Context, chart *api.Chart, config api.Config) error {
	// do nothing
	return nil
}

func (cw *nilchartwriter) Finish(ctx context.Context, config api.Config) error {
	return nil
}
