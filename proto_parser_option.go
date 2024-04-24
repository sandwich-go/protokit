package protokit

import (
	"fmt"

	"github.com/jhump/protoreflect/desc"
	protokit2 "github.com/sandwich-go/protokit/option/gen/golang/protokit"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func getRpcServiceOption(sd *desc.ServiceDescriptor) *protokit2.RpcServiceOptions {
	msgO := sd.GetServiceOptions()
	opts, ok := proto.GetExtension(msgO, protokit2.E_RpcService).(*protokit2.RpcServiceOptions)
	if ok {
		return opts
	}
	return nil
}

func isJobService(sd *desc.ServiceDescriptor) bool {
	msgO := sd.GetServiceOptions()
	opts, ok := proto.GetExtension(msgO, protokit2.E_JobService).(bool)
	if ok {
		return opts
	}
	return false
}

func getBackOfficeServiceOption(sd *desc.ServiceDescriptor) *protokit2.BackOfficeServiceOptions {
	msgO := sd.GetServiceOptions()
	opts, ok := proto.GetExtension(msgO, protokit2.E_BackOfficeService).(*protokit2.BackOfficeServiceOptions)
	if ok {
		return opts
	}
	return nil
}

func getRpcMethodOption(md *descriptorpb.MethodDescriptorProto) *protokit2.RpcMethodOptions {
	msgO := md.GetOptions()
	opts, ok := proto.GetExtension(msgO, protokit2.E_RpcMethod).(*protokit2.RpcMethodOptions)
	if ok {
		return opts
	}
	return nil
}

func getJobMethodOption(md *descriptorpb.MethodDescriptorProto) *protokit2.JobMethodOptions {
	msgO := md.GetOptions()
	opts, ok := proto.GetExtension(msgO, protokit2.E_JobMethod).(*protokit2.JobMethodOptions)
	if ok {
		return opts
	}
	return nil
}

func getBackOfficeMethodOption(md *descriptorpb.MethodDescriptorProto) *protokit2.BackOfficeMethodOptions {
	msgO := md.GetOptions()
	opts, ok := proto.GetExtension(msgO, protokit2.E_BackOfficeMethod).(*protokit2.BackOfficeMethodOptions)
	if ok {
		return opts
	}
	return nil
}

type serviceAnnotation interface {
	Bool(key string, defaultVal ...bool) (bool, error)
	String(key string, defaultVal ...string) string
	Contains(key string) bool
}

type serviceOptionAnnotation struct {
	*protokit2.RpcServiceOptions
}

func (so *serviceOptionAnnotation) Bool(key string, defaultVal ...bool) (bool, error) {
	switch key {
	case QueryPathSnakeCase:
		return so.QueryPathSnakeCase == protokit2.QueryPathType_SNAKE_CASE, nil
	case ServiceUriAutoAlias:
		return false, nil
	case ServiceTagActor:
		return so.Actor, nil
	case ServiceTagERPC:
		return so.Erpc, nil
	case ServiceTagRPC:
		return so.Rpc, nil
	case Tell:
		return so.AskTell == protokit2.MethodAskType_TELL, nil
	case ActorAskReentrant:
		return so.ActorAskReentrant, nil
	default:
		panic(fmt.Sprintf("RpcServiceOptions get bool unknown key: %s", key))
	}
	return false, nil
}

func (so *serviceOptionAnnotation) String(key string, defaultVal ...string) string {
	switch key {
	case QueryPath:
		if so.QueryPath == "" && len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return so.QueryPath
	case LangOff:
		return so.LangOff
	default:
		panic(fmt.Sprintf("RpcServiceOptions get string unknown key: %s", key))
	}
	return ""
}

func (so *serviceOptionAnnotation) Contains(key string) bool { return false }

type methodeAnnotation interface {
	Bool(key string, defaultVal ...bool) (bool, error)
	String(key string, defaultVal ...string) string
	Contains(key string) bool
}

type methodOptionAnnotation struct {
	*protokit2.RpcMethodOptions
}

func (so *methodOptionAnnotation) Bool(key string, defaultVal ...bool) (bool, error) {
	var dft bool
	if len(defaultVal) > 0 {
		dft = defaultVal[0]
	}
	switch key {
	case ServiceTagActor:
		return so.Actor || dft, nil
	case ServiceTagERPC:
		return so.Erpc || dft, nil
	case ServiceTagRPC:
		return so.Rpc || dft, nil
	case Tell:
		return (so.AskTell == protokit2.MethodAskType_TELL) || dft, nil
	case ActorAskReentrant:
		return so.ActorAskReentrant, nil
	case ServiceTagQuit:
		return so.Quit, nil
	default:
		panic(fmt.Sprintf("RpcMethodOptions get bool unknown key: %s", key))
	}
	return false, nil
}

func (so *methodOptionAnnotation) String(key string, defaultVal ...string) string {
	switch key {
	case Alias:
		return so.Alias
	case ActorAlias:
		return so.ActorAlias
	case LangOff:
		return so.LangOff
	default:
		panic(fmt.Sprintf("RpcMethodOptions get string unknown key: %s", key))
	}
	return ""
}

func (so *methodOptionAnnotation) Contains(key string) bool {
	switch key {
	case Alias:
		return so.Alias != ""
	case ActorAlias:
		return so.ActorAlias != ""
	}
	return false
}
