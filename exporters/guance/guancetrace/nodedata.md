
```go
func main() {
	message := "{\"trace_id\":\"RTXnv/SvgjdqYNQiWOT0kZSkluKPyIKR/lKR52i5\",\"span_id\":\"7peXV49LwnY=\",\"name\":\"startRootSpan\",\"kind\":1,\"start_time_unix_nano\":1686707957120476562,\"end_time_unix_nano\":1686707957357489730,\"attributes\":[{\"key\":\"resource.name\",\"value\":{\"Value\":{\"StringValue\":\"/\"}}},{\"key\":\"span.type\",\"value\":{\"Value\":{\"StringValue\":\"web\"}}}],\"status\":{}}"
	_ = message
	metricName := "otel"

	opts := point.DefaultMetricOptions()
	tags := map[string]string{}
	fields := map[string]interface{}{}

	tags["host"] = "HP-serer"
	tags["user"] = "root"

	fields["float64Data"] = float64(1.234)
	fields["int64Data"] = int64(5678)
	// fields["stringData"] = message

	pt := point.NewPointV2([]byte(metricName), append(point.NewTags(tags), point.NewKVs(fields)...), opts...)
	_ = pt

	kvs := point.NewKVs(map[string]any{"f1": 123})

	kvs = kvs.Add([]byte(`t2`), []byte(message), false, true)

	for _, v := range kvs {
		pt.AddKV(v)
	}

	fmt.Println("pt.LineProto() == ", pt.LineProto())
	return
}
// ===============================================================
                spattrs := extractAtrributes(span.Attributes)

				dkspan := &itrace.DatakitSpan{
					TraceID:   hex.EncodeToString(span.GetTraceId()),
					ParentID:  byteToString(span.GetParentSpanId()),
					SpanID:    byteToString(span.GetSpanId()),
					Resource:  span.Name,
					Operation: span.Name,
					Source:    inputName,
					Tags:      make(map[string]string),
					Metrics:   make(map[string]interface{}),
					Start:     int64(span.StartTimeUnixNano),
					Duration:  int64(span.EndTimeUnixNano - span.StartTimeUnixNano),
					Status:    getDKSpanStatus(span.GetStatus()),
				}
				dkspan.SpanType = itrace.FindSpanTypeStrSpanID(dkspan.SpanID, dkspan.ParentID, spanIDs, parentIDs)

				attrs := newAttributes(resattrs).merge(scpattrs...).merge(spattrs...)
				if kv, i := attrs.find(otelResourceServiceKey); i != -1 {
					dkspan.Service = kv.Value.GetStringValue()
				}
				if kv, i := attrs.find(otelResourceServiceVersionKey); i != -1 {
					dkspan.Tags[itrace.TAG_VERSION] = kv.Value.GetStringValue()
				}
				if kv, i := attrs.find(otelResourceProcessIDKey); i != -1 {
					dkspan.Tags[itrace.TAG_PID] = kv.Value.GetStringValue()
				}
				if kv, i := attrs.find(otelResourceContainerNameKey); i != -1 {
					dkspan.Tags[itrace.TAG_CONTAINER_HOST] = kv.Value.GetStringValue()
				}
				if kv, i := attrs.find(otelHTTPMethodKey); i != -1 {
					dkspan.Tags[itrace.TAG_HTTP_METHOD] = kv.Value.GetStringValue()
					attrs.remove(otelHTTPMethodKey)
				}
				if kv, i := attrs.find(otelHTTPStatusCodeKey); i != -1 {
					dkspan.Tags[itrace.TAG_HTTP_STATUS_CODE] = kv.Value.GetStringValue()
					attrs.remove(otelHTTPStatusCodeKey)
				}

				for i := range span.Events {
					if span.Events[i].Name == ExceptionEventName {
						for o, d := range otelErrKeyToDkErrKey {
							if attr, ok := getAttribute(o, span.Events[i].Attributes); ok {
								dkspan.Metrics[d] = attr.Value.GetStringValue()
							}
						}
						break
					}
				}

				attrtags, attrfields := attrs.splite()
				dkspan.Tags = itrace.MergeTags(tags, dkspan.Tags, attrtags)
				dkspan.Metrics = itrace.MergeFields(dkspan.Metrics, attrfields)

				dkspan.SourceType = getSourceType(dkspan.Tags)

				if buf, err := json.Marshal(span); err != nil {
					log.Warn(err.Error())
				} else {
					dkspan.Content = string(buf)
				}

				dktrace = append(dktrace, dkspan)

```

