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

package guancemetric // import "go.opentelemetry.io/otel/exporters/guance/guancemetric"

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/GuanceCloud/cliutils/point"

	"go.opentelemetry.io/otel/exporters/guance/internal/feed"
	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

var errShutdown = errors.New("exporter shutdown")

// exporter is an OpenTelemetry metric exporter.
type exporter struct {
	host    string
	token   string
	pointCh chan []*point.Point
	encVal  atomic.Value // encoderHolder
	// wg         sync.WaitGroup // Shutdown wait deed return
	stopped atomic.Bool

	temporalitySelector metric.TemporalitySelector
	aggregationSelector metric.AggregationSelector

	redactTimestamps bool
}

// Like : exporters/stdout/stdoutmetric/exporter.go
// New returns a configured metric exporter.
//
// If no options are passed, the default exporter returned will use a JSON
// encoder with tab indentations that output to STDOUT.
func New(host, token string, options ...Option) (metric.Exporter, error) {
	cfg := newConfig(options...)

	exp := &exporter{
		host:  host,
		token: token,
		// pointCh:             make(chan []*point.Point, 1),
		temporalitySelector: cfg.temporalitySelector,
		aggregationSelector: cfg.aggregationSelector,
		redactTimestamps:    cfg.redactTimestamps,
	}
	exp.encVal.Store(*cfg.convertor)

	return exp, nil
}

// Temporality returns the Temporality to use for an instrument kind.
//Temporality返回用于乐器种类的临时性。
func (e *exporter) Temporality(k metric.InstrumentKind) metricdata.Temporality {
	return e.temporalitySelector(k)
}

// Aggregation returns the Aggregation to use for an instrument kind.
//Aggregation返回要用于工具类型的Aggregation。
func (e *exporter) Aggregation(k metric.InstrumentKind) aggregation.Aggregation {
	return e.aggregationSelector(k)
}

// Export serializes and transmits metric data to a internal/feed.
//导出串行化度量数据并将其传输到 internal/feed 。
// This is called synchronously, there is no concurrency safety requirement. Because of this, it is critical that all timeouts and cancellations of the passed context be honored.
//这被称为同步，没有并发安全要求。因此，遵守传递的上下文的所有超时和取消是至关重要的。
// All retry logic must be contained in this function. The SDK does not implement any retry logic. All errors returned by this function are considered unrecoverable and will be reported to a configured error Handler.
//此函数中必须包含所有重试逻辑。SDK没有实现任何重试逻辑。此函数返回的所有错误都被认为是不可恢复的，并将报告给配置的错误处理程序。
// The passed ResourceMetrics may be reused when the call completes. If an exporter needs to hold this data after it returns, it needs to make a copy.
//当调用完成时，可以重用传递的ResourceMetrics。如果出口商在返回后需要保存这些数据，则需要制作一份副本。
func (e *exporter) Export(ctx context.Context, data *metricdata.ResourceMetrics) error {
	if e.stopped.Load() {
		return errShutdown
	}

	select {
	case <-ctx.Done():
		// Don't do anything if the context has already timed out.
		return ctx.Err()
	default:
		// Context is still valid, continue.
	}

	if e.redactTimestamps {
		redactTimestamps(data)
	}

	convertorHolder := e.encVal.Load().(ConvertorHolder)
	points, err := convertorHolder.Convert(data)
	if err != nil {
		return err
	}

	urlStr := fmt.Sprintf("%s/v1/write/%s?token=%s", e.host, "metric", e.token)
	if e.token == "" {
		urlStr = fmt.Sprintf("%s/v1/write/%s", e.host, "metric")
	}

	feedInfos := make([]feed.FeedInfo, 0)
	for _, pt := range points {
		feedInfos = append(feedInfos, feed.FeedInfo{pt.LineProto() + "\n", urlStr})
	}
	fmt.Println("Export发送给chan")
	feed.FeedCh <- feedInfos

	return nil
}

