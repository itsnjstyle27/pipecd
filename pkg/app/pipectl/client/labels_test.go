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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseLabels(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name    string
		labels  []string
		want    map[string]string
		wantErr string
	}{
		{
			name: "no labels",
			want: map[string]string{},
		},
		{
			name:   "several labels",
			labels: []string{"env:prod", "team:backend"},
			want:   map[string]string{"env": "prod", "team": "backend"},
		},
		{
			name:   "a colon in the value",
			labels: []string{"image:nginx:1.27"},
			want:   map[string]string{"image": "nginx:1.27"},
		},
		{
			name:   "an empty value",
			labels: []string{"env:"},
			want:   map[string]string{"env": ""},
		},
		{
			name:    "no colon",
			labels:  []string{"env:prod", "team"},
			wantErr: `invalid label "team"`,
		},
		{
			name:    "the wrong separator",
			labels:  []string{"env=prod"},
			wantErr: `invalid label "env=prod"`,
		},
		{
			name:    "an empty key",
			labels:  []string{":prod"},
			wantErr: `invalid label ":prod"`,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseLabels(tc.labels)
			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
