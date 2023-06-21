// Copyright The OpenTelemetry Authors
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

package guancetrace // import "go.opentelemetry.io/otel/exporters/guance/guancetrace "

// "go.opentelemetry.io/otel/sdk/metric"
// "go.opentelemetry.io/otel/sdk/metric/aggregation"

type spanInfo struct {
	name      string
	fieldType string
	isTag     bool
}

var covertRule map[string]spanInfo = map[string]spanInfo{
	"container_host":   spanInfo{"container_host", "", true},
	"endpoint":         spanInfo{"endpoint", "", true},
	"env":              spanInfo{"env", "", true},
	"http_host":        spanInfo{"http_host", "", true},
	"http_method":      spanInfo{"http_method", "", true},
	"http_route":       spanInfo{"http_route", "", true},
	"http_status_code": spanInfo{"http_status_code", "", true},
	"http_url":         spanInfo{"http_url", "", true},
	"operation":        spanInfo{"operation", "", true},
	"pid":              spanInfo{"pid", "", true},
	"project":          spanInfo{"project", "", true},
	"service":          spanInfo{"service", "", true},
	"source_type":      spanInfo{"source_type", "", true},
	"status":           spanInfo{"status", "", true},
	"span_type":        spanInfo{"span_type", "", true},
	"duration":         spanInfo{"duration", "Field_I", false},
	"message":          spanInfo{"message", "Field_I", false},
	"parent_id":        spanInfo{"parent_id", "Field_I", false},
	"priority":         spanInfo{"priority", "Field_I", false},
	"resource":         spanInfo{"resource", "Field_I", false},
	"sample_rate":      spanInfo{"sample_rate", "Field_I", false},
	"span_id":          spanInfo{"span_id", "Field_I", false},
	"start":            spanInfo{"start", "Field_I", false},
	"trace_id":         spanInfo{"trace_id", "Field_I", false},
}
