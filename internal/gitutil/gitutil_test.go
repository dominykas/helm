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
	"testing"
)

func TestIsGitUrl(t *testing.T) {
	// Test table: Given url, IsGitURL should return expect.
	tests := []struct {
		url    string
		expect bool
	}{
		{"oci://example.com/example/chart", false},
		{"git://example.com/example/chart", true},
		{"git://https://example.com/example/chart", true},
	}

	for _, test := range tests {
		if IsGitURL(test.url) != test.expect {
			t.Errorf("Expected %t for %s", test.expect, test.url)
		}
	}
}
