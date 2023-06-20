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
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/GuanceCloud/cliutils/point"

	"go.opentelemetry.io/otel/exporters/guance/internal/feed"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

var errShutdown = errors.New("exporter shutdown")
var zeroTime time.Time

var _ trace.SpanExporter = &Exporter{}

// exporter is an OpenTelemetry metric exporter.
type Exporter struct {
	url     string
	pointCh chan []*point.Point
	encVal  atomic.Value // encoderHolder
	// wg         sync.WaitGroup // Shutdown wait deed return
	stopped atomic.Bool

	redactTimestamps bool

	encoderMu sync.Mutex
	stoppedMu sync.RWMutex
}

// Like : exporters/stdout/stdoutmetric/exporter.go
// New returns a configured metric exporter.
//
// If no options are passed, the default exporter returned will use a JSON
// encoder with tab indentations that output to STDOUT.
func New(url string, options ...Option) (trace.SpanExporter, error) {
	cfg := newConfig(options...)

	exp := &Exporter{
		url: url,
		// pointCh:             make(chan []*point.Point, 1),
		// temporalitySelector: cfg.temporalitySelector,
		// aggregationSelector: cfg.aggregationSelector,
		redactTimestamps: cfg.redactTimestamps,
	}
	exp.encVal.Store(*cfg.convertor)

	return exp, nil
}

// ExportSpans exports a batch of spans.
// This function is called synchronously, so there is no concurrency safety requirement. However, due to the synchronous calling pattern, it is critical that all timeouts and cancellations contained in the passed context must be honored.
// Any retry logic must be contained in this function. The SDK that calls this function will not implement any retry logic. All errors returned by this function are considered unrecoverable and will be reported to a configured error Handler.
// ExportSpans导出一批跨度。
// 此函数是同步调用的，因此没有并发安全要求。但是，由于同步调用模式，必须遵守传递的上下文中包含的所有超时和取消，这一点至关重要。
// 此函数中必须包含任何重试逻辑。调用此函数的SDK不会实现任何重试逻辑。此函数返回的所有错误都被认为是不可恢复的，并将报告给配置的错误处理程序。
func (e *Exporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	e.stoppedMu.RLock()
	stopped := e.stopped
	e.stoppedMu.RUnlock()
	if stopped.Load() {
		return nil
	}

	if len(spans) == 0 {
		return nil
	}

	stubs := tracetest.SpanStubsFromReadOnlySpans(spans)

	convertorHolder := e.encVal.Load().(ConvertorHolder)

	e.encoderMu.Lock()
	defer e.encoderMu.Unlock()
	for i := range stubs {
		if e.redactTimestamps {
			stubs[i].StartTime = zeroTime
			stubs[i].EndTime = zeroTime
			for j := range stubs[i].Events {
				ev := &stubs[i].Events[j]
				ev.Time = zeroTime
			}
		}
	}

	// Encode span stubs, one by one
	points, err := convertorHolder.Convert(stubs) //  TODO 需要改
	if err != nil {
		return err
	}
	feedInfos := make([]feed.FeedInfo, 0)

	for _, pt := range points {
		feedInfos = append(feedInfos, feed.FeedInfo{LineProto: pt.LineProto() + "\n", URL: e.url})
		// s := pt.LineProto()
		// fmt.Println("# pt.LineProto() == ", s)
		// break

	}

	fmt.Println("Export发送给chan")
	feed.FeedCh <- feedInfos

	return nil
}

// Shutdown notifies the exporter of a pending halt to operations. The exporter is expected to perform any cleanup or synchronization it requires while honoring all timeouts and cancellations contained in the passed context.
//关闭会通知导出程序操作暂停。导出器应执行所需的任何清理或同步，同时遵守传递上下文中包含的所有超时和取消。
// Shutdown is called to stop the exporter, it performs no action.
func (e *Exporter) Shutdown(ctx context.Context) error {
	if e.stopped.Load() {
		return errShutdown
	}
	e.stopped.Swap(true) // Set exporter shutdown

	return ctx.Err()
}
