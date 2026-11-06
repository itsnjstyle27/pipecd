// Copyright 2026 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"fmt"
	"strings"
)

// ParseLabels converts KEY:VALUE flag values into a map of labels.
// A value in another form is an error rather than being skipped, because the
// labels are filters: without one, a command would act on every application.
func ParseLabels(labels []string) (map[string]string, error) {
	parsed := make(map[string]string, len(labels))
	for _, label := range labels {
		key, value, ok := strings.Cut(label, ":")
		if !ok || key == "" {
			return nil, fmt.Errorf("invalid label %q, expected the form KEY:VALUE", label)
		}
		parsed[key] = value
	}
	return parsed, nil
}
