```go
// Exporter handles the delivery of metric data to external receivers. This is the final component in the metric push pipeline.
//Exporter负责将度量数据传递到外部接收器。这是度量推送管道中的最后一个组件。
type Exporter interface {
// Temporality returns the Temporality to use for an instrument kind.
//Temporality返回用于乐器种类的临时性。
Temporality(InstrumentKind) metricdata.Temporality
临时性（InstrumentKind）度量数据。临时性

// Aggregation returns the Aggregation to use for an instrument kind.
//Aggregation返回要用于工具类型的Aggregation。
Aggregation(InstrumentKind) aggregation.Aggregation
聚合（InstrumentKind）聚合。聚合

// Export serializes and transmits metric data to a receiver.
//导出串行化度量数据并将其传输到接收器。
// This is called synchronously, there is no concurrency safety requirement. Because of this, it is critical that all timeouts and cancellations of the passed context be honored.
//这被称为同步，没有并发安全要求。因此，遵守传递的上下文的所有超时和取消是至关重要的。
// All retry logic must be contained in this function. The SDK does not implement any retry logic. All errors returned by this function are considered unrecoverable and will be reported to a configured error Handler.
//此函数中必须包含所有重试逻辑。SDK没有实现任何重试逻辑。此函数返回的所有错误都被认为是不可恢复的，并将报告给配置的错误处理程序。
// The passed ResourceMetrics may be reused when the call completes. If an exporter needs to hold this data after it returns, it needs to make a copy.
//当调用完成时，可以重用传递的ResourceMetrics。如果出口商在返回后需要保存这些数据，则需要制作一份副本。
Export(context.Context, *metricdata.ResourceMetrics) error
导出（context.context，*metricdata.ResourceMetrics）错误

// ForceFlush flushes any metric data held by an exporter.
//ForceFlush刷新导出程序所持有的所有度量数据。
// The deadline or cancellation of the passed context must be honored. An appropriate error should be returned in these situations.
//必须遵守截止日期或取消已传递的上下文。在这些情况下，应该返回适当的错误。
ForceFlush(context.Context) error
ForceFlush（context.context）错误

// Shutdown flushes all metric data held by an exporter and releases any held computational resources.
//关闭将刷新导出程序所持有的所有度量数据，并释放所有持有的计算资源。
// The deadline or cancellation of the passed context must be honored. An appropriate error should be returned in these situations.
//必须遵守截止日期或取消已传递的上下文。在这些情况下，应该返回适当的错误。
// After Shutdown is called, calls to Export will perform no operation and instead will return an error indicating the shutdown state.
//调用Shutdown后，对Export的调用将不执行任何操作，而是返回一个指示关闭状态的错误。
Shutdown(context.Context) error
关闭（context.context）错误
}
```

```go
// MeterProvider provides access to named Meter instances, for instrumenting an application or package.
//MeterProvider提供对命名Meter实例的访问，用于检测应用程序或包。
// Warning: Methods may be added to this interface in minor releases. See package documentation on API implementation for information on how to set default behavior for unimplemented methods.
//警告：方法可能会在小版本中添加到此接口。有关如何为未实现的方法设置默认行为的信息，请参阅有关API实现的包文档。
type MeterProvider interface 
类型MeterProvider接口{
// MeterProvider is embedded in the OpenTelemetry metric API [MeterProvider].
//MeterProvider嵌入OpenTelemetry度量API[MeterProvider]中。
// Embed this interface in your implementation of the [MeterProvider] if you want users to experience a compilation error, signaling they need to update to your latest implementation, when the [MeterProvider] interface is extended (which is something that can happen without a major version bump of the API package).
//如果您希望用户在扩展[MeterProvider]接口时遇到编译错误，请将此接口嵌入[MeterProvider][计量提供者]的实现中，这表明他们需要更新到您的最新实现（这可能在没有API包的重大版本冲突的情况下发生）。
// [MeterProvider]: go.opentelemetry.io/otel/metric.MeterProvider
//[计量提供者]：go.opentelemetry.io/otel/metric.MeterProvider
embedded.MeterProvider
嵌入式.MeterProvider

// Meter returns a new Meter with the provided name and configuration.
//Meter返回一个具有提供的名称和配置的新Meter。
// A Meter should be scoped at most to a single package. The name needs to be unique so it does not collide with other names used by an application, nor other applications. To achieve this, the import path of the instrumentation package is recommended to be used as name.
//仪表的范围最多应为一个包。名称必须是唯一的，这样它就不会与应用程序或其他应用程序使用的其他名称冲突。为了实现这一点，建议使用检测包的导入路径作为名称。
// If the name is empty, then an implementation defined default name will be used instead.
//如果名称为空，则将使用实现定义的默认名称。
Meter(name string, opts ...MeterOption) Meter
仪表（名称字符串，选项…仪表选项）仪表
}

// Meter provides access to instrument instances for recording metrics.
//仪表提供对仪表实例的访问，以记录度量。
// Warning: Methods may be added to this interface in minor releases. See package documentation on API implementation for information on how to set default behavior for unimplemented methods.
//警告：方法可能会在小版本中添加到此接口。有关如何为未实现的方法设置默认行为的信息，请参阅有关API实现的包文档。
type Meter interface {
```