```go
category doFeed ==  /v1/write/tracing
opentelemetry,host=zhangub-OMEN-by-HP-Laptop-15-dc1xxx,operation=startRootSpan,resource_name=/,service=dktrace-otel-agent,service_name=dktrace-otel-agent,source_type=custom,span_type=web,status=ok duration=237013i,message="{\"trace_id\":\"RTXnv/SvgjdqYNQiWOT0kZSkluKPyIKR/lKR52i5\",\"span_id\":\"7peXV49LwnY=\",\"name\":\"startRootSpan\",\"kind\":1,\"start_time_unix_nano\":1686707957120476562,\"end_time_unix_nano\":1686707957357489730,\"attributes\":[{\"key\":\"resource.name\",\"value\":{\"Value\":{\"StringValue\":\"/\"}}},{\"key\":\"span.type\",\"value\":{\"Value\":{\"StringValue\":\"web\"}}}],\"status\":{}}",parent_id="0",priority=1i,resource="startRootSpan",span_id="ee9797578f4bc276",start=1686707957120476i,trace_id="4535e7bff4af82376a60d42258e4f49194a496e28fc88291fe5291e768b9" 1686707957120476562
opentelemetry,host=zhangub-OMEN-by-HP-Laptop-15-dc1xxx,operation=set,resource_name=redis.set,service=dktrace-otel-agent,service_name=dktrace-otel-agent,source_type=custom,span_type=cache,status=ok,ttl=3690 duration=151216i,message="{\"trace_id\":\"RTXnv/SvgjdqYNQiWOT0kZSkluKPyIKR/lKR52i5\",\"span_id\":\"vcAa9ePVYAQ=\",\"parent_span_id\":\"7peXV49LwnY=\",\"name\":\"set\",\"kind\":1,\"start_time_unix_nano\":1686707957318953357,\"end_time_unix_nano\":1686707957470170161,\"attributes\":[{\"key\":\"resource.name\",\"value\":{\"Value\":{\"StringValue\":\"redis.set\"}}},{\"key\":\"span.type\",\"value\":{\"Value\":{\"StringValue\":\"cache\"}}},{\"key\":\"ttl\",\"value\":{\"Value\":{\"StringValue\":\"3690\"}}}],\"status\":{}}",parent_id="ee9797578f4bc276",resource="set",span_id="bdc01af5e3d56004",start=1686707957318953i,trace_id="4535e7bff4af82376a60d42258e4f49194a496e28fc88291fe5291e768b9" 1686707957318953357
opentelemetry,host=zhangub-OMEN-by-HP-Laptop-15-dc1xxx,operation=user.auth,resource_name=/authenticate,service=dktrace-otel-agent,service_name=dktrace-otel-agent,source_type=custom,span_type=web,status=ok duration=200973i,message="{\"trace_id\":\"RTXnv/SvgjdqYNQiWOT0kZSkluKPyIKR/lKR52i5\",\"span_id\":\"2YCHg2OkqRo=\",\"parent_span_id\":\"2Uq0efEqtFI=\",\"name\":\"user.auth\",\"kind\":1,\"start_time_unix_nano\":1686707957339493149,\"end_time_unix_nano\":1686707957540466788,\"attributes\":[{\"key\":\"resource.name\",\"value\":{\"Value\":{\"StringValue\":\"/authenticate\"}}},{\"key\":\"span.type\",\"value\":{\"Value\":{\"StringValue\":\"web\"}}}],\"status\":{}}",parent_id="d94ab479f12ab452",resource="user.auth",span_id="d980878363a4a91a",start=1686707957339493i,trace_id="4535e7bff4af82376a60d42258e4f49194a496e28fc88291fe5291e768b9" 1686707957339493149
opentelemetry,db_shard=xxx-xx-xxxx-xx,host=zhangub-OMEN-by-HP-Laptop-15-dc1xxx,operation=mysql.query,resource_name=select\ name\,\ age\,\ ts\ from\ 'user'\ where\ id\=5678765678,service=dktrace-otel-agent,service_name=dktrace-otel-agent,source_type=custom,status=ok duration=301093i,message="{\"trace_id\":\"RTXnv/SvgjdqYNQiWOT0kZSkluKPyIKR/lKR52i5\",\"span_id\":\"5zwmoh80WF4=\",\"parent_span_id\":\"7peXV49LwnY=\",\"name\":\"mysql.query\",\"kind\":1,\"start_time_unix_nano\":1686707957319135749,\"end_time_unix_nano\":1686707957620229742,\"attributes\":[{\"key\":\"resource.name\",\"value\":{\"Value\":{\"StringValue\":\"select name, age, ts from 'user' where id=5678765678\"}}},{\"key\":\"span.type\",\"value\":{\"Value\":{\"StringValue\":\"\"}}},{\"key\":\"db.shard\",\"value\":{\"Value\":{\"StringValue\":\"xxx-xx-xxxx-xx\"}}}],\"status\":{}}",parent_id="ee9797578f4bc276",resource="mysql.query",span_id="e73c26a21f34585e",start=1686707957319135i,trace_id="4535e7bff4af82376a60d42258e4f49194a496e28fc88291fe5291e768b9" 1686707957319135749
opentelemetry,host=zhangub-OMEN-by-HP-Laptop-15-dc1xxx,operation=user.getUserName,resource_name=/get/user/name,service=dktrace-otel-agent,service_name=dktrace-otel-agent,source_type=custom,span_type=web,status=ok duration=301371i,message="{\"trace_id\":\"RTXnv/SvgjdqYNQiWOT0kZSkluKPyIKR/lKR52i5\",\"span_id\":\"2Uq0efEqtFI=\",\"parent_span_id\":\"7peXV49LwnY=\",\"name\":\"user.getUserName\",\"kind\":1,\"start_time_unix_nano\":1686707957318905540,\"end_time_unix_nano\":1686707957620277293,\"attributes\":[{\"key\":\"resource.name\",\"value\":{\"Value\":{\"StringValue\":\"/get/user/name\"}}},{\"key\":\"span.type\",\"value\":{\"Value\":{\"StringValue\":\"web\"}}}],\"status\":{}}",parent_id="ee9797578f4bc276",resource="user.getUserName",span_id="d94ab479f12ab452",start=1686707957318905i,trace_id="4535e7bff4af82376a60d42258e4f49194a496e28fc88291fe5291e768b9" 1686707957318905540
opentelemetry,host=zhangub-OMEN-by-HP-Laptop-15-dc1xxx,operation=user.school,resource_name=/get/user/school,service=dktrace-otel-agent,service_name=dktrace-otel-agent,source_type=custom,span_type=web,status=ok duration=201254i,message="{\"trace_id\":\"RTXnv/SvgjdqYNQiWOT0kZSkluKPyIKR/lKR52i5\",\"span_id\":\"sa4qyutf5aI=\",\"parent_span_id\":\"q/MsOf224PY=\",\"name\":\"user.school\",\"kind\":1,\"start_time_unix_nano\":1686707957569634671,\"end_time_unix_nano\":1686707957770888917,\"attributes\":[{\"key\":\"resource.name\",\"value\":{\"Value\":{\"StringValue\":\"/get/user/school\"}}},{\"key\":\"span.type\",\"value\":{\"Value\":{\"StringValue\":\"web\"}}}],\"status\":{}}",parent_id="abf32c39fdb6e0f6",resource="user.school",span_id="b1ae2acaeb5fe5a2",start=1686707957569634i,trace_id="4535e7bff4af82376a60d42258e4f49194a496e28fc88291fe5291e768b9" 1686707957569634671
opentelemetry,host=zhangub-OMEN-by-HP-Laptop-15-dc1xxx,operation=user.id,resource_name=/get/user/id,service=dktrace-otel-agent,service_name=dktrace-otel-agent,source_type=custom,span_type=web,status=ok duration=302065i,message="{\"trace_id\":\"RTXnv/SvgjdqYNQiWOT0kZSkluKPyIKR/lKR52i5\",\"span_id\":\"q/MsOf224PY=\",\"parent_span_id\":\"2YCHg2OkqRo=\",\"name\":\"user.id\",\"kind\":1,\"start_time_unix_nano\":1686707957485185754,\"end_time_unix_nano\":1686707957787251658,\"attributes\":[{\"key\":\"resource.name\",\"value\":{\"Value\":{\"StringValue\":\"/get/user/id\"}}},{\"key\":\"span.type\",\"value\":{\"Value\":{\"StringValue\":\"web\"}}}],\"status\":{}}",parent_id="d980878363a4a91a",resource="user.id",span_id="abf32c39fdb6e0f6",start=1686707957485185i,trace_id="4535e7bff4af82376a60d42258e4f49194a496e28fc88291fe5291e768b9" 1686707957485185754
opentelemetry,host=zhangub-OMEN-by-HP-Laptop-15-dc1xxx,operation=user.class,resource_name=/get/user/class,service=dktrace-otel-agent,service_name=dktrace-otel-agent,source_type=custom,span_type=web,status=ok duration=201951i,message="{\"trace_id\":\"RTXnv/SvgjdqYNQiWOT0kZSkluKPyIKR/lKR52i5\",\"span_id\":\"kEXiCAbcoeA=\",\"parent_span_id\":\"sa4qyutf5aI=\",\"name\":\"user.class\",\"kind\":1,\"start_time_unix_nano\":1686707957614905682,\"end_time_unix_nano\":1686707957816857063,\"attributes\":[{\"key\":\"resource.name\",\"value\":{\"Value\":{\"StringValue\":\"/get/user/class\"}}},{\"key\":\"span.type\",\"value\":{\"Value\":{\"StringValue\":\"web\"}}}],\"status\":{}}",parent_id="b1ae2acaeb5fe5a2",resource="user.class",span_id="9045e20806dca1e0",start=1686707957614905i,trace_id="4535e7bff4af82376a60d42258e4f49194a496e28fc88291fe5291e768b9" 1686707957614905682
opentelemetry,group-access-token=kjskdafhcFertyuiknbvj,host=zhangub-OMEN-by-HP-Laptop-15-dc1xxx,operation=user.number,resource_name=/get/user/number,service=dktrace-otel-agent,service_name=dktrace-otel-agent,source_type=custom,span_type=web,status=ok duration=200827i,message="{\"trace_id\":\"RTXnv/SvgjdqYNQiWOT0kZSkluKPyIKR/lKR52i5\",\"span_id\":\"SZFc93wuGgw=\",\"parent_span_id\":\"kEXiCAbcoeA=\",\"name\":\"user.number\",\"kind\":1,\"start_time_unix_nano\":1686707957620291333,\"end_time_unix_nano\":1686707957821118663,\"attributes\":[{\"key\":\"resource.name\",\"value\":{\"Value\":{\"StringValue\":\"/get/user/number\"}}},{\"key\":\"span.type\",\"value\":{\"Value\":{\"StringValue\":\"web\"}}},{\"key\":\"group-access-token\",\"value\":{\"Value\":{\"StringValue\":\"kjskdafhcFertyuiknbvj\"}}}],\"status\":{}}",parent_id="9045e20806dca1e0",resource="user.number",span_id="49915cf77c2e1a0c",start=1686707957620291i,trace_id="4535e7bff4af82376a60d42258e4f49194a496e28fc88291fe5291e768b9" 1686707957620291333
opentelemetry,host=zhangub-OMEN-by-HP-Laptop-15-dc1xxx,operation=user.score,resource_name=/get/user/score,service=dktrace-otel-agent,service_name=dktrace-otel-agent,source_type=custom,span_type=web,status=error duration=200867i,error_message="access deny",error_stack="goroutine 43 [running]:

```