// ForceFlush flushes any metric data held by an exporter.
//ForceFlush刷新导出程序所持有的所有度量数据。
// The deadline or cancellation of the passed context must be honored. An appropriate error should be returned in these situations.
//必须遵守截止日期或取消已传递的上下文。在这些情况下，应该返回适当的错误。
func (e *exporter) ForceFlush(ctx context.Context) error {
	if e.stopped.Load() {
		return errShutdown
	}

	// exporter holds no state, nothing to flush.
	return ctx.Err()
}

// Shutdown flushes all metric data held by an exporter and releases any held computational resources.
//关闭将刷新导出程序所持有的所有度量数据，并释放所有持有的计算资源。
// The deadline or cancellation of the passed context must be honored. An appropriate error should be returned in these situations.
//必须遵守截止日期或取消已传递的上下文。在这些情况下，应该返回适当的错误。
// After Shutdown is called, calls to Export will perform no operation and instead will return an error indicating the shutdown state.
//调用Shutdown后，对Export的调用将不执行任何操作，而是返回一个指示关闭状态的错误。
func (e *exporter) Shutdown(ctx context.Context) error {
	if e.stopped.Load() {
		return errShutdown
	}
	e.stopped.Swap(true) // Set exporter shutdown

	return ctx.Err()
}

// 改时间戳
func redactTimestamps(orig *metricdata.ResourceMetrics) {
	for i, sm := range orig.ScopeMetrics {
		metrics := sm.Metrics
		for j, m := range metrics {
			data := m.Data
			orig.ScopeMetrics[i].Metrics[j].Data = redactAggregationTimestamps(data)
		}
	}
}

var (
	errUnknownAggType = errors.New("unknown aggregation type")
)

// 改时间戳
func redactAggregationTimestamps(orig metricdata.Aggregation) metricdata.Aggregation {
	switch a := orig.(type) {
	case metricdata.Sum[float64]:
		return metricdata.Sum[float64]{
			Temporality: a.Temporality,
			DataPoints:  redactDataPointTimestamps(a.DataPoints),
			IsMonotonic: a.IsMonotonic,
		}
	case metricdata.Sum[int64]:
		return metricdata.Sum[int64]{
			Temporality: a.Temporality,
			DataPoints:  redactDataPointTimestamps(a.DataPoints),
			IsMonotonic: a.IsMonotonic,
		}
	case metricdata.Gauge[float64]:
		return metricdata.Gauge[float64]{
			DataPoints: redactDataPointTimestamps(a.DataPoints),
		}
	case metricdata.Gauge[int64]:
		return metricdata.Gauge[int64]{
			DataPoints: redactDataPointTimestamps(a.DataPoints),
		}
	case metricdata.Histogram[int64]:
		return metricdata.Histogram[int64]{
			Temporality: a.Temporality,
			DataPoints:  redactHistogramTimestamps(a.DataPoints),
		}
	case metricdata.Histogram[float64]:
		return metricdata.Histogram[float64]{
			Temporality: a.Temporality,
			DataPoints:  redactHistogramTimestamps(a.DataPoints),
		}
	default:
		global.Error(errUnknownAggType, fmt.Sprintf("%T", a))
		return orig
	}
}

// 改时间戳
func redactHistogramTimestamps[T int64 | float64](hdp []metricdata.HistogramDataPoint[T]) []metricdata.HistogramDataPoint[T] {
	out := make([]metricdata.HistogramDataPoint[T], len(hdp))
	for i, dp := range hdp {
		out[i] = metricdata.HistogramDataPoint[T]{
			Attributes:   dp.Attributes,
			Count:        dp.Count,
			Sum:          dp.Sum,
			Bounds:       dp.Bounds,
			BucketCounts: dp.BucketCounts,
			Min:          dp.Min,
			Max:          dp.Max,
		}
	}
	return out
}

// 改时间戳
func redactDataPointTimestamps[T int64 | float64](sdp []metricdata.DataPoint[T]) []metricdata.DataPoint[T] {
	out := make([]metricdata.DataPoint[T], len(sdp))
	for i, dp := range sdp {
		out[i] = metricdata.DataPoint[T]{
			Attributes: dp.Attributes,
			Value:      dp.Value,
		}
	}
	return out
}
