// Copyright © 2023 Bright Zheng <bright.zheng@outlook.com>
// SPDX-License-Identifier: Apache-2.0

package chartloader

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/brightzheng100/helm-packager/pkg/api"
	"helm.sh/helm/v3/pkg/chart/loader"
)

type embedchartloader struct {
	fs fs.FS
}

func NewEmbedChartLoader(fs fs.FS) *embedchartloader {
	return &embedchartloader{
		fs: fs,
	}
}

func (cl *embedchartloader) Load(ctx context.Context, config api.Config) ([]*api.Chart, error) {
	files := []*loader.BufferedFile{}

	err := fs.WalkDir(cl.fs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		data, err := fs.ReadFile(cl.fs, path)
		if err != nil {
			return fmt.Errorf("could not read manifest %s: %w", path, err)
		}

		files = append(files, &loader.BufferedFile{
			Name: path,
			Data: data,
		})

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not walk chart directory: %w", err)
	}

	chart, err := loader.LoadFiles(files)
	if err != nil {
		return nil, fmt.Errorf("could not load chart from files: %w", err)
	}

	return []*api.Chart{{C: chart}}, nil
}

func (cl *embedchartloader) Finish(ctx context.Context, config api.Config) error {
	return nil
}
