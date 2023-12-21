/*
Copyright The Helm Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package getter

import (
	"bytes"
	"fmt"
	"os"

	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"

	"github.com/Masterminds/vcs"
	securejoin "github.com/cyphar/filepath-securejoin"

	"helm.sh/helm/v3/internal/gitutil"
)

// GitGetter is the default git backend handler
type GitGetter struct {
	opts options
}

func (g *GitGetter) ChartName() string {
	return g.opts.chartName
}

// Get performs a Get from repo.Getter and returns the body.
func (g *GitGetter) Get(href string, options ...Option) (*bytes.Buffer, error) {
	for _, opt := range options {
		opt(&g.opts)
	}
	return g.get(href)
}

func (g *GitGetter) get(href string) (*bytes.Buffer, error) {
	gitURL, err := gitutil.ParseGitRepositoryURL(href)
	if err != nil {
		return nil, err
	}
	version := g.opts.version
	chartName := g.opts.chartName
	if version == "" {
		return nil, fmt.Errorf("the version must be a valid tag or branch name for the git repo, not nil")
	}
	tmpDir, err := os.MkdirTemp("", "helm-git-")
	if err != nil {
		return nil, err
	}

	gitTmpDir, err := securejoin.SecureJoin(tmpDir, chartName)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(gitTmpDir, 0755); err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	repo, err := vcs.NewRepo(gitURL.GitRemoteURL.String(), gitTmpDir)
	if err != nil {
		return nil, err
	}
	if err := repo.Get(); err != nil {
		return nil, err
	}
	if err := repo.UpdateVersion(version); err != nil {
		return nil, err
	}

	chartDir, err := securejoin.SecureJoin(gitTmpDir, gitURL.PathUnderGitRepository)
	if err != nil {
		return nil, err
	}

	ch, err := loader.LoadDir(chartDir)
	if err != nil {
		return nil, err
	}

	tarballPath, err := chartutil.Save(ch, tmpDir)
	if err != nil {
		return nil, err
	}

	buf, err := os.ReadFile(tarballPath)
	return bytes.NewBuffer(buf), err
}

// NewGitGetter constructs a valid git client as a Getter
func NewGitGetter(ops ...Option) (Getter, error) {

	client := GitGetter{}

	for _, opt := range ops {
		opt(&client.opts)
	}

	return &client, nil
}
