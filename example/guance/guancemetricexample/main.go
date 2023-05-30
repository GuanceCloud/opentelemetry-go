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
	"os"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/guance/guancemetric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

const text string = `
You can set host and token by env
linux:   export GUANCE_HOST="localhost" ; export GUANCE_TOKEN="tkn_XXXXXXXX"
windows: set GUANCE_HOST="localhost" ; set GUANCE_TOKEN="tkn_XXXXXXXX"

You can set host and token in file ~/.guance/hostandtoken
localhost
tkn_XXXXXXXX

You can set host and token here`

func main() {
	host, token := getHostAndToken()
	if host == "" {
		panic("host is empty string.")
	}
	fmt.Println("host and token:", host, token) // 回头删除

	enc := guancemetric.ConvertorHolder{}

	exporter, err := guancemetric.New(
		host,
		token,
		guancemetric.WithConvertor(enc),
		// guancemetric.WithoutTimestamps(), // 必须删掉，不然没有时间戳
	)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	//Example one. This is where the sdk would be used to create a Meter and from that instruments that would make measurements of your code. To simulate that behavior, call export directly with mocked data.
	//样例一。在这里，sdk将用于创建一个Meter，并从中对代码进行测量。要模拟这种行为，请使用模拟数据直接调用export。
	fmt.Println("样例一： ============")

	for i := 0; i < 1; i++ {
		duration := (i % 3) + 1
		f := float64(i+40) + float64(duration+1)/100
		mockData := creatMockData(time.Now(), f)
		_ = exporter.Export(ctx, &mockData)
		fmt.Println("完成：", i)
		time.Sleep(time.Second * time.Duration(duration))
	}

	fmt.Println("关闭sdk ============")
	// Ensure the periodic reader is cleaned up by shutting down the sdk.
	//通过关闭sdk来确保定期读卡器已清理干净。
	fmt.Println(exporter.Shutdown(ctx))
	fmt.Println(exporter.Shutdown(ctx)) // 故意重复关闭
	mockData := creatMockData(time.Now(), 12.34)
	fmt.Println(exporter.Export(ctx, &mockData)) // 故意关闭后尝试发送数据
	time.Sleep(time.Second * 2)
	fmt.Println("退出 main ============")
}

func getHostAndToken() (string, string) {
	var host, token string
	host, token = getFromEnv()
	fmt.Println("ENV host and token:", host, token) // 回头删除
	if host == "" {
		host, token = getFromFile()
	}
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

func creatMockData(now time.Time, f float64) metricdata.ResourceMetrics {
	res := resource.NewSchemaless(
		semconv.ServiceName("stdoutmetric-example"),
	)
	return metricdata.ResourceMetrics{
		Resource: res,
		ScopeMetrics: []metricdata.ScopeMetrics{
			{
				Scope: instrumentation.Scope{Name: "example", Version: "v0.0.1"},
				Metrics: []metricdata.Metrics{
					{
						Name:        "requests",
						Description: "Number of requests received",
						Unit:        "1",
						Data: metricdata.Sum[int64]{
							IsMonotonic: true,
							Temporality: metricdata.DeltaTemporality,
							DataPoints: []metricdata.DataPoint[int64]{
								{
									Attributes: attribute.NewSet(attribute.String("server", "central")),
									StartTime:  now,
									Time:       now.Add(1 * time.Second),
									Value:      5,
								},
							},
						},
					},
					{
						Name:        "system.cpu.time",
						Description: "Accumulated CPU time spent",
						Unit:        "s",
						Data: metricdata.Sum[float64]{
							IsMonotonic: true,
							Temporality: metricdata.CumulativeTemporality,
							DataPoints: []metricdata.DataPoint[float64]{
								{
									Attributes: attribute.NewSet(attribute.String("state", "user")),
									StartTime:  now,
									Time:       now.Add(1 * time.Second),
									Value:      0.5,
								},
							},
						},
					},
					{
						Name:        "latency",
						Description: "Time spend processing received requests",
						Unit:        "ms",
						Data: metricdata.Histogram[float64]{
							Temporality: metricdata.DeltaTemporality,
							DataPoints: []metricdata.HistogramDataPoint[float64]{
								{
									Attributes:   attribute.NewSet(attribute.String("server", "central")),
									StartTime:    now,
									Time:         now.Add(1 * time.Second),
									Count:        10,
									Bounds:       []float64{1, 5, 10},
									BucketCounts: []uint64{1, 3, 6, 0},
									Sum:          57,
								},
							},
						},
					},
					{
						Name:        "system.memory.usage",
						Description: "Memory usage",
						Unit:        "By",
						Data: metricdata.Gauge[int64]{
							DataPoints: []metricdata.DataPoint[int64]{
								{
									Attributes: attribute.NewSet(attribute.String("state", "used")),
									Time:       now.Add(1 * time.Second),
									Value:      100,
								},
							},
						},
					},
					{
						Name:        "temperature",
						Description: "CPU global temperature",
						Unit:        "cel(1 K)",
						Data: metricdata.Gauge[float64]{
							DataPoints: []metricdata.DataPoint[float64]{
								{
									Attributes: attribute.NewSet(attribute.String("server", "central")),
									Time:       now.Add(1 * time.Second),
									Value:      f, // 这个值有变化
								},
							},
						},
					},
				},
			},
		},
	}
}