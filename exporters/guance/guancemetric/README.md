# OpenTelemetry-Go Guancemetric Exporter



## Installation

```
go get -u go.opentelemetry.io/otel/exporters/guance/guancemetric
```

## Example

See [../../example/jaeger](../../example/guance/guancemetricexample).

## Configuration


### Environment Variables


## Contributing

This exporter uses a vendored copy of the Apache Thrift library (v0.14.1) at a custom import path.
When re-generating Thrift code in the future, please adapt import paths as necessary.

## References

# 开发笔记
## 按照otel项目要求 使用了go1.19
## metric项目
  expoter：在exporters/guance/guancemetric包
  feed(数据上传、重试6次)：在 exporters/guance/internal/feed包
  example：在example/guance/guancemetricexample/main.go
### 参考范例
  expoter，参考： exporters/stdout/stdoutmetric/exporter.go
  数据源 ，参考： exporters/stdout/stdoutmetric/example_test.go
  example，参考： exporters/stdout/stdoutmetric/example_test.go
  数据转换，参考： exporters/prometheus/exporter.go

# 开发指导文档
  otel文档根目录： github.com/open-telemetry/docs-cn
  快速入门： github.com/open-telemetry/docs-cn/blob/main/QUICKSTART.md
  标准规范： github.com/open-telemetry/docs-cn/tree/main/specification
  贡献者指南： github.com/open-telemetry/docs-cn/blob/main/community/CONTRIBUTING.md
  处理和导出数据： github.com/open-telemetry/docs-cn/blob/main/community/opentelemetryGoInstrumentation/exporting_data.md
  (英文)metric文档： github.com/open-telemetry/docs-cn/blob/main/specification/metrics/sdk.md