```go

// ResourceMetrics is a collection of ScopeMetrics and the associated Resource that created them.
//ResourceMetrics是ScopeMetrics和创建它们的相关资源的集合。
type ResourceMetrics struct {
	// Resource represents the entity that collected the metrics.
	//Resource表示收集度量的实体。
	Resource *resource.Resource
	// ScopeMetrics are the collection of metrics with unique Scopes.
	//ScopeMetrics是具有唯一Scopes的度量的集合。
	ScopeMetrics []ScopeMetrics
}

// ScopeMetrics is a collection of Metrics Produces by a Meter.
//ScopeMetrics是由Meter生成的度量的集合。
type ScopeMetrics struct {
	// Scope is the Scope that the Meter was created with.
	//范围是创建仪表时使用的范围。
	Scope instrumentation.Scope
	// Metrics are a list of aggregations created by the Meter.
	//度量是由仪表创建的聚合列表。
	Metrics []Metrics
}

// Metrics is a collection of one or more aggregated timeseries from an Instrument.
//度量是一个工具的一个或多个聚合时间序列的集合。
type Metrics struct {
	// Name is the name of the Instrument that created this data.
	//Name是创建该数据的仪器的名称。
	Name string
	// Description is the description of the Instrument, which can be used in documentation.
	//说明是对仪器的说明，可用于文件编制。
	Description string
	// Unit is the unit in which the Instrument reports.
	//单位是仪器报告的单位。
	Unit string
	// Data is the aggregated data from an Instrument.
	//数据是来自仪器的汇总数据。
	Data Aggregation
}
```
```go
const (
// undefinedTemporality represents an unset Temporality.
//undefinedTemporality表示一种未设置的时态。
undefinedTemporality Temporality = iota
// CumulativeTemporality defines a measurement interval that continues to expand forward in time from a starting point. New measurements are added to all previous measurements since a start time.
//累积时间定义了从一个起点开始在时间上继续向前扩展的测量间隔。新的测量值将添加到自开始时间以来的所有先前测量值中。
CumulativeTemporality
// DeltaTemporality defines a measurement interval that resets each cycle.
//DeltaTemporality定义了重置每个周期的测量间隔。
// Measurements from one cycle are recorded independently, measurements from other cycles do not affect them.
//一个周期的测量值是独立记录的，其他周期的测量不影响它们。
DeltaTemporality
)
```
```go
// HistogramDataPoint is a single histogram data point in a timeseries.
//直方图数据点是时间序列中的单个直方图数据点。
type HistogramDataPoint[N int64 | float64] struct {
// Attributes is the set of key value pairs that uniquely identify the timeseries.
//属性是唯一标识时间序列的一组键值对。
Attributes attribute.Set
// StartTime is when the timeseries was started.
//StartTime是时间序列开始的时间。
StartTime time.Time
// Time is the time when the timeseries was recorded.
//时间是记录时间序列的时间。
Time time.Time
// Count is the number of updates this histogram has been calculated with.
//Count是计算此直方图的更新次数。
Count uint64
// Bounds are the upper bounds of the buckets of the histogram. Because the last boundary is +infinity this one is implied.
//边界是直方图的桶的上限。因为最后一个边界是+无穷大，所以这个边界是隐含的。
Bounds []float64
// BucketCounts is the count of each of the buckets.
//BucketCounts是每个 桶 的计数。
BucketCounts []uint64
// Min is the minimum value recorded. (optional)
//Min是记录的最小值。（可选）
Min Extrema[N]
// Max is the maximum value recorded. (optional)
//Max是记录的最大值。（可选）
Max Extrema[N]
// Sum is the sum of the values recorded.
//总和是记录的值的总和。
Sum N
// Exemplars is the sampled Exemplars collected during the timeseries.
//示例是在时间序列期间收集的采样示例。
Exemplars []Exemplar[N] `json:",omitempty"`
}
```
