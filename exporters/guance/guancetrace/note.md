```go
type SpanExporter interface {
	// DO NOT CHANGE: any modification will not be backwards compatible and must never be done outside of a new major release.
	// ExportSpans exports a batch of spans.
	// This function is called synchronously, so there is no concurrency safety requirement. However, due to the synchronous calling pattern, it is critical that all timeouts and cancellations contained in the passed context must be honored.
	// Any retry logic must be contained in this function. The SDK that calls this function will not implement any retry logic. All errors returned by this function are considered unrecoverable and will be reported to a configured error Handler.
    // 不要更改：任何修改都是不向后兼容的，并且永远不能在新的主要版本之外进行。
    // ExportSpans导出一批跨度。
    // 此函数是同步调用的，因此没有并发安全要求。但是，由于同步调用模式，必须遵守传递的上下文中包含的所有超时和取消，这一点至关重要。
    // 此函数中必须包含任何重试逻辑。调用此函数的SDK不会实现任何重试逻辑。此函数返回的所有错误都被认为是不可恢复的，并将报告给配置的错误处理程序。
	ExportSpans(ctx context.Context, spans []ReadOnlySpan) error
	// DO NOT CHANGE: any modification will not be backwards compatible and must never be done outside of a new major release.
	// Shutdown notifies the exporter of a pending halt to operations. The exporter is expected to perform any cleanup or synchronization it requires while honoring all timeouts and cancellations contained in the passed context.
    //不要更改：任何修改都是不向后兼容的，并且永远不能在新的主要版本之外进行。
    //关闭会通知导出程序操作暂停。导出器应执行所需的任何清理或同步，同时遵守传递上下文中包含的所有超时和取消。
	Shutdown(ctx context.Context) error
	// DO NOT CHANGE: any modification will not be backwards compatible and must never be done outside of a new major release.
    //不要更改：任何修改都是不向后兼容的，并且永远不能在新的主要版本之外进行。
}
```


```go

| Tag     | container_host   |      //  host name of container                                                                              |
| Tag     | endpoint         |      //  end point of resource                                                                               |
| Tag     | env              |      //  environment arguments                                                                               |
| Tag     | http_host        |      //  HTTP host                                                                                           |
| Tag     | http_method      |      //  HTTP method                                                                                         |
| Tag     | http_route       |      //  HTTP route                                                                                          |
| Tag     | http_status_code |      //  HTTP status code                                                                                    |
| Tag     | http_url         |      //  HTTP URL                                                                                            |
| Tag     | operation        |      //  operation of resource                                                                               |
| Tag     | pid              |      //  process id                                                                                          |
| Tag     | project          |      //  project name                                                                                        |
| Tag     | service          |      //  service name                                                                                        |
| Tag     | source_type      |      //  source types [app, framework, cache, message_queue, custom, db, web]                                |
2 Tag     | status           |      //  span status [ok, info, warning, error, critical]                                                    |
2 Tag     | span_type        |      //  span types [entry, local, exit, unknown]                                                            |
2 Field   | duration         |      // 微秒 | span duration                                                                                       |
| Field   | message          |      //  raw data content                                                                                    |
2 Field   | parent_id        |      //  parent ID of span                                                                                   |
| Field   | priority         |      //  priority rules (PRIORITY_USER_REJECT, PRIORITY_AUTO_REJECT, PRIORITY_AUTO_KEEP, PRIORITY_USER_KEEP) |
2 Field   | resource         |      //  resource of service                                                                                 |
| Field   | sample_rate      |      //  global sampling ratio (0.1 means roughly 10 percent will send to data center)                       |
2 Field   | span_id          |      //  span ID                                                                                             |
2 Field   | start            |      // 微秒 | span start timestamp                                                                                |
2 Field   | trace_id         |      //  trace ID                                                                                            |

	// Types that are assignable to Val:
	//	*Field_I
	//	*Field_U
	//	*Field_F
	//	*Field_B
	//	*Field_D
	//	*Field_A


type isField_Val interface {
		isField_Val()
}
type Field_I struct {
		I int64 `protobuf:"varint,2,opt,name=i,proto3,oneof"` // signed int
}
type Field_U struct {
		U uint64 `protobuf:"varint,3,opt,name=u,proto3,oneof"` // unsigned int
}
type Field_F struct {
		F float64 `protobuf:"fixed64,4,opt,name=f,proto3,oneof"` // float64
}
type Field_B struct {
		B bool `protobuf:"varint,5,opt,name=b,proto3,oneof"` // bool
}
type Field_D struct {
		D []byte `protobuf:"bytes,6,opt,name=d,proto3,oneof"` // bytes, for string or binary data
}
type Field_A struct {
		// XXX: not used
	A *anypb.Any `protobuf:"bytes,7,opt,name=a,proto3,oneof"` // any data
}

```


```go
Span Type 为当前 span 在 trace 中的相对位置，其取值说明如下：

- entry：当前 api 为入口即链路进入进入服务后的第一个调用
- local: 当前 api 为入口后出口前的 api
- exit: 当前 api 为链路在服务上最后一个调用
- unknown: 当前 api 的相对位置状态不明确

Priority Rules 为客户端采样优先级规则

- `PRIORITY_USER_REJECT = -1` 用户选择拒绝上报
- `PRIORITY_AUTO_REJECT = 0` 客户端采样器选择拒绝上报
- `PRIORITY_AUTO_KEEP = 1` 客户端采样器选择上报
- `PRIORITY_USER_KEEP = 2` 用户选择上报
```


