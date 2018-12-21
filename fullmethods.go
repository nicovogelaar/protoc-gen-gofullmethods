package main

import (
	"bytes"
	"go/format"
	"html/template"
	"path"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	googlegen "github.com/golang/protobuf/protoc-gen-go/generator"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/pkg/errors"
)

const (
	tmpl = `
// Code generated by protoc-gen-go. DO NOT EDIT.
// source: {{ .FileName }}

package {{ .GoPackage }}

const (
{{- range $service := .Services }}
{{- range $method := $service.Methods }}
	{{ $service.Name }}_{{ $method }} = "/{{ $.Package }}.{{ $service.Name }}/{{ $method }}"
{{- end}}
{{- end}}
)

var (
	FullMethods = []string{{ "{" }}
{{- range $service := .Services }}
{{- range $method := $service.Methods }}
		{{ $service.Name }}_{{ $method }},
{{- end}}
{{- end}}
	{{ "}" }}
)
`
)

type service struct {
	Name    string
	Methods []string
}

type data struct {
	FileName  string
	GoPackage string
	Package   string
	Services  []service
}

type fullMethodsGenerator struct {
	*generator
}

func newFullMethodsGenerator() *fullMethodsGenerator {
	return &fullMethodsGenerator{generator: newGenerator()}
}

func (g *fullMethodsGenerator) generate() error {
	return g.generator.generate(g)
}

func (g *fullMethodsGenerator) generateFile(protoFile *descriptor.FileDescriptorProto) (*plugin.CodeGeneratorResponse_File, error) {
	if protoFile.Name == nil {
		return nil, errors.New("missing filename")
	}
	if protoFile.GetOptions().GetGoPackage() == "" {
		return nil, errors.New("missing go_package")
	}

	dat := data{
		FileName:  *protoFile.Name,
		GoPackage: protoFile.GetOptions().GetGoPackage(),
		Package:   protoFile.GetPackage(),
		Services:  make([]service, len(protoFile.Service)),
	}

	for _, svc := range protoFile.Service {
		methods := make([]string, len(svc.Method))
		for i, method := range svc.Method {
			methods[i] = *method.Name
		}
		dat.Services = append(dat.Services, service{
			Name:    googlegen.CamelCase(svc.GetName()),
			Methods: methods,
		})
	}

	buf := bytes.NewBuffer(nil)

	err := template.Must(template.New("").Parse(tmpl)).Execute(buf, dat)
	if err != nil {
		return nil, err
	}

	g.P(buf.String())

	formatted, err := format.Source(g.Bytes())
	if err != nil {
		return nil, err
	}

	return &plugin.CodeGeneratorResponse_File{
		Name:    proto.String(protoFileBaseName(*protoFile.Name) + ".fullmethods.pb.go"),
		Content: proto.String(string(formatted)),
	}, nil
}

func protoFileBaseName(name string) string {
	if ext := path.Ext(name); ext == ".proto" {
		name = name[:len(name)-len(ext)]
	}
	return name
}
