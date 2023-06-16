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

// config contains options for the exporter.
type config struct {
	convertor *ConvertorHolder
	// temporalitySelector metric.TemporalitySelector
	// aggregationSelector metric.AggregationSelector
	redactTimestamps bool
}

// newConfig creates a validated config configured with options.
//newConfig创建一个已验证的配置，配置有选项。
func newConfig(options ...Option) config {
	cfg := config{}
	for _, opt := range options {
		cfg = opt.apply(cfg)
	}

	// // // 好像是默认的编码，暂时屏蔽
	// // if cfg.convertor == nil {
	// // 	enc := json.NewEncoder(os.Stdout)
	// // 	enc.SetIndent("", "\t")
	// // 	cfg.convertor = &convertorHolder{convertor: enc}
	// // }

	// if cfg.temporalitySelector == nil {
	// 	cfg.temporalitySelector = metric.DefaultTemporalitySelector
	// }

	// if cfg.aggregationSelector == nil {
	// 	cfg.aggregationSelector = metric.DefaultAggregationSelector
	// }

	return cfg
}

// Option sets exporter option values.
type Option interface {
	apply(config) config
}

type optionFunc func(config) config

func (o optionFunc) apply(c config) config {
	return o(c)
}

// WithEncoder sets the exporter to use encoder to encode all the metric
// data-types to an output.
func WithConvertor(convertor Convertor) Option {
	return optionFunc(func(c config) config {
		if convertor != nil {
			c.convertor = &ConvertorHolder{convertor: convertor}
		}
		return c
	})
}

// // WithTemporalitySelector sets the TemporalitySelector the exporter will use
// // to determine the Temporality of an instrument based on its kind. If this
// // option is not used, the exporter will use the DefaultTemporalitySelector
// // from the go.opentelemetry.io/otel/sdk/metric package.
// func WithTemporalitySelector(selector metric.TemporalitySelector) Option {
// 	return temporalitySelectorOption{selector: selector}
// }

// type temporalitySelectorOption struct {
// 	selector metric.TemporalitySelector
// }

// func (t temporalitySelectorOption) apply(c config) config {
// 	c.temporalitySelector = t.selector
// 	return c
// }

// // WithAggregationSelector sets the AggregationSelector the exporter will use
// // to determine the aggregation to use for an instrument based on its kind. If
// // this option is not used, the exporter will use the
// // DefaultAggregationSelector from the go.opentelemetry.io/otel/sdk/metric
// // package or the aggregation explicitly passed for a view matching an
// // instrument.
// func WithAggregationSelector(selector metric.AggregationSelector) Option {
// 	// Deep copy and validate before using.
// 	wrapped := func(ik metric.InstrumentKind) aggregation.Aggregation {
// 		a := selector(ik)
// 		cpA := a.Copy()
// 		if err := cpA.Err(); err != nil {
// 			cpA = metric.DefaultAggregationSelector(ik)
// 			global.Error(
// 				err, "using default aggregation instead",
// 				"aggregation", a,
// 				"replacement", cpA,
// 			)
// 		}
// 		return cpA
// 	}

// 	return aggregationSelectorOption{selector: wrapped}
// }

// type aggregationSelectorOption struct {
// 	selector metric.AggregationSelector
// }

// func (t aggregationSelectorOption) apply(c config) config {
// 	c.aggregationSelector = t.selector
// 	return c
// }

// WithoutTimestamps sets all timestamps to zero in the output stream.
func WithoutTimestamps() Option {
	return optionFunc(func(c config) config {
		c.redactTimestamps = true
		return c
	})
}

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
