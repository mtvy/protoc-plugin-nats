package generator

import (
	"strings"

	"github.com/mtvy/protoc-plugin-nats/stringutils"
	descriptor "google.golang.org/protobuf/types/descriptorpb"
)

func (g *generator) generateRegister(service *descriptor.ServiceDescriptorProto) {
	g.generateRegisterBuilder(service)

	g.generateHandlers(service)

	g.generateMethodsArray(service)
}

func (g *generator) generateHandlers(service *descriptor.ServiceDescriptorProto) {
	serviceName := stringutils.CamelCase(service.GetName())
	for _, method := range service.Method {
		inputType := g.goTypeName(method.GetInputType())
		outputType := g.goTypeName(method.GetOutputType())

		g.P(`func _`, serviceName, `_`, method.GetName(), `_Handler(ctx `, g.pkgs["context"], `.Context, req micro.Request, impl interface{}) (protoreflect.ProtoMessage, error) {`)
		g.P(`	in := new(`, inputType, `)`)
		g.P(`	if err := proto.Unmarshal(req.Data(), in); err != nil {`)
		g.P(`		return nil, err`)
		g.P(`	}`)
		g.P(`	out := new(`, outputType, `)`)
		g.P(`	err := impl.(`, serviceName, `Server).`, method.GetName(), `(ctx, in, out)`)
		g.P(`	return out, err`)
		g.P(`}`)
		g.P()
	}
}

func (g *generator) generateRegisterBuilder(service *descriptor.ServiceDescriptorProto) {
	serviceName := stringutils.CamelCase(service.GetName())
	g.P(`type Registrar interface {`)
	g.P(`	// TODO: поменять group interface{} на type`)
	g.P(`	RegisterService(ctx context.Context, group interface{}, impl interface{}) error`)
	g.P(`	Close()`)
	g.P(`}`)
	g.P()
	g.P(`// New`, serviceName, `Server builds a new micro.Service that will be registered with the instance provided`)
	g.P(`// Each RPC on the service will be mapped to a new endpoint within the micro service`)
	g.P(`func Register`, serviceName, `Server(ctx `, g.pkgs["context"], `.Context, service Registrar, impl `, serviceName, `Server) error {`)
	g.P(`	return service.RegisterService(ctx, `, serviceName, `_GroupDesc, impl)`)
	g.P(`}`)
	g.P()
}

func (g *generator) generateMethodsArray(service *descriptor.ServiceDescriptorProto) {
	serviceName := stringutils.CamelCase(service.GetName())

	g.P(`// TODO: сделать лучше`)
	g.P(`type GroupDescer interface {`)
	g.P(`	GetMethods() []interface{}`)
	g.P(`}`)
	g.P(``)
	g.P(`type GroupDesc struct {`)
	g.P(`	Methods []interface{}`)
	g.P(`}`)
	g.P(``)
	g.P(`func (g *GroupDesc) GetMethods() []interface{} {`)
	g.P(`	return g.Methods`)
	g.P(`}`)
	g.P(``)
	g.P(`// TODO: сделать лучше`)
	g.P(`type MethodDescer interface {`)
	g.P(`	GetName() string`)
	g.P(`	GetHandler() func(ctx context.Context, req micro.Request, impl interface{}) (protoreflect.ProtoMessage, error)`)
	g.P(`}`)
	g.P(``)
	g.P(`type MethodDesc struct {`)
	g.P(`	Name string`)
	g.P(`	Handler func(ctx context.Context, req micro.Request, impl interface{}) (protoreflect.ProtoMessage, error)`)
	g.P(`}`)
	g.P(``)
	g.P(`func (m *MethodDesc) GetName() string {`)
	g.P(`	return m.Name`)
	g.P(`}`)
	g.P(``)
	g.P(`func (m *MethodDesc) GetHandler() func(ctx context.Context, req micro.Request, impl interface{}) (protoreflect.ProtoMessage, error) {`)
	g.P(`	return m.Handler`)
	g.P(`}`)
	g.P(``)
	g.P(`var `, serviceName, `_GroupDesc = &GroupDesc{`)
	g.P(`	Methods: []interface{}{`)

	for _, method := range service.Method {
		g.P(`		&MethodDesc{`)
		g.P(`			Name:    "`, strings.ToLower(method.GetName()), `",`)
		g.P(`			Handler: _`, serviceName, `_`, method.GetName(), `_Handler,`)
		g.P(`		},`)
	}
	g.P(`	},`)
	g.P(`}`)
	g.P()
}
