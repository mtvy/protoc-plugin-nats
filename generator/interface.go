package generator

import (
	"fmt"
	"strings"

	"github.com/mtvy/protoc-plugin-nats/stringutils"
	"github.com/mtvy/protoc-plugin-nats/typemap"
	descriptor "google.golang.org/protobuf/types/descriptorpb"
)

func (g *generator) generateInterfaceClient(file *descriptor.FileDescriptorProto, service *descriptor.ServiceDescriptorProto) {
	serviceName := stringutils.CamelCase(service.GetName())
	comments, err := g.reg.ServiceComments(file, service)
	if err == nil {
		g.printComments(comments)
		g.P(`// `)
	}
	g.P(`// Service will use the "`, file.GetPackage(), `.`, strings.ToLower(serviceName), `" service group`)

	g.P(`type `, serviceName, ` interface {`)
	for _, method := range service.Method {
		comments, err := g.reg.MethodComments(file, service, method)
		if err == nil {
			g.printComments(comments)
			g.P(`// `)
		}
		g.P(`// Method subject will be "`, file.GetPackage(), `.`, strings.ToLower(serviceName), `.`, strings.ToLower(method.GetName()), `"`)

		g.P(g.generateSignatureClient(method))
		g.P()
	}
	g.P(`}`)
}

func (g *generator) generateSignatureClient(method *descriptor.MethodDescriptorProto) string {
	methName := stringutils.CamelCase(method.GetName())
	inputType := g.goTypeName(method.GetInputType())
	outputType := g.goTypeName(method.GetOutputType())
	return fmt.Sprintf(`	%s(%s.Context, *%s) (*%s, error)`, methName, g.pkgs["context"], inputType, outputType)
}

func (g *generator) generateInterfaceServer(file *descriptor.FileDescriptorProto, service *descriptor.ServiceDescriptorProto) {
	serviceName := stringutils.CamelCase(service.GetName())
	comments, err := g.reg.ServiceComments(file, service)
	if err == nil {
		g.printComments(comments)
		g.P(`// `)
	}
	g.P(`// Service will use the "`, file.GetPackage(), `.`, strings.ToLower(serviceName), `" service group`)

	g.P(`type `, serviceName, `Server interface {`)
	for _, method := range service.Method {
		comments, err := g.reg.MethodComments(file, service, method)
		if err == nil {
			g.printComments(comments)
			g.P(`// `)
		}
		g.P(`// Method subject will be "`, file.GetPackage(), `.`, strings.ToLower(serviceName), `.`, strings.ToLower(method.GetName()), `"`)

		g.P(g.generateSignatureServer(method))
		g.P()
	}
	g.P(`}`)
}

func (g *generator) generateSignatureServer(method *descriptor.MethodDescriptorProto) string {
	methName := stringutils.CamelCase(method.GetName())
	inputType := g.goTypeName(method.GetInputType())
	outputType := g.goTypeName(method.GetOutputType())
	return fmt.Sprintf(`	%s(%s.Context, *%s, *%s) error`, methName, g.pkgs["context"], inputType, outputType)
}

func (g *generator) printComments(comments typemap.DefinitionComments) bool {
	text := strings.TrimSuffix(comments.Leading, "\n")
	if len(strings.TrimSpace(text)) == 0 {
		return false
	}
	split := strings.Split(text, "\n")
	for _, line := range split {
		g.P("// ", strings.TrimPrefix(line, " "))
	}
	return len(split) > 0
}
