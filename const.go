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

const Alias = "alias"
const ActorAlias = "actor_alias"
const ActorAskReentrant = "actor_ask_reentrant"
const GrpcStyle = "grpc_style"
const CsProxyDefault = "cs_proxy_default"
const ReturnPacket = "return_packet"
const AsyncCall = "async_call"
