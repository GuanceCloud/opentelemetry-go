// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package guancetrace // import "go.opentelemetry.io/otel/exporters/guance/guancetrace "

const (
	// datakit tracing customer tags.
	UNKNOWN_SERVICE = "unknown_service"
	CONTAINER_HOST  = "container_host"
	ENV             = "env"
	PROJECT         = "project"
	VERSION         = "version"

	// span status.
	STATUS_OK       = "ok"
	STATUS_INFO     = "info"
	STATUS_DEBUG    = "debug"
	STATUS_WARN     = "warning"
	STATUS_ERR      = "error"
	STATUS_CRITICAL = "critical"

	// span position in trace.
	SPAN_TYPE_ENTRY   = "entry"
	SPAN_TYPE_LOCAL   = "local"
	SPAN_TYPE_EXIT    = "exit"
	SPAN_TYPE_UNKNOWN = "unknown"

	// span source type.
	SPAN_SOURCE_APP       = "app"
	SPAN_SOURCE_FRAMEWORK = "framework"
	SPAN_SOURCE_CACHE     = "cache"
	SPAN_SOURCE_MSGQUE    = "message_queue"
	SPAN_SOURCE_CUSTOMER  = "custom"
	SPAN_SOURCE_DB        = "db"
	SPAN_SOURCE_WEB       = "web"

	// line protocol tags.
	TAG_CONTAINER_HOST   = "container_host"
	TAG_ENDPOINT         = "endpoint"
	TAG_ENV              = "env"
	TAG_HTTP_HOST        = "http_host"
	TAG_HTTP_METHOD      = "http_method"
	TAG_HTTP_ROUTE       = "http_route"
	TAG_HTTP_STATUS_CODE = "http_status_code"
	TAG_HTTP_URL         = "http_url"
	TAG_OPERATION        = "operation"
	TAG_PID              = "pid"
	TAG_PROJECT          = "project"
	TAG_SERVICE          = "service"
	TAG_SOURCE_TYPE      = "source_type"
	TAG_SPAN_STATUS      = "status"
	TAG_SPAN_TYPE        = "span_type"
	TAG_VERSION          = "version"

	// line protocol fields.
	FIELD_DURATION    = "duration"
	FIELD_MESSAGE     = "message"
	FIELD_PARENTID    = "parent_id"
	FIELD_PRIORITY    = "priority"
	FIELD_RESOURCE    = "resource"
	FIELD_SAMPLE_RATE = "sample_rate"
	FIELD_SPANID      = "span_id"
	FIELD_START       = "start"
	FIELD_TRACEID     = "trace_id"
	FIELD_ERR_MESSAGE = "error_message"
	FIELD_ERR_STACK   = "error_stack"
	FIELD_ERR_TYPE    = "error_type"
	FIELD_CALL_TREE   = "calling_tree"

	TRACE_128_BIT_ID = "trace_128_bit_id"
)

// Attributes binding to resource.
const (
	otelResourceServiceKey        = "service.name"
	otelResourceServiceVersionKey = "service.version"
	otelResourceContainerNameKey  = "container.name"
	otelResourceProcessIDKey      = "process.pid"
)

//nolint:deadcode,unused,varcheck
const (
	// HTTP.
	otelHTTPSchemeKey     = "http.scheme"
	otelHTTPMethodKey     = "http.method"
	otelHTTPStatusCodeKey = "http.status_code"
	// 从 otel.span 对象解析到 datakit.span 中的时候，有些字段无法没有对应，不应当主动丢弃，暂时放进tags中
	// see : vendor/go.opentelemetry.io/proto/otlp/trace/v1/trace.pb.go:383.
	DroppedAttributesCount = "dropped_attributes_count"
	DroppedEventsCount     = "dropped_events_count"
	DroppedLinksCount      = "dropped_links_count"
	Events                 = "events_count"
	Links                  = "links_count"
	// database.
	otelDBSystemKey = "db.system"
	// message queue.
	otelMessagingSystemKey = "messaging.system"
	// rpc system.
	otelRPCSystemKey = "rpc.system"
)

// Attributes binding to event
const (
	ExceptionEventName     = "exception"
	ExceptionTypeKey       = "exception.type"
	ExceptionMessageKey    = "exception.message"
	ExceptionStacktraceKey = "exception.stacktrace"
)

var otelErrKeyToDkErrKey = map[string]string{
	ExceptionTypeKey:       FIELD_ERR_TYPE,
	ExceptionMessageKey:    FIELD_ERR_MESSAGE,
	ExceptionStacktraceKey: FIELD_ERR_STACK,
}
