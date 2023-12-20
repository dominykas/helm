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

package gitutil

import (
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	giturls "github.com/whilp/git-urls"

	"github.com/Masterminds/vcs"
)

var gitRepositoryURLRe = regexp.MustCompile(`^git(\+\w+)?://`)

type GitRepositoryURL struct {
	RepositoryURL string
	GitRemoteURL  *url.URL
}

// HasGitReference returns true if a git repository contains a specified ref (branch/tag)
func HasGitReference(gitRepo, ref string) (bool, error) {
	local, err := os.MkdirTemp("", "helm-git-")
	if err != nil {
		return false, err
	}
	repo, err := vcs.NewRepo(gitRepo, local)

	if err != nil {
		return false, err
	}

	if err := repo.Get(); err != nil {
		return false, err
	}
	defer os.RemoveAll(local)
	return repo.IsReference(ref), nil
}

// IsGitRepository determines whether a URL is to be treated as a git repository URL
func IsGitRepository(url string) bool {
	return gitRepositoryURLRe.MatchString(url)
}

// ParseGitRepositoryURL creates a new GitRepositoryURL from a string
func ParseGitRepositoryURL(repositoryURL string) (*GitRepositoryURL, error) {
	gitRemoteURL, err := giturls.Parse(strings.TrimPrefix(repositoryURL, "git+"))

	if err != nil {
		return nil, err
	}

	if gitRemoteURL.User != nil {
		return nil, errors.Errorf("git repository URL should not contain credentials - please use git credential helpers")
	}

	return &GitRepositoryURL{
		RepositoryURL: repositoryURL,
		GitRemoteURL:  gitRemoteURL,
	}, err
}
