// Copyright © 2023 Bright Zheng <bright.zheng@outlook.com>
// SPDX-License-Identifier: Apache-2.0

package chartloader

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/brightzheng100/helm-packager/pkg/api"
	securejoin "github.com/cyphar/filepath-securejoin"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/registry"
	"helm.sh/helm/v3/pkg/repo"
)

var settings = cli.New()

type remotechartloader struct {
	fromChartRepo string
	fromCharts    []string
	toDir         string
}

func NewRemoteChartLoader(fromChartRepo string, fromCharts []string, toDir string) *remotechartloader {
	return &remotechartloader{
		fromChartRepo: fromChartRepo,
		fromCharts:    fromCharts,
		toDir:         toDir,
	}
}

func (cl *remotechartloader) Load(ctx context.Context, config api.Config) ([]*api.Chart, error) {
	var charts []*api.Chart

	actionConfig := new(action.Configuration)
	client := action.NewPullWithOpts(action.WithConfig(actionConfig))

	client.Settings = settings

	registryClient, err := newRegistryClient(client.CertFile, client.KeyFile, client.CaFile,
		client.InsecureSkipTLSverify, client.PlainHTTP)
	if err != nil {
		return nil, fmt.Errorf("missing registry client: %w", err)
	}
	client.SetRegistryClient(registryClient)

	for i := 0; i < len(cl.fromCharts); i++ {
		client.Version = ">0.0.0-0"
		chartVer := strings.Split(cl.fromCharts[i], ":")
		chartName := chartVer[0]
		if len(chartVer) == 2 {
			chartVersion := chartVer[1]
			client.Version = chartVersion
		}

		url := fmt.Sprintf("%s/%s", cl.fromChartRepo, chartName)

		client.DestDir = fmt.Sprintf("%s/%s", cl.toDir, chartName)
		client.Untar = true       // always untar for image processing
		client.UntarDir = "chart" // fmt.Sprintf("%s/chart", chartName)

		if err := os.MkdirAll(fmt.Sprintf("%s/%s/%s", cl.toDir, chartName, "chart"), 0755); err != nil {
			return nil, errors.Wrap(err, "failed to untar (mkdir)")
		}

		//output, err := client.Run(url)
		output, chartpath, err := cl.run(client, registryClient, url, config)
		if err != nil {
			return nil, err
		}
		fmt.Fprint(os.Stdout, output)

		chart, err := loader.Load(chartpath)
		if err != nil {
			fmt.Println(fmt.Sprintf("could not load chart '%s' from %s: %s", chartName, cl.fromChartRepo, err))
		}
		charts = append(charts, &api.Chart{C: chart})
	}

	return charts, nil
}

// copied from https://github.com/helm/helm/blob/main/cmd/helm/root.go
func newRegistryClient(certFile, keyFile, caFile string, insecureSkipTLSverify, plainHTTP bool) (*registry.Client, error) {
	if certFile != "" && keyFile != "" || caFile != "" || insecureSkipTLSverify {
		registryClient, err := newRegistryClientWithTLS(certFile, keyFile, caFile, insecureSkipTLSverify)
		if err != nil {
			return nil, err
		}
		return registryClient, nil
	}
	registryClient, err := newDefaultRegistryClient(plainHTTP)
	if err != nil {
		return nil, err
	}
	return registryClient, nil
}

// copied from https://github.com/helm/helm/blob/main/cmd/helm/root.go
func newDefaultRegistryClient(plainHTTP bool) (*registry.Client, error) {
	opts := []registry.ClientOption{
		registry.ClientOptDebug(false),
		registry.ClientOptEnableCache(true),
		registry.ClientOptWriter(os.Stderr),
		registry.ClientOptCredentialsFile(settings.RegistryConfig),
	}
	if plainHTTP {
		opts = append(opts, registry.ClientOptPlainHTTP())
	}

	// Create a new registry client
	registryClient, err := registry.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return registryClient, nil
}

// copied from https://github.com/helm/helm/blob/main/cmd/helm/root.go
func newRegistryClientWithTLS(certFile, keyFile, caFile string, insecureSkipTLSverify bool) (*registry.Client, error) {
	// Create a new registry client
	registryClient, err := registry.NewRegistryClientWithTLS(os.Stderr, certFile, keyFile, caFile, insecureSkipTLSverify,
		settings.RegistryConfig, settings.Debug,
	)
	if err != nil {
		return nil, err
	}
	return registryClient, nil
}

