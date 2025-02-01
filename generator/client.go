package generator

import (
	"strings"

	"github.com/mtvy/protoc-plugin-nats/stringutils"
	descriptor "google.golang.org/protobuf/types/descriptorpb"
)

func (g *generator) generateClient(service *descriptor.ServiceDescriptorProto) {
	serviceName := stringutils.CamelCase(service.GetName())
	clientTypeName := stringutils.PascalCase(service.GetName()) + "Client"

	g.P(`	   	type ClientConnector interface {`)
	g.P(`	   		RequestWithContext(ctx context.Context, subject string, in, out protoreflect.ProtoMessage) error`)
	g.P(`	   		Close()`)
	g.P(`	   	}`)
	g.P(``)
	g.P(`	   	func New`, serviceName, `Client(client ClientConnector, name string) `, serviceName, ` {`)
	g.P(`	   		return &`, clientTypeName, `{`)
	g.P(`	   			client: client,`)
	g.P(`	   			name:   name,`)
	g.P(`	   		}`)
	g.P(`	   	}`)
	g.P(``)
	g.P(`	   	func (u *`, clientTypeName, `) getSubject(endpoint string) string {`)
	g.P(`	   		return fmt.Sprintf("%s.%s", u.name, endpoint)`)
	g.P(`	   	}`)
	g.P(``)
	g.P(`	   var _ `, serviceName, ` = (*`, clientTypeName, `)(nil)`)
	g.P(``)
	g.P(`	   	type `, clientTypeName, ` struct {`)
	g.P(`	   		client ClientConnector`)
	g.P(`	   		name string`)
	g.P(`	   	}`)
	g.P(``)

	for _, method := range service.Method {
		inputType := g.goTypeName(method.GetInputType())
		outputType := g.goTypeName(method.GetOutputType())
		endpointName := strings.ToLower(method.GetName())

		g.P(`	   	func (impl *`, clientTypeName, `) `, method.GetName(), `(ctx context.Context, in *`, inputType, `) (*`, outputType, `, error) {`)
		g.P(`	   		out := new(`, outputType, `)`)
		g.P(`	   		err := impl.client.RequestWithContext(ctx, impl.getSubject("`, endpointName, `"), in, out)`)
		g.P(`	   		if err != nil {`)
		g.P(`	   			return nil, err`)
		g.P(`	   		}`)
		g.P(`	   		return out, nil`)
		g.P(`	   	}`)
		g.P()
	}
	g.P()
}
