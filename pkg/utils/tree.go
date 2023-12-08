// Copyright © 2023 Bright Zheng <bright.zheng@outlook.com>
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/brightzheng100/helm-packager/pkg/api"
	"github.com/xlab/treeprint"
	"helm.sh/helm/v3/pkg/chart"
)

type tree struct {
	t api.Tree
}

func NewRootTree() *api.Tree {
	r := treeprint.New()
	return &api.Tree{
		T: r,
	}
}

func NewTree(t *api.Tree) *api.Tree {
	tree := treeprint.NewWithRoot(t.T)
	return &api.Tree{T: tree}
}

func AddChart(t *api.Tree, chartName, fileName string) {
	t.T.AddBranch(chartName).AddNode(fileName)
}

func AddChartFiles(t *api.Tree, chartName string, files []*chart.File) {
	chartBranch := t.T.AddBranch("chart")
	for _, f := range files {
		chartBranch.AddNode(f.Name)
	}
}

func AddChartImages(t *api.Tree, chartName string, images []string) {
	chartBranch := t.T.FindByValue(chartName)
	imageBranch := chartBranch.AddBranch("images")

	for _, imgref := range images {
		nametag := strings.Split(filepath.Base(imgref), ":")
		imageBranch.AddNode(fmt.Sprintf("%s-%s.tar (%s)", nametag[0], nametag[1], imgref))
	}
}

func Print(t *api.Tree) {
	fmt.Println(t.T.String())
}