// run is a modified version of Helm Pull command's Run function
// run downloads the chart and untar it always for necessary image processing
// but whether the untar files are kept or not depends on the config of api.Config.IncludeChartFiles
func (cl *remotechartloader) run(p *action.Pull, rc *registry.Client, chartRef string, config api.Config) (string, string, error) {
	var out strings.Builder

	c := downloader.ChartDownloader{
		Out:     &out,
		Keyring: p.Keyring,
		Verify:  downloader.VerifyNever,
		Getters: getter.All(p.Settings),
		Options: []getter.Option{
			getter.WithBasicAuth(p.Username, p.Password),
			getter.WithPassCredentialsAll(p.PassCredentialsAll),
			getter.WithTLSClientConfig(p.CertFile, p.KeyFile, p.CaFile),
			getter.WithInsecureSkipVerifyTLS(p.InsecureSkipTLSverify),
			getter.WithPlainHTTP(p.PlainHTTP),
		},
		RegistryClient:   rc,
		RepositoryConfig: p.Settings.RepositoryConfig,
		RepositoryCache:  p.Settings.RepositoryCache,
	}

	if registry.IsOCI(chartRef) {
		c.Options = append(c.Options,
			getter.WithRegistryClient(rc))
		c.RegistryClient = rc
	}

	if p.Verify {
		c.Verify = downloader.VerifyAlways
	} else if p.VerifyLater {
		c.Verify = downloader.VerifyLater
	}

	if p.RepoURL != "" {
		chartURL, err := repo.FindChartInAuthAndTLSAndPassRepoURL(p.RepoURL, p.Username, p.Password, chartRef, p.Version, p.CertFile, p.KeyFile, p.CaFile, p.InsecureSkipTLSverify, p.PassCredentialsAll, getter.All(p.Settings))
		if err != nil {
			return out.String(), "", err
		}
		chartRef = chartURL
	}

	// always download to the configured folder
	saved, v, err := c.DownloadTo(chartRef, p.Version, p.DestDir)
	if err != nil {
		return out.String(), "", err
	}

	if p.Verify {
		for name := range v.SignedBy.Identities {
			fmt.Fprintf(&out, "Signed by: %v\n", name)
		}
		fmt.Fprintf(&out, "Using Key With Fingerprint: %X\n", v.SignedBy.PrimaryKey.Fingerprint)
		fmt.Fprintf(&out, "Chart Hash Verified: %s\n", v.FileHash)
	}

	// After verification, untar the chart into the requested directory.
	ud := p.UntarDir
	if !filepath.IsAbs(ud) {
		ud = filepath.Join(p.DestDir, ud)
	}
	// Let udCheck to check conflict file/dir without replacing ud when untarDir is the current directory(.).
	udCheck := ud
	if udCheck == "." {
		_, udCheck = filepath.Split(chartRef)
	} else {
		_, chartName := filepath.Split(chartRef)
		udCheck = filepath.Join(udCheck, chartName)
	}

	if _, err := os.Stat(udCheck); err != nil {
		if err := os.MkdirAll(udCheck, 0755); err != nil {
			return out.String(), "", errors.Wrap(err, "failed to untar (mkdir)")
		}
	} else {
		return out.String(), "", errors.Errorf("failed to untar: a file or directory with the name %s already exists", udCheck)
	}

	return out.String(), saved, expandFile(ud, saved)
}

// expandFile expands the src file into the dest directory.
func expandFile(dest, src string) error {
	h, err := os.Open(src)
	if err != nil {
		return err
	}
	defer h.Close()
	return expand(dest, h)
}

// expand uncompresses and extracts a chart into the specified directory.
func expand(dir string, r io.Reader) error {
	files, err := loader.LoadArchiveFiles(r)
	if err != nil {
		return err
	}

	// Copy all files verbatim. We don't parse these files because parsing can remove
	// comments.
	for _, file := range files {
		outpath, err := securejoin.SecureJoin(dir, file.Name)
		if err != nil {
			return err
		}

		// Make sure the necessary subdirs get created.
		basedir := filepath.Dir(outpath)
		if err := os.MkdirAll(basedir, 0755); err != nil {
			return err
		}

		if err := os.WriteFile(outpath, file.Data, 0644); err != nil {
			return err
		}
	}

	return nil
}

func (cl *remotechartloader) Finish(ctx context.Context, config api.Config) error {
	// we need to clean up if chart files are not included as they will be downloaded by default
	if !config.ChartFilesIncluded {
		for i := 0; i < len(cl.fromCharts); i++ {
			chartVer := strings.Split(cl.fromCharts[i], ":")
			chartName := chartVer[0]

			// remote ${toDIr}/{chartName}/chart/*
			os.RemoveAll(fmt.Sprintf("%s/%s/%s", cl.toDir, chartName, "chart"))
		}
	}

	return nil
}
