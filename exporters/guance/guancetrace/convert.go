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

package guancetrace // import "go.opentelemetry.io/otel/exporters/guance/guancetrace"

import (
	"github.com/GuanceCloud/cliutils/point"

	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

var scopeInfoKeys = [2]string{"scope_name", "scope_version"}

// Encoder encodes and outputs OpenTelemetry metric data-types as human readable text.
//编码器将OpenTelemetry度量数据类型编码并输出为人类可读文本。
type Convertor interface {
	// Encode handles the encoding and writing of OpenTelemetry metric data.
	//Encode处理OpenTelemetry度量数据的编码和写入。
	Convert(spanStubs []tracetest.SpanStub) ([]*point.Point, error)
}

// encoderHolder is the concrete type used to wrap an Encoder so it can be used as a atomic.Value type.
//encoderHolder是用于包装编码器的具体类型，因此它可以用作原子.Value类型。
type ConvertorHolder struct {
	convertor Convertor
}

// type tempSpan struct {
// 	TraceID    [16]byte               `json:"trace_id"`
// 	ParentID   [8]byte                `json:"parent_id"`
// 	SpanID     [8]byte                `json:"span_id"`
// 	Service    string                 `json:"service"`     // service name
// 	Resource   string                 `json:"resource"`    // resource or api under service
// 	Operation  string                 `json:"operation"`   // api name
// 	Source     string                 `json:"source"`      // client tracer name
// 	SpanType   string                 `json:"span_type"`   // relative span position in tracing: entry, local, exit or unknow
// 	SourceType string                 `json:"source_type"` // service type
// 	Tags       map[string]string      `json:"tags"`
// 	Metrics    map[string]interface{} `json:"metrics"`
// 	Start      int64                  `json:"start"`    // unit: nano sec
// 	Duration   int64                  `json:"duration"` // unit: nano sec
// 	Status     string                 `json:"status"`   // span status like error, ok, info etc.
// 	Content    string                 `json:"content"`  // raw tracing data in json
// }

// Convert otelTrace to lineProtos.
func (c ConvertorHolder) Convert(spanStubs []tracetest.SpanStub) ([]*point.Point, error) {

	dkSpans, err := buildDKSpans(spanStubs)
	if err != nil {
		return nil, err
	}

	pts, err := buildPoints(dkSpans)
	if err != nil {
		return nil, err
	}

	return pts, nil
}

// addAttributes add attributes info as tags.
func addAttributes(k, v string, tags map[string]string) int {
	sInfo, ok := covertRule[k]
	if !ok || !sInfo.isTag {
		return 1
	}
	tags[k] = v

	return 0
}

// newPoint create a new point.
func newPoint(name string, tags map[string]string, fields map[string]interface{}) *point.Point {
	opts := point.DefaultMetricOptions()
	return point.NewPointV2([]byte(name),
		append(point.NewTags(tags), point.NewKVs(fields)...),
		opts...)
}

func findSpanTypeStrSpanID(spanID, parentID string, spanIDs, parentIDs map[string]bool) string {
	if parentID != "0" && parentID != "" {
		if spanIDs[parentID] {
			if parentIDs[spanID] {
				return SPAN_TYPE_LOCAL
			} else {
				return SPAN_TYPE_EXIT
			}
		}
	}

	return SPAN_TYPE_ENTRY
}
