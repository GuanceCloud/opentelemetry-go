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
	"log"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"

	"go.opentelemetry.io/otel/exporters/guance/internal/feed"
	"go.opentelemetry.io/otel/sdk/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

var errShutdown = errors.New("exporter shutdown")
var zeroTime time.Time

// exporter is an OpenTelemetry metric exporter.
type Exporter struct {
	url    string
	client *http.Client
	logger logr.Logger
	// pointCh chan []*point.Point
	feedCh chan []feed.FeedInfo
	encVal atomic.Value // encoderHolder
	// wg         sync.WaitGroup // Shutdown wait deed return
	stopped atomic.Bool

	redactTimestamps bool

	encoderMu sync.Mutex
	stoppedMu sync.RWMutex

	stopedCh chan interface{}
}

var _ sdktrace.SpanExporter = &Exporter{}

var emptyLogger = logr.Logger{}

// config contains options for the exporter.
type config struct {
	client    *http.Client
	logger    logr.Logger
	convertor *ConvertorHolder
	// temporalitySelector metric.TemporalitySelector
	// aggregationSelector metric.AggregationSelector
	redactTimestamps bool
}

// Option defines a function that configures the exporter.
type Option interface {
	apply(config) config
}

type optionFunc func(config) config

func (fn optionFunc) apply(cfg config) config {
	return fn(cfg)
}

// WithLogger configures the exporter to use the passed logger.
// WithLogger and WithLogr will overwrite each other.
func WithLogger(logger *log.Logger) Option {
	return WithLogr(stdr.New(logger))
}

// WithLogr configures the exporter to use the passed logr.Logger.
// WithLogr and WithLogger will overwrite each other.
func WithLogr(logger logr.Logger) Option {
	return optionFunc(func(cfg config) config {
		cfg.logger = logger
		return cfg
	})
}

// WithClient configures the exporter to use the passed HTTP client.
func WithClient(client *http.Client) Option {
	return optionFunc(func(cfg config) config {
		cfg.client = client
		return cfg
	})
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

// WithoutTimestamps sets all timestamps to zero in the output stream.
func WithoutTimestamps() Option {
	return optionFunc(func(c config) config {
		c.redactTimestamps = true
		return c
	})
}

// Like : exporters/stdout/stdoutmetric/exporter.go
// New returns a configured metric exporter.
//
// If no options are passed, the default exporter returned will use a JSON
// encoder with tab indentations that output to STDOUT.
func New(collectorURL string, opts ...Option) (*Exporter, error) {
	u, err := url.Parse(collectorURL)
	if err != nil {
		return nil, fmt.Errorf("invalid collector URL %q: %v", collectorURL, err)
	}
	if u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("invalid collector URL %q: no scheme or host", collectorURL)
	}

	cfg := config{}
	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}

	exp := &Exporter{
		url:              collectorURL,
		client:           cfg.client,
		logger:           cfg.logger,
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
	e.feedCh <- feedInfos

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
	close(e.stopCh)
	return ctx.Err()
}
