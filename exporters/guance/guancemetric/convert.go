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

package guancemetric // import "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"

import (
	"fmt"

	"github.com/GuanceCloud/cliutils/point"

	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

var scopeInfoKeys = [2]string{"scope_name", "scope_version"}

// Encoder encodes and outputs OpenTelemetry metric data-types as human readable text.
//编码器将OpenTelemetry度量数据类型编码并输出为人类可读文本。
type Convertor interface {
	// Encode handles the encoding and writing of OpenTelemetry metric data.
	//Encode处理OpenTelemetry度量数据的编码和写入。
	Convert(v *metricdata.ResourceMetrics) ([]*point.Point, error)
}

// encoderHolder is the concrete type used to wrap an Encoder so it can be used as a atomic.Value type.
//encoderHolder是用于包装编码器的具体类型，因此它可以用作原子.Value类型。
type ConvertorHolder struct {
	convertor Convertor
}

// Convert otelMetrics to lineProtos.
func (c ConvertorHolder) Convert(otelMetrics *metricdata.ResourceMetrics) ([]*point.Point, error) {
	fmt.Println("这里进行转换")
	points := make([]*point.Point, 0)

	for _, scopeMetric := range otelMetrics.ScopeMetrics {
		for _, m := range scopeMetric.Metrics {
			switch v := m.Data.(type) {
			case metricdata.Histogram[int64]:
				points = append(points, addHistogramMetric(v, m, scopeMetric.Scope.Name, scopeMetric.Scope.Version)...)
			case metricdata.Histogram[float64]:
				points = append(points, addHistogramMetric(v, m, scopeMetric.Scope.Name, scopeMetric.Scope.Version)...)
			case metricdata.Sum[int64]:
				points = append(points, addSumMetric(v, m, scopeMetric.Scope.Name, scopeMetric.Scope.Version)...)
			case metricdata.Sum[float64]:
				points = append(points, addSumMetric(v, m, scopeMetric.Scope.Name, scopeMetric.Scope.Version)...)
			case metricdata.Gauge[int64]:
				points = append(points, addGaugeMetric(v, m, scopeMetric.Scope.Name, scopeMetric.Scope.Version)...)
			case metricdata.Gauge[float64]:
				points = append(points, addGaugeMetric(v, m, scopeMetric.Scope.Name, scopeMetric.Scope.Version)...)
			}
		}
	}

	return points, nil
}

// Create a new point.
func newPoint(name string, tags map[string]string, fields map[string]interface{}) *point.Point {
	opts := point.DefaultMetricOptions()
	return point.NewPointV2([]byte(name),
		append(point.NewTags(tags), point.NewKVs(fields)...),
		opts...)
}

func addHistogramMetric[N int64 | float64](histogram metricdata.Histogram[N], m metricdata.Metrics, scopeName, scopeVersion string) []*point.Point {
	points := make([]*point.Point, 0)
	for _, dp := range histogram.DataPoints {
		if len(dp.Bounds)+1 != len(dp.BucketCounts) {
			return make([]*point.Point, 0) // TODO 这里出错处理要斟酌，几乎不会出错的
		}

		tags := make(map[string]string)
		fields := make(map[string]interface{})

		// Add tags.
		tags["scope_version"] = scopeVersion
		// tags["description"] = m.Description // 这个不想要了
		tags["unit"] = m.Unit
		kvs := dp.Attributes.ToSlice()
		for _, kv := range kvs {
			tags[string(kv.Key)] = kv.Value.AsString()
		}

		// Add bucket points.
		var bound string
		for i := 0; i < len(dp.Bounds)+1; i++ {
			// Add tags.
			if i == len(dp.Bounds) {
				bound = "+Inf"
			} else {
				bound = fmt.Sprintf("%f", dp.Bounds[i])
			}
			tags["le"] = bound // tags["le"] value will be overwritten time and time

			// Add fields.
			fields = make(map[string]interface{})
			fields[m.Name+"_bucket"] = dp.BucketCounts[i]

			// Create point.
			pt := newPoint(scopeName, tags, fields)
			pt.SetTime(dp.Time)

			points = append(points, pt)
		}
		delete(tags, "le")

		// Add sum points.
		fields = make(map[string]interface{})
		fields[m.Name+"_sum"] = dp.Sum
		pt := newPoint(scopeName, tags, fields)
		pt.SetTime(dp.Time)
		points = append(points, pt)

		// Add count points.
		fields = make(map[string]interface{})
		fields[m.Name+"_count"] = dp.Count
		pt = newPoint(scopeName, tags, fields)
		pt.SetTime(dp.Time)
		points = append(points, pt)
	}

	return points
}

func addSumMetric[N int64 | float64](sum metricdata.Sum[N], m metricdata.Metrics, scopeName, scopeVersion string) []*point.Point {
	points := make([]*point.Point, 0)
	for _, dp := range sum.DataPoints {
		tags := make(map[string]string)
		fields := make(map[string]interface{})

		// Add tags.
		tags["scope_version"] = scopeVersion
		// tags["description"] = m.Description // 这个不想要了
		tags["unit"] = m.Unit
		kvs := dp.Attributes.ToSlice()
		for _, kv := range kvs {
			tags[string(kv.Key)] = kv.Value.AsString()
		}

		// Add fields.
		fields[m.Name] = dp.Value

		// Create point.
		pt := newPoint(scopeName, tags, fields)
		pt.SetTime(dp.Time)

		points = append(points, pt)
	}

	return points
}

func addGaugeMetric[N int64 | float64](gauge metricdata.Gauge[N], m metricdata.Metrics, scopeName, scopeVersion string) []*point.Point {
	points := make([]*point.Point, 0)
	for _, dp := range gauge.DataPoints {
		tags := make(map[string]string)
		fields := make(map[string]interface{})

		// Add tags.
		tags["scope_version"] = scopeVersion
		// tags["description"] = m.Description // 这个不想要了
		tags["unit"] = m.Unit
		kvs := dp.Attributes.ToSlice()
		for _, kv := range kvs {
			tags[string(kv.Key)] = kv.Value.AsString()
		}

		// Add fields.
		fields[m.Name] = dp.Value

		// Create point.
		pt := newPoint(scopeName, tags, fields)
		pt.SetTime(dp.Time)

		points = append(points, pt)
	}

	return points
}
