module go.opentelemetry.io/otel/exporters/guance/guancetrace

go 1.19

require (
	github.com/GuanceCloud/cliutils v0.1.1
	go.opentelemetry.io/otel v1.15.0-rc.2
// go.opentelemetry.io/otel/exporters/guance/internal/feed v1.15.0-rc.2
)

require github.com/influxdata/line-protocol/v2 v2.2.1 // indirect

require (
	github.com/aliyun/aliyun-oss-go-sdk v2.1.2+incompatible // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/influxdata/influxdb1-client v0.0.0-20200827194710-b269163b24ab // indirect
	github.com/kr/pretty v0.2.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.42.0 // indirect
	github.com/rs/xid v1.2.1 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	go.opentelemetry.io/otel/sdk v1.15.0-rc.2
	go.opentelemetry.io/otel/trace v1.15.0-rc.2 // indirect
	golang.org/x/sys v0.7.0 // indirect
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)

replace go.opentelemetry.io/otel/metric => ../../../metric

// replace go.opentelemetry.io/otel/exporters/guance/internal/feed  => ../internal/feed

replace go.opentelemetry.io/otel => ../../..

replace go.opentelemetry.io/otel/sdk/metric => ../../../sdk/metric

replace go.opentelemetry.io/otel/trace => ../../../trace

replace go.opentelemetry.io/otel/sdk => ../../../sdk