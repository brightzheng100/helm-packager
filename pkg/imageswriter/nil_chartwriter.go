// Copyright © 2023 Bright Zheng <bright.zheng@outlook.com>
// SPDX-License-Identifier: Apache-2.0

package imageswriter

import (
	"context"

	"github.com/brightzheng100/helm-packager/pkg/api"
)

type nilimageswriter struct {
}

func NewNilImagesWriter() *nilimageswriter {
	return &nilimageswriter{}
}

func (iw *nilimageswriter) Write(ctx context.Context, chart *api.Chart, config api.Config) error {
	// do nothing
	return nil
}

func (iw *nilimageswriter) Finish(ctx context.Context, config api.Config) error {
	return nil
}
