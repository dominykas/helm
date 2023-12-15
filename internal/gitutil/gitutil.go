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
	"os"
	"strings"

	"github.com/Masterminds/vcs"
)

// HasGitReference returns true if a git repository contains a specified ref (branch/tag)
func HasGitReference(gitRepo, ref, repoName string) (bool, error) {
	local, err := os.MkdirTemp("", repoName)
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
	return strings.HasPrefix(url, "git://")
}

// RepositoryURLToGitURL converts a repository URL into a URL that `git clone` could consume
func RepositoryURLToGitURL(url string) string {
	return strings.TrimPrefix(url, "git://")
}
