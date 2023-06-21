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

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/guance/guancetrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

const text string = `+
If send Info through local DK, token=""
You can set host and token by env
linux:   export GUANCE_HOST="localhost" ; export GUANCE_TOKEN="tkn_XXXXXXXX"
windows: set GUANCE_HOST="localhost" ; set GUANCE_TOKEN="tkn_XXXXXXXX"

You can set host and token in file ~/.guance/hostandtoken
localhost
tkn_XXXXXXXX

You can set host and token here`

var logger = log.New(os.Stderr, "guance-example", log.Ldate|log.Ltime|log.Llongfile)

// initTracer creates a new trace provider instance and registers it as global trace provider.
func initTracer(urlStr string) (func(context.Context) error, error) {
	// Create Zipkin Exporter and install it as a global tracer.
	//
	// For demoing purposes, always sample. In a production application, you should
	// configure the sampler to a trace.ParentBased(trace.TraceIDRatioBased) set at the desired
	// ratio.
	enc := guancetrace.ConvertorHolder{}

	exporter, err := guancetrace.New(
		urlStr,
		guancetrace.WithConvertor(enc),
		guancetrace.WithLogger(logger),
		// guancetrace.WithoutTimestamps(), // 必须删掉，不然没有时间戳
	)
	if err != nil {
		return nil, err
	}

	batcher := sdktrace.NewBatchSpanProcessor(exporter)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("guance_trace_example"),
		)),
	)
	otel.SetTracerProvider(tp)

	return tp.Shutdown, nil
}

func main() {
	host, token := getHostAndToken()

	// host = "http://127.0.0.1:9529" // TODO 测试代码回头删除
	// token = ""                     // TODO 测试代码回头删除

	if host == "" {
		panic("host is empty string.")
	}
	urlStr := fmt.Sprintf("%s/v1/write/%s?token=%s", host, "tracing", token)
	if token == "" {
		urlStr = fmt.Sprintf("%s/v1/write/%s", host, "tracing")
	}
	// fmt.Println("host and token:", urlStr) // 回头删除

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	stopCh := make(chan interface{})
	stopedCh := make(chan interface{})
	shutdown, err := initTracer(urlStr, stopCh, stopedCh)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	tr := otel.GetTracerProvider().Tracer("component-main")
	// === 第一组trace
	ctx, span := tr.Start(ctx, "foo", trace.WithSpanKind(trace.SpanKindServer))
	<-time.After(6 * time.Millisecond)
	bar(ctx)
	<-time.After(6 * time.Millisecond)
	span.End()
	// === 第二组trace

	// === END 第二组trace
	close(stopCh)
	<-stopedCh // 等待 feed 停止完成
	fmt.Println("END")
}

func bar(ctx context.Context) {
	tr := otel.GetTracerProvider().Tracer("component-bar")
	_, span := tr.Start(ctx, "bar")
	span.SetAttributes(attribute.KeyValue{Key: "http_host", Value: attribute.StringValue("http:/127.0.0.1:123")})
	span.SetAttributes(attribute.KeyValue{Key: "pid", Value: attribute.StringValue("2468")})
	span.SetAttributes(attribute.KeyValue{Key: "abcd", Value: attribute.StringValue("xyz")}) // not useful
	<-time.After(6 * time.Millisecond)
	span.End()
}

func getHostAndToken() (string, string) {
	var host, token string
	host, token = getFromEnv()
	if host != "" {
		fmt.Println("got host and token from ENV")
	}

	if host == "" {
		host, token = getFromFile()
		if host != "" {
			fmt.Println("got host and token from FILE")
		}
	}

	// host = "" // TODO 测试代码，要删除的

	if host == "" {
		fmt.Println(text)
		fmt.Println("your host?")
		fmt.Scanf("%s", &host)
		fmt.Println("your token?")
		fmt.Scanf("%s", &token)
	}

	return host, token
}

func getFromEnv() (string, string) {
	var host, token string
	host = os.Getenv("GUANCE_HOST")
	token = os.Getenv("GUANCE_TOKEN")
	return host, token
}

func getFromFile() (string, string) {
	var host, token string

	dir := "~/.guance/hostandtoken"
	expandedDir, err := homedir.Expand(dir)
	if err != nil {
		fmt.Println("error: ", err)
		return "", ""
	}
	fmt.Printf("Expand of %s is: %s\n", dir, expandedDir) // 回头删除

	content, err := ioutil.ReadFile(expandedDir)
	if err != nil {
		fmt.Println("read file error: ", err)
		return "", ""
	}
	s := strings.Split(string(content), "\n")
	if len(s) > 0 {
		host = s[0]
	}
	if len(s) > 1 {
		token = s[1]
	}

	return host, token
}
