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

import (
	"encoding/json"
	"fmt"

	"github.com/GuanceCloud/cliutils/lineproto"
	// "github.com/prometheus/log"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// 来源于 internal/trace/trace.go
type DatakitSpan struct {
	TraceID    string                 `json:"trace_id"`
	ParentID   string                 `json:"parent_id"`
	SpanID     string                 `json:"span_id"`
	Service    string                 `json:"service"`     // service name
	Resource   string                 `json:"resource"`    // resource or api under service
	Operation  string                 `json:"operation"`   // api name
	Source     string                 `json:"source"`      // client tracer name
	SpanType   string                 `json:"span_type"`   // relative span position in tracing: entry, local, exit or unknow
	SourceType string                 `json:"source_type"` // service type
	Tags       map[string]string      `json:"tags"`
	Metrics    map[string]interface{} `json:"metrics"`
	Start      int64                  `json:"start"`    // unit: nano sec
	Duration   int64                  `json:"duration"` // unit: nano sec
	Status     string                 `json:"status"`   // span status like error, ok, info etc.
	Content    string                 `json:"content"`  // raw tracing data in json
}

func buildDKSpans(spanStubs []tracetest.SpanStub) ([]*DatakitSpan, error) {
	fmt.Println("这里进行转换==》")
	dkSpans := make([]*DatakitSpan, 0)
	spanIDs, parentIDs := getSpanIDsAndParentIDs(spanStubs)
	for _, spanStub := range spanStubs {

		dkSpan := &DatakitSpan{
			TraceID:   spanStub.SpanContext.TraceID().String(), // TraceID:   hex.EncodeToString(spanStub.GetTraceId()),
			ParentID:  spanStub.Parent.SpanID().String(),       // ParentID:  byteToString(span.GetParentSpanId()),
			SpanID:    spanStub.SpanContext.SpanID().String(),  // SpanID:    byteToString(span.GetSpanId()),
			Resource:  spanStub.Resource.String(),              // TODO 不保证对,//	Resource:  span.Name,
			Operation: spanStub.Name,                           //	Operation: span.Name,
			Source:    "otel-exporter",                         //	Source:    inputName,
			Tags:      make(map[string]string),
			Metrics:   make(map[string]interface{}),
			Start:     int64(spanStub.StartTime.Nanosecond()),                 //	Start:     int64(span.StartTimeUnixNano),
			Duration:  spanStub.EndTime.Sub(spanStub.StartTime).Nanoseconds(), //	Duration:  int64(span.EndTimeUnixNano - span.StartTimeUnixNano),
			Status:    spanStub.Status.Code.String(),                          //	TODO 不保证对, Status:    getDKSpanStatus(span.GetStatus()),
		}
		dkSpan.SpanType = findSpanTypeStrSpanID(dkSpan.SpanID, dkSpan.ParentID, spanIDs, parentIDs)

		attrs := spanStub.Attributes
		if kv, i := findAttr(attrs, otelResourceServiceKey); i != -1 {
			dkSpan.Service = kv.Value.AsString()
		}
		if kv, i := findAttr(attrs, otelResourceServiceVersionKey); i != -1 {
			dkSpan.Tags[TAG_VERSION] = kv.Value.AsString()
		}
		if kv, i := findAttr(attrs, otelResourceProcessIDKey); i != -1 {
			dkSpan.Tags[TAG_PID] = kv.Value.AsString()
		}
		if kv, i := findAttr(attrs, otelResourceContainerNameKey); i != -1 {
			dkSpan.Tags[TAG_CONTAINER_HOST] = kv.Value.AsString()
		}
		if kv, i := findAttr(attrs, otelHTTPMethodKey); i != -1 {
			dkSpan.Tags[TAG_HTTP_METHOD] = kv.Value.AsString()
			// attrs.remove(otelHTTPMethodKey)
			attrs = append(attrs[0:i], attrs[i+1:]...)
		}
		if kv, i := findAttr(attrs, otelHTTPStatusCodeKey); i != -1 {
			dkSpan.Tags[TAG_HTTP_STATUS_CODE] = kv.Value.AsString()
			// attrs.remove(otelHTTPStatusCodeKey)
			attrs = append(attrs[0:i], attrs[i+1:]...)
		}

		// 处理 Events
		for i := range spanStub.Events {
			if spanStub.Events[i].Name == ExceptionEventName {
				for o, d := range otelErrKeyToDkErrKey {
					if kv, ok := getAttribute(o, spanStub.Events[i].Attributes); ok {
						dkSpan.Metrics[d] = kv.Value.AsString()
					}
				}
				break
			}
		}

		// 处理 Attributes
		attrtags, attrfields := spliteAttrs(attrs)
		dkSpan.Tags = mergeTags(dkSpan.Tags, attrtags)
		dkSpan.Metrics = mergeFields(dkSpan.Metrics, attrfields)

		dkSpan.SourceType = getSourceType(dkSpan.Tags)

		if buf, err := json.Marshal(spanStub); err != nil {
			fmt.Println("", err)
		} else {
			dkSpan.Content = string(buf)
		}

		dkSpans = append(dkSpans, dkSpan)

		/*
			// ==============================
			fmt.Println(dkSpan)
			traceName := spanStub.Name
			// // if spanStub.SpanContext.HasTraceID() {
			// var traceID [16]byte
			// traceID = spanStub.SpanContext.TraceID()
			// fields["trace_id"] = traceID
			stringFields["trace_id"] = dkSpan.TraceID
			// }
			// if spanStub.SpanContext.HasSpanID() {
			stringFields["span_id"] = dkSpan.SpanID
			// }
			// if spanStub.Parent.HasSpanID() {
			if dkSpan.ParentID != "" {
				stringFields["parent_id"] = dkSpan.ParentID
			}
			// }
			// "span_type" DK文档这样说的，但是对不上茬口[entry, local, exit, unknown] vs ["internal","server","client","producer","consumer","unspecified"]
			tags["span_type"] = spanStub.SpanKind.String()
			fields["start"] = spanStub.StartTime.UnixNano()                             // TODO unit?: nano sec 这个默认纳秒吗？
			fields["duration"] = spanStub.EndTime.Sub(spanStub.StartTime).Nanoseconds() // TODO unit?: nano sec 这个需要强转一次吗？
			droppedAttributesCount := 0
			for _, attr := range spanStub.Attributes {
				droppedAttributesCount += addAttributes(string(attr.Key), attr.Value.AsString(), tags)
			}
			// Events                 []tracesdk.Event // TODO 好像不用弄了，丢弃，只保留丢弃条数
			// Links                  []tracesdk.Link // TODO 好像不用弄了，丢弃，只保留丢弃条数
			tags["status"] = spanStub.Status.Code.String() // 这个放在

			if spanStub.DroppedAttributes+droppedAttributesCount > 0 {
				fields["dropped_attributes_count"] = spanStub.DroppedAttributes + droppedAttributesCount
			}
			if spanStub.DroppedEvents > 0 {
				fields["dropped_events_count"] = spanStub.DroppedEvents
			}
			if spanStub.DroppedLinks > 0 {
				fields["dropped_links_count"] = spanStub.DroppedLinks
			}
			if spanStub.ChildSpanCount > 0 {
				fields["child_span_count"] = spanStub.ChildSpanCount // TODO 这个私自加的，DK文档没有
			}
			if len(spanStub.Events) > 0 {
				fields["events_count"] = len(spanStub.Events)
			}
			if len(spanStub.Links) > 0 {
				fields["links_count"] = len(spanStub.Links)
			}
			stringFields["resource"] = spanStub.Resource.String() // TODO 不保证对
			// ==============================
		*/

	}
	return dkSpans, nil
}

func getSpanIDsAndParentIDs(spans []tracetest.SpanStub) (map[string]bool, map[string]bool) {
	var (
		spanIDs   = make(map[string]bool)
		parentIDs = make(map[string]bool)
	)
	for _, span := range spans {
		spanIDs[span.SpanContext.SpanID().String()] = true
		parentIDs[span.Parent.SpanID().String()] = true
	}

	return spanIDs, parentIDs
}

func findAttr(attrs []attribute.KeyValue, key string) (*attribute.KeyValue, int) {
	for i := len(attrs) - 1; i >= 0; i-- {
		if string(attrs[i].Key) == key {
			return &attrs[i], i
		}
	}

	return nil, -1
}

func spliteAttrs(attrs []attribute.KeyValue) (map[string]string, map[string]interface{}) {
	tags := make(map[string]string)
	metrics := make(map[string]interface{})
	maxTagValueLen := lineproto.NewDefaultOption().MaxTagValueLen
	for _, v := range attrs {
		switch v.Value.Type() {
		case attribute.STRING:
			if s := v.Value.AsString(); len(s) > maxTagValueLen {
				metrics[string(v.Key)] = s
			} else {
				tags[string(v.Key)] = s
			}
		case attribute.FLOAT64:
			metrics[string(v.Key)] = v.Value.AsFloat64()
		case attribute.INT64:
			metrics[string(v.Key)] = v.Value.AsInt64()
		}
	}

	return tags, metrics
}

func getAttribute(key string, attributes []attribute.KeyValue) (*attribute.KeyValue, bool) {
	for _, attr := range attributes {
		if string(attr.Key) == key {
			return &attr, true
		}
	}

	return nil, false
}

func mergeTags(input ...map[string]string) map[string]string {
	tags := make(map[string]string)
	for i := range input {
		for k, v := range input[i] {
			tags[k] = v
		}
	}

	return tags
}

func mergeFields(input ...map[string]interface{}) map[string]interface{} {
	fields := make(map[string]interface{})
	for i := range input {
		for k, v := range input[i] {
			fields[k] = v
		}
	}

	return fields
}

func getSourceType(tags map[string]string) string {
	for key := range tags {
		switch key {
		case otelHTTPSchemeKey, otelHTTPMethodKey, otelRPCSystemKey:
			return SPAN_SOURCE_WEB
		case otelDBSystemKey:
			return SPAN_SOURCE_DB
		case otelMessagingSystemKey:
			return SPAN_SOURCE_MSGQUE
		}
	}

	return SPAN_SOURCE_CUSTOMER
}
