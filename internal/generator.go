package internal

import (
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/pluginpb"
)

type generator struct {
	*protogen.GeneratedFile
	*protogen.File
}

func normalizeFullname(fn protoreflect.FullName) string {
	return strings.ReplaceAll(string(fn), ".", "_")
}

func (g *generator) Generate() {
	g.P("package ", g.GoPackageName)
	g.P("const (")
	for _, s := range g.Services {
		for _, m := range s.Methods {
			g.P("\tMethod_", normalizeFullname(s.Desc.FullName()), "__", m.Desc.Name(), ` = "/`, s.Desc.FullName(), "/", m.Desc.Name(), `"`)
		}
	}
	g.P(")")
}

func Run(opt protogen.Options) {
	opt.Run(func(p *protogen.Plugin) error {
		p.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range p.Files {
			if !f.Generate {
				continue
			}
			filename := f.GeneratedFilenamePrefix + "_methods.pb.go"
			file := p.NewGeneratedFile(filename, f.GoImportPath)
			g := generator{file, f}
			g.Generate()
		}
		return nil
	})
}
