# Helm Packager

Helm Packager (`helm-packager`) is a tool for packaging Helm charts and their images.

In the real world, there are some common scenarios like:

1. The Helm charts might be embedded in the custom CLI with `embed.FS`.
2. The Helm charts might be located somewhere in Internet.

Now, we want to:

1. List out all images that a Helm chart uses.
2. Download and export the Helm chart and its images to local folder so that we can install in the air-gapped env.
3. Import the downloaded Helm chart and its images into Helm chart artifactory and OCI-compliant registry accordingly.


## CLI Usage

### Pull

Pull command pulls the Helm charts and their images to local.

**Usage:**

```sh
helm-packager pull \
  --from-chart-repo <REMOTE_REPOSITORY_URL> \
  --from-charts <CHART_NAME>[:<CHART_VERSION>][,<CHART_NAME>[:<CHART_VERSION>]] \
 [--to-dir <CHARTS_DIR>]
```

Note:
- The chart version is optional. When no version is specified, the latest version will be used.
- When no `--to-dir` is specified, the output will be printed to `stdout` so it's convenient when you want to have a peak at what the Helm chart images are.

For example:

```sh
helm-packager pull \
  --from-chart-repo oci://registry-1.docker.io/bitnamicharts \
  --from-charts apache:10.2.3,nginx
```

OUTPUT:

```
.
├── apache
│   ├── apache-10.2.3.tgz
│   └── images
│       └── apache-2.4.58-debian-11-r1.tar (docker.io/bitnami/apache:2.4.58-debian-11-r1)
└── nginx
    ├── nginx-15.4.4.tgz
    └── images
        └── nginx-1.25.3-debian-11-r1.tar (docker.io/bitnami/nginx:1.25.3-debian-11-r1)
```

When `--to-dir` is specified, the output will be writen to the directory with a structured folders and files:

```sh
helm-packager pull \
  --from-chart-repo oci://registry-1.docker.io/bitnamicharts \
  --from-charts apache:10.2.3,nginx \
  --to-dir ./_charts
```

Same output but now the charts and their images are downloaded, which can be verified by `tree` command:

```console
$ tree _charts
_charts
├── apache
│   ├── apache-10.2.3.tgz
│   └── images
│       └── apache-2.4.58-debian-11-r1.tar
└── nginx
    ├── images
    │   └── nginx-1.25.3-debian-11-r1.tar
    └── nginx-15.4.4.tgz

4 directories, 4 files
```

### Push

**...WIP...**

Push command pushs the exported local Helm charts and their images to remote Helm repository (e.g. [ChartMuseum](https://github.com/helm/chartmuseum)) and image registry.

**Usage:**

```sh
helm-packager push \
  --from-dir <EXPORTED DIR WITH CHARTS AND IMAGES> \
 [--from-charts <CHART_NAME>[,<CHART_NAME>] \
  --to-chart-repo <TARGETED HELM REPOSITORY TO PUSH CHARTS TO> \
  --to-image-registry <TARGETED IMAGE REGISTRY TO PUSH IMAGES TO>
```

For example, to push all charts and their images witin a specified `./_charts` folder:

```sh
helm-packager push \
  --from-dir ./_charts \
  --to-chart-repo https://my.chart.repo \
  --to-image-resitry https://my.docker.registry
```

Or to selectively push one or some charts and their images witin a specified `./_charts` folder:

```sh
helm-packager push \
  --from-dir ./_charts \
  --from-charts apache,nginx \
  --to-chart-repo https://my.chart.repo \
  --to-image-resitry https://my.docker.registry
```

### Copy

**...WIP...**

Copy command is to copy Helm charts and their images from source Helm chart repository / image registry to target Helm chart repository / image registry.

**Usage:**

```sh
helm-packager copy \
  --from-chart-repo <REMOTE_REPOSITORY_URL> \
  --from-charts <CHART_NAME>[:<CHART_VERSION>][,<CHART_NAME>[:<CHART_VERSION>]] \
  --to-chart-repo <TARGETED HELM REPOSITORY TO PUSH CHARTS TO> \
  --to-image-registry <TARGETED IMAGE REGISTRY TO PUSH IMAGES TO>
```

For example, to copy Helm charts `apache` with specific version `10.2.3` and another Helm chart `nginx` from Bitnami repository to the private Helm chart repository / image registry.

```sh
helm-packager copy \
  --from-chart-repo oci://registry-1.docker.io/bitnamicharts \
  --from-charts apache:10.2.3,nginx \
  --to-chart-repo https://my.chart.repo \
  --to-image-resitry https://my.docker.registry
```


## As an SDK

There are a few examples provided under [examples](./examples/) folder.

### Example: Process Charts from `embed.FS`

```go
package main

import (...)

//go:embed charts
var embeddedCharts embed.FS

var chartName = flag.String("chart", "robotshop", "Chart name.")

func main() {
	flag.Parse()

	chartFS, err := fs.Sub(embeddedCharts, fmt.Sprintf("charts/%s", *chartName))
	if err != nil {
		panic(err)
	}

	cl := chartloader.NewEmbedChartLoader(chartFS)
	cw := chartwriter.NewStdoutChartWriter()
	iw := imageswriter.NewStdoutImagesWriter()

	ctx := context.Background()

	cp := pipeline.NewBuilder(ctx).
		WithChartLoader(cl).
		WithChartWriter(cw).
		WithImagesWriter(iw).
		ConfigureChartFilesIncluded(false).
		Complete()

	err = cp.Process()
	if err != nil {
		panic(err)
	}
}
```

### Example: Process Helm charts from remote repository

```go
package main

import (...)

var (
	fromRepo   = "oci://registry-1.docker.io/bitnamicharts"
	fromCharts = []string{"apache:10.2.3", "nginx"}
	toDir      = "./_charts"
)

func main() {
	ctx := context.Background()

	// create folder if needed
	if _, err := os.Stat(toDir); os.IsNotExist(err) {
		err := os.Mkdir(toDir, 0766)
		if err != nil {
			panic(err)
		}
	}

	cl := chartloader.NewRemoteChartLoader(fromRepo, fromCharts, toDir)
	cw := chartwriter.NewStdoutChartWriter()
	iw := imageswriter.NewFileImagesWriter(toDir)

	cp := pipeline.NewBuilder(ctx).
		WithChartLoader(cl).
		WithChartWriter(cw).
		WithImagesWriter(iw).
		ConfigureChartFilesIncluded(false).
		Complete()

	err := cp.Process()
	if err != nil {
		panic(err)
	}
}
```