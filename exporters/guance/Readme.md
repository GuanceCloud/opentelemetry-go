# otel-exporter 支持

## 项目目标
- 按照 otel 项目规范，把 metric 和 trace 直接打到 dataway。
- 按照 otel 项目规范，把 metric 和 trace 打到 Datakit，转发 dataway。

## 一阶段完成目标
- exporters/guance/guancemetric 包实现 metric 转为行协议，打给 feed 包。
- example/guance/guancemetricexample 项目生成 mock metric 打给 guancemetric 包。
- exporters/guance/guancetrace 包把 otel-trace 转为 dktrice 再转为行协议，打给 feed 包。
- example/guance/guancetraceexample 项目生成 mock trice 打给 guancetrace 包。
- exporters/guance/internal/feed 包通过 init() 启动唯一协程，定时通过通道接收行协议，汇总、压缩、上传。（这样做当 exporter 来大量碎片数据的时候，可以聚合上传，加大数据承载能力）
- 分支名 step3-add-gzon

### 参考代码
- exporters/guance/guancemetric 参考 exporters/stdout/stdoutmetric 较多。
- exporters/guance/guancetrace 参考 exporters/zipkin 较多
- example/guance/guancemetricexample 参考 exporters/stdout/stdoutmetric/example_test.go 较多。
- example/guance/guancetraceexample 参考 example/zipkin 较多。

## 二阶设定目标（只在 trace 上开了头）
- feed 包大改造，复刻 Datakit 项目的 internal/io/dataway/endpoint.go。
- 每一个 exporter 实例持有一个 endpoint 实例。
- 每一个 exporter 实例拉起来一个 feed 协程，透过 endpoint 实例上传数据。（一阶段，feed 协程是永远不关闭的）
- 分支名 step4-add-avast-retry-go