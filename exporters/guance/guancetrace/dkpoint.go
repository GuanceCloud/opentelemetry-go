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

package guancetrace // import "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"

import (
	"fmt"
	"strings"
	"time"

	"github.com/GuanceCloud/cliutils/point"
)

func buildPoints(dkSpans []*DatakitSpan) ([]*point.Point, error) {
	pts := make([]*point.Point, 0)

	for _, dkSpan := range dkSpans {
		processUnknown(dkSpan)

		tags := map[string]string{
			TAG_SERVICE:     dkSpan.Service,
			TAG_OPERATION:   dkSpan.Operation,
			TAG_SOURCE_TYPE: dkSpan.SourceType,
			TAG_SPAN_TYPE:   dkSpan.SpanType,
			TAG_SPAN_STATUS: dkSpan.Status,
		}

		for k, v := range dkSpan.Tags {
			tags[strings.ReplaceAll(k, ".", "_")] = v
		}

		fields := map[string]interface{}{
			FIELD_START:    dkSpan.Start / int64(time.Microsecond),
			FIELD_DURATION: dkSpan.Duration / int64(time.Microsecond),
		}

		stringFields := map[string]interface{}{
			FIELD_TRACEID:  dkSpan.TraceID,
			FIELD_PARENTID: dkSpan.ParentID,
			FIELD_SPANID:   dkSpan.SpanID,
			FIELD_RESOURCE: dkSpan.Resource,
			FIELD_MESSAGE:  dkSpan.Content,
		} // string 型字段，要额外传入 pt

		for k, v := range dkSpan.Metrics {
			if _, ok := v.(string); ok {
				stringFields[strings.ReplaceAll(k, ".", "_")] = v
			} else {
				fields[strings.ReplaceAll(k, ".", "_")] = v
			}
		}

		// Create point.
		opts := point.DefaultMetricOptions()
		pt := point.NewPointV2([]byte(dkSpan.Source), append(point.NewTags(tags), point.NewKVs(fields)...), opts...)
		kvs := point.NewKVs(stringFields)
		for _, v := range kvs {
			pt.AddKV(v)
		}

		fmt.Println("pt.LineProto() == ", pt.LineProto())
		pts = append(pts, pt)
	}

	return pts, nil
}

func processUnknown(dkSpan *DatakitSpan) {
	if dkSpan != nil {
		if dkSpan.Service == "" {
			dkSpan.Service = UNKNOWN_SERVICE
		}
		if dkSpan.SourceType == "" {
			dkSpan.SourceType = SPAN_SOURCE_CUSTOMER
		}
		if dkSpan.SpanType == "" {
			dkSpan.SpanType = SPAN_TYPE_UNKNOWN
		}
	}
}
