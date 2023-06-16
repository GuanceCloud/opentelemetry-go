// Copyright The OpenTelemetry Authors
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

package guancetrace

import (
	"reflect"
	"testing"
)

func Test_addAttributes(t *testing.T) {
	type args struct {
		k    string
		v    string
		tags map[string]string
	}
	tests := []struct {
		name     string
		args     args
		want     int
		wantTags map[string]string
	}{
		{
			name: "normal_not_in_covertRule",
			args: args{
				k:    "foo",
				v:    "bar",
				tags: make(map[string]string),
			},
			want:     1,
			wantTags: make(map[string]string),
		},
		{
			name: "normal_in_covertRule",
			args: args{
				k:    "http_host",
				v:    "http:/127.0.0.1:1234",
				tags: make(map[string]string),
			},
			want:     0,
			wantTags: map[string]string{"http_host": "http:/127.0.0.1:1234"},
		},
		{
			name: "normal_in_covertRule_have_old_tags",
			args: args{
				k:    "http_host",
				v:    "http:/127.0.0.1:1234",
				tags: map[string]string{"pid": "1357"},
			},
			want:     0,
			wantTags: map[string]string{"http_host": "http:/127.0.0.1:1234", "pid": "1357"},
		},
		{
			name: "normal_in_covertRule_have_same_old_tags",
			args: args{
				k:    "http_host",
				v:    "http:/127.0.0.1:1234",
				tags: map[string]string{"http_host": "http:/127.0.0.1:2468"},
			},
			want:     0,
			wantTags: map[string]string{"http_host": "http:/127.0.0.1:1234"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := addAttributes(tt.args.k, tt.args.v, tt.args.tags)
			if got != tt.want {
				t.Errorf("addAttributes() got %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(tt.args.tags, tt.wantTags) {
				t.Errorf("addAttributes() got %v, want %v", tt.args.tags, tt.wantTags)
			}
		})
	}
}
