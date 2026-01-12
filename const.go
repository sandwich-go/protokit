package protokit

const AnnotationPrefix = "annotation@"
const AnnotationSpaceAnnotationType = "annotation_type"
const AnnotationSpaceRawdata = "rawdata"
const AnnotationSpaceAB = "ab"
const AnnotationSpaceFilter = "filter"

const AnnotationKeyType = "type"
const AnnotationProfoFileTypeRawdata = "rawdata"
const AnnotationProfoFileTypeRawdataConst = "rawdata_const"

const AnnotationService = "service"
const AnnotationGlobal = "global"

type ServiceTag = string

const (
	ServiceTagALL   ServiceTag = "all"
	ServiceTagRPC   ServiceTag = "rpc"
	ServiceTagActor ServiceTag = "actor"
	ServiceTagERPC  ServiceTag = "erpc"
	ServiceTagJob   ServiceTag = "job"
	ServiceTagQuit  ServiceTag = "quit"
)

const QueryPathSnakeCase = "query_path_snake_case"
const QueryPath = "query_path"
const ServiceUriAutoAlias = "service_uri_auto_alias"
const Tell = "tell"
const LangOff = "lang_off"
const ActorSystemName = "actor_system_name"
const Labels = "labels"
const ShortId = "short_id"
const HandleTimeout = "handle_timeout"

const Alias = "alias"
const ActorAlias = "actor_alias"
const ActorAskReentrant = "actor_ask_reentrant"
const GrpcStyle = "grpc_style"
const CsProxyDefault = "cs_proxy_default"
const ReturnPacket = "return_packet"
const AsyncCall = "async_call"
const CsActorIdSource = "cs_actor_id_source"
const CsAutoResend = "cs_auto_resend"
const CsRpcBlocking = "cs_rpc_blocking"
const CsWeakNetworkThreshold = "cs_weak_network_threshold"
const CsDisconnectionThreshold = "cs_disconnection_threshold"
const Custom = "custom"
