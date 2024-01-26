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
)

const QueryPathSnakeCase = "query_path_snake_case"
const QueryPath = "query_path"
const ServiceUriAutoAlias = "service_uri_auto_alias"
const Tell = "tell"
const LangOff = "lang_off"