### Datakit Tracing Span 数据结构 {#span-struct}

``` golang
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
```

Datakit Span 是 Datakit 内部使用的数据结构。第三方 Tracing Agent 数据结构会转换成 Datakit Span 结构后发送到数据中心。

> 以下简称 DKSpan

| Field Name | Data Type                  | Unit | Description                                   | Correspond To              |
| ---------- | ------------------------   | ---- | -------------------------------------------   | ------------------------   |
| TraceID    | `string`                   |      | Trace ID                                      | `dkproto.fields.trace_id`  |
| ParentID   | `string`                   |      | Parent Span ID                                | `dkproto.fields.parent_id` |
| SpanID     | `string`                   |      | Span ID                                       | `dkproto.fields.span_id`   |
| Service    | `string`                   |      | Service Name                                  | `dkproto.tags.service`     |
| Resource   | `string`                   |      | Resource Name(.e.g `/get/data/from/some/api`) | `dkproto.fields.resource`  |
| Operation  | `string`                   |      | 生产此条 Span 的方法名                        | `dkproto.tags.operation`   |
| Source     | `string`                   |      | Span 接入源(.e.g `ddtrace`)                   | `dkproto.name`             |
| SpanType   | `string`                   |      | Span Type(.e.g `entry`)                       | `dkproto.tags.span_type`   |
| SourceType | `string`                   |      | Span Source Type(.e.g `web`)                  | `dkproto.tags.type`        |
| Tags       | `map[string, string]`      |      | Span Tags                                     | `dkproto.tags`             |
| Metrics    | `map[string, interface{}]` |      | Span Metrics(计算用)                          | `dkproto.fields`           |
| Start      | `int64`                    | 纳秒 | Span 起始时间                                 | `dkproto.fields.start`     |
| Duration   | `int64`                    | 纳秒 | 耗时                                          | `dkproto.fields.duration`  |
| Status     | `string`                   |      | Span 状态字段                                 | `dkproto.tags.status`      |
| Content    | `string`                   |      | Span 原始数据                                 | `dkproto.fields.message`   |



OpenTelemetry 中的 `resource_spans` 和 DKSpan 的对应关系如下：

| Field Name           | Data Type           | Unit | Description    | Correspond To                  |
| ---                  | ---                 | ---  | ---            | ---                            |
| trace_id             | `[16]byte`          |      | Trace ID       | `dkspan.TraceID`               |
| span_id              | `[8]byte`           |      | Span ID        | `dkspan.SpanID`                |
| parent_span_id       | `[8]byte`           |      | Parent Span ID | `dkspan.ParentID`              |
| name                 | `string`            |      | Span Name      | `dkspan.Operation`             |
| kind                 | `string`            |      | Span Type      | `dkspan.SpanType`              |
| start_time_unix_nano | `int64`             | 纳秒 | Span 起始时间  | `dkspan.Start`                   |
| end_time_unix_nano   | `int64`             | 纳秒 | Span 终止时间  | `dkspan.Duration = end -start`   |
| status               | `string`            |      | Span Status    | `dkspan.Status`                |
| name                 | `string`            |      | resource Name  | `dkspan.Resource`              |
| resource.attributes  | `map[string]string` |      | resource 标签  | XXX                             |
| span.attributes      | `map[string]string` |      | Span 标签      | `dkspan.tags`                   |

`dkspan.tags.service, dkspan.tags.project, dkspan.tags.env, dkspan.tags.version, dkspan.tags.container_host, dkspan.tags.http_method, dkspan.tags.http_status_code`


OpenTelemetry 有些独有字段， 但 DKSpan 没有字段与之对应，所以就放在了标签中，只有这些值非 0 时才会显示，如：

| Field                         | Date Type | Uint | Description             | Correspond                             |
| :---                          | :---      | :--- | :---                    | :---                                   |
2 span.dropped_attributes_count | `int`     |      | Span 被删除的标签数量   | `dkspan.tags.dropped_attributes_count` |
2 span.dropped_events_count     | `int`     |      | Span 被删除的事件数量   | `dkspan.tags.dropped_events_count`     |
2 span.dropped_links_count      | `int`     |      | Span 被删除的连接数量   | `dkspan.tags.dropped_links_count`      |
2 span.events_count             | `int`     |      | Span 关联事件数量       | `dkspan.tags.events_count`             |
2 span.links_count              | `int`     |      | Span 所关联的 span 数量 | `dkspan.tags.links_count`              |

---

实际的otel-trace结构

```go
// SpanStub is a stand-in for a Span.
type SpanStub struct {
	Name                   string
	SpanContext            trace.SpanContext
	Parent                 trace.SpanContext
	SpanKind               trace.SpanKind
	StartTime              time.Time
	EndTime                time.Time
	Attributes             []attribute.KeyValue
	Events                 []tracesdk.Event
	Links                  []tracesdk.Link
	Status                 tracesdk.Status
	DroppedAttributes      int
	DroppedEvents          int
	DroppedLinks           int
	ChildSpanCount         int
	Resource               *resource.Resource
	InstrumentationLibrary instrumentation.Library
}



```