// Copyright © 2023 Bright Zheng <bright.zheng@outlook.com>
// SPDX-License-Identifier: Apache-2.0

package imageswriter

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/brightzheng100/helm-packager/pkg/api"
	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"gopkg.in/op/go-logging.v1"
	"helm.sh/helm/v3/pkg/action"
)

func Templatize(ctx context.Context, chart *api.Chart) (string, error) {
	// Create chart renderer.
	client := action.NewInstall(&action.Configuration{})
	client.ClientOnly = true
	client.DryRun = true
	client.ReleaseName = chart.C.Name()
	client.IncludeCRDs = false
	client.Namespace = "fake-namespace-name"

	// Render chart.
	rel, err := client.Run(chart.C, nil)
	if err != nil {
		return "", fmt.Errorf("could not render helm chart correctly: %w", err)
	}

	return rel.Manifest, nil
}

// extractImages extracts the images from the templatized Helm chart
// Ref: https://mikeperry.io/posts/copy-helm-images/
// helm template . \
// | yq '..|.image? | select(.)' \
// | sort \
// | uniq \
// | xargs -I % -n 1 -P 4 bash -c "docker pull % && docker tag % my.registry.com/% docker push my.registry.com/%"
func ExtractImages(ctx context.Context, yaml string) ([]string, error) {
	images := []string{}

	// yqlib logging
	backendLeveled := logging.AddModuleLevel(logging.NewLogBackend(os.Stdout, "", 0))
	backendLeveled.SetLevel(logging.WARNING, "")
	yqlib.GetLogger().SetBackend(backendLeveled)

	encoder := yqlib.NewYamlEncoder(2, false, yqlib.ConfiguredYamlPreferences)
	decoder := yqlib.NewYamlDecoder(yqlib.ConfiguredYamlPreferences)

	result, err := yqlib.NewStringEvaluator().Evaluate(expression, yaml, encoder, decoder)
	if err != nil {
		return images, fmt.Errorf("could not evaluate with yq: %w", err)
	}

	//fmt.Println(result)
	scanner := bufio.NewScanner(strings.NewReader(result))
	for scanner.Scan() {
		s := scanner.Text()
		if exists := slices.Contains(images, s); !exists && r.MatchString(s) {
			images = append(images, s)
		}
	}
	slices.Sort(images)

	return images, nil
}
