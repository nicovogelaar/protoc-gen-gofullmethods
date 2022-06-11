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
	methods := make([]string, 0)
	for _, s := range g.Services {
		for _, m := range s.Methods {
			methodName := normalizeFullname(protoreflect.FullName(s.Desc.Name())) + "_" + string(m.Desc.Name())
			g.P("\t", methodName, ` = "/`, s.Desc.FullName(), "/", m.Desc.Name(), `"`)
			methods = append(methods, methodName)
		}
	}
	g.P(")")
	g.P("var (")
	g.P("\tFullMethods = []string{")
	for _, m := range methods {
		g.P("\t\t", m, ",")
	}
	g.P("}")
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
