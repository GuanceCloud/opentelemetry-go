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

// import (
// 	"regexp"

// 	commonpb "gitlab.jiagouyun.com/cloudcare-tools/datakit/plugins/inputs/opentelemetry/compiled/v1/common"
// )

// type getAttributeFunc func(key string, attributes []*commonpb.KeyValue) (*commonpb.KeyValue, bool)

// func getAttr(key string, attributes []*commonpb.KeyValue) (*commonpb.KeyValue, bool) {
// 	for _, attr := range attributes {
// 		if attr.Key == key {
// 			return attr, true
// 		}
// 	}

// 	return nil, false
// }

// func getAttrWrapper(ignore []*regexp.Regexp) getAttributeFunc {
// 	if len(ignore) == 0 {
// 		return getAttr
// 	} else {
// 		return func(key string, attributes []*commonpb.KeyValue) (*commonpb.KeyValue, bool) {
// 			for _, rexp := range ignore {
// 				if rexp.MatchString(key) {
// 					return nil, false
// 				}
// 			}

// 			return getAttr(key, attributes)
// 		}
// 	}
// }

// type extractAttributesFunc func(src []*commonpb.KeyValue) (dest []*commonpb.KeyValue)

// func extractAttrs(src []*commonpb.KeyValue) (dest []*commonpb.KeyValue) {
// 	dest = append(dest, src...)

// 	return
// }

// func extractAttrsWrapper(ignore []*regexp.Regexp) extractAttributesFunc {
// 	if len(ignore) == 0 {
// 		return extractAttrs
// 	} else {
// 		return func(src []*commonpb.KeyValue) (dest []*commonpb.KeyValue) {
// 		NEXT_ATTR:
// 			for _, v := range src {
// 				for _, rexp := range ignore {
// 					if rexp.MatchString(v.Key) {
// 						continue NEXT_ATTR
// 					}
// 				}
// 				dest = append(dest, v)
// 			}

// 			return
// 		}
// 	}
// }

// func newAttributes(attrs []*commonpb.KeyValue) *attributes {
// 	a := &attributes{}
// 	a.attrs = append(a.attrs, attrs...)

// 	return a
// }

// type attributes struct {
// 	attrs []*commonpb.KeyValue
// }

// // nolint: deadcode,unused
// func (a *attributes) loop(proc func(i int, k string, v *commonpb.KeyValue) bool) {
// 	for i, v := range a.attrs {
// 		if !proc(i, v.Key, v) {
// 			break
// 		}
// 	}
// }

// func (a *attributes) merge(attrs ...*commonpb.KeyValue) *attributes {
// 	for _, v := range attrs {
// 		if _, i := a.find(v.Key); i != -1 {
// 			a.attrs[i] = v
// 		} else {
// 			a.attrs = append(a.attrs, v)
// 		}
// 	}

// 	return a
// }

// func (a *attributes) find(key string) (*commonpb.KeyValue, int) {
// 	for i := len(a.attrs) - 1; i >= 0; i-- {
// 		if a.attrs[i].Key == key {
// 			return a.attrs[i], i
// 		}
// 	}

// 	return nil, -1
// }

// func (a *attributes) remove(key string) *attributes {
// 	if _, i := a.find(key); i != -1 {
// 		a.attrs = append(a.attrs[:i], a.attrs[i+1:]...)
// 	}

// 	return a
// }

// func (a *attributes) splite() (map[string]string, map[string]interface{}) {
// 	tags := make(map[string]string)
// 	metrics := make(map[string]interface{})
// 	for _, v := range a.attrs {
// 		switch v.Value.Value.(type) {
// 		case *commonpb.AnyValue_BytesValue, *commonpb.AnyValue_StringValue:
// 			if s := v.Value.GetStringValue(); len(s) > point.MaxTagValueLen {
// 				metrics[v.Key] = s
// 			} else {
// 				tags[v.Key] = s
// 			}
// 		case *commonpb.AnyValue_DoubleValue:
// 			metrics[v.Key] = v.Value.GetDoubleValue()
// 		case *commonpb.AnyValue_IntValue:
// 			metrics[v.Key] = v.Value.GetIntValue()
// 		}
// 	}

// 	return tags, metrics
// }
