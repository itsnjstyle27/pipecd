// Copyright 2024 The PipeCD Authors.
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

package notifier

import (
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func Test_getAccountsAsString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		accounts []string
		want     string
	}{
		{
			name:     "empty",
			accounts: []string{},
			want:     "",
		},
		{
			name:     "single",
			accounts: []string{"foo"},
			want:     "<@foo>",
		},
		{
			name:     "multiple",
			accounts: []string{"foo", "bar"},
			want:     "<@foo> <@bar>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := getAccountsAsString(tt.accounts)
			if got != tt.want {
				t.Errorf("getAccountsAsString(): got %s, want %s", got, tt.want)
			}
		})
	}
}

func Test_getGroupsAsString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		groups []string
		want   string
	}{
		{
			name:   "empty",
			groups: []string{},
			want:   "",
		},
		{
			name:   "single",
			groups: []string{"foo"},
			want:   "<!subteam^foo>",
		},
		{
			name:   "with correct format <!subteam^foo>",
			groups: []string{"<!subteam^foo>"},
			want:   "<!subteam^foo>",
		},
		{
			name:   "multiple",
			groups: []string{"foo", "bar"},
			want:   "<!subteam^foo> <!subteam^bar>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := getGroupsAsString(tt.groups)
			if got != tt.want {
				t.Errorf("getGroupsAsString(): got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestTruncateText(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name string
		text string
		max  int
		want string
	}{
		{
			name: "shorter than max",
			text: "abc",
			max:  8,
			want: "abc",
		},
		{
			name: "exactly max",
			text: "abcdefgh",
			max:  8,
			want: "abcdefgh",
		},
		{
			name: "longer than max",
			text: "abcdefghi",
			max:  8,
			want: "abcdefgh...",
		},
		{
			name: "multi-byte text that fits in max is kept",
			text: "\u65e5\u672c\u8a9e",
			max:  8,
			want: "\u65e5\u672c\u8a9e",
		},
		{
			name: "multi-byte text is cut between characters",
			text: "\u65e5\u672c\u8a9e\u306e\u30b3\u30df\u30c3\u30c8\u30e1\u30c3\u30bb\u30fc\u30b8",
			max:  8,
			want: "\u65e5\u672c\u8a9e\u306e\u30b3\u30df\u30c3\u30c8...",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := truncateText(tc.text, tc.max)
			assert.Equal(t, tc.want, got)
			assert.True(t, utf8.ValidString(got))
		})
	}
}
