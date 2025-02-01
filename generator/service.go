package generator

import (
	"strings"

	"github.com/mtvy/protoc-plugin-nats/stringutils"

	descriptor "google.golang.org/protobuf/types/descriptorpb"
)

func (g *generator) generateService(file *descriptor.FileDescriptorProto, service *descriptor.ServiceDescriptorProto) {
	serviceName := stringutils.CamelCase(service.GetName())

	g.sectionComment(serviceName + ` Interface`)
	g.generateInterfaceClient(file, service)

	g.sectionComment(serviceName + ` NATS Client`)
	g.generateClient(service)

	g.sectionComment(serviceName + `Server Interface`)
	g.generateInterfaceServer(file, service)

	g.sectionComment(serviceName + ` Micro Service Handlers`)
	g.generateRegister(service)
}

// Big header comments to makes it easier to visually parse a generated file.
func (g *generator) sectionComment(sectionTitle string) {
	g.P()
	g.P(`// `, strings.Repeat("=", len(sectionTitle)))
	g.P(`// `, sectionTitle)
	g.P(`// `, strings.Repeat("=", len(sectionTitle)))
	g.P()
}
