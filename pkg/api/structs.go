// Copyright © 2023 Bright Zheng <bright.zheng@outlook.com>
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"github.com/xlab/treeprint"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/registry"
)

// Chart represents a loaded Helm chart.
type Chart struct {
	C *chart.Chart
}

// Config represents the configuration in the pipeline
type Config struct {
	ChartFilesIncluded bool
	Dryrun             bool

	TreeRoot *Tree
}

// Tree is a wrapper of treeprint.Tree for tree view display
type Tree struct {
	T treeprint.Tree
}

// RemoteChart represents a Helm chart from repot repository
type RemoteChart struct {
	Chart

	CaFile                string // --ca-file
	CertFile              string // --cert-file
	KeyFile               string // --key-file
	InsecureSkipTLSverify bool   // --insecure-skip-verify
	PlainHTTP             bool   // --plain-http
	Keyring               string // --keyring
	Password              string // --password
	PassCredentialsAll    bool   // --pass-credentials
	RepoURL               string // --repo
	Username              string // --username
	Verify                bool   // --verify
	Version               string // --version

	// registryClient provides a registry client but is not added with
	// options from a flag
	registryClient *registry.Client
}
