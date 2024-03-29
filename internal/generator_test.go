package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	fds   *descriptorpb.FileDescriptorSet = &descriptorpb.FileDescriptorSet{}
	files *protoregistry.Files            = &protoregistry.Files{}
	p     *protogen.Plugin
)

func TestMain(m *testing.M) {
	os.RemoveAll("test/test.pb.descriptor")
	cmd := exec.Command("protoc", "-o", "test/test.pb.descriptor", "--include_imports", "-I", "test", "test.proto")
	err := cmd.Start()
	if err != nil {
		panic(fmt.Errorf("failed to generate protobuf descriptor: %s, %q", err.Error(), strings.Join(cmd.Args, " ")))
	}
	err = cmd.Wait()
	if err != nil {
		panic(fmt.Errorf("failed to generate protobuf descriptor: %s, %q", err.Error(), strings.Join(cmd.Args, " ")))
	}
	gengocmd := exec.Command("protoc", "--go_out=:test", "--go_opt=paths=source_relative", "-I", "test", "test.proto")
	err = gengocmd.Start()
	if err != nil {
		panic(fmt.Errorf("failed to generate protobuf descriptor: %s, %q", err.Error(), strings.Join(gengocmd.Args, " ")))
	}
	err = gengocmd.Wait()
	if err != nil {
		panic(fmt.Errorf("failed to generate protobuf descriptor: %s, %q", err.Error(), strings.Join(gengocmd.Args, " ")))
	}
	b, err := ioutil.ReadFile("./test/test.pb.descriptor")
	if err != nil {
		panic(fmt.Errorf("failed to parse proto descriptor: %w", err))
	}
	err = proto.Unmarshal(b, fds)
	if err != nil {
		panic(fmt.Errorf("failed to unmarshall proto descriptor: %w", err))
	}
	files, err = protodesc.NewFiles(fds)
	if err != nil {
		panic(fmt.Errorf("failed to read FileDescriptorSet: %w", err))
	}
	p, err = protogen.Options{}.New(&pluginpb.CodeGeneratorRequest{
		ProtoFile: fds.File,
	})
	if err != nil {
		panic(fmt.Errorf("failed to generate plugin: %w", err))
	}

	os.Exit(m.Run())
}

func TestPlugin(t *testing.T) {
	importPath := p.FilesByPath["test.proto"].GoImportPath
	gen := p.NewGeneratedFile("test.pb.fullmethods.go", importPath)
	g := &generator{gen, p.FilesByPath["test.proto"]}
	var (
		complierMajor int32 = 1
		complierMinor int32 = 26
		complierPatch int32 = 5
	)
	os.Setenv("BUILD_VERSION", "v1.0.0")
	req := pluginpb.CodeGeneratorRequest{
		CompilerVersion: &pluginpb.Version{
			Major: &complierMajor,
			Minor: &complierMinor,
			Patch: &complierPatch,
		},
	}
	g.Generate("test.proto", &req)
	raw, err := g.Content()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	wantStr := `// Code generated by protoc-gen-gofullmethods. DO NOT EDIT.
// versions:
// - protoc-gen-gofullmethods (unknown)
// - protoc                   v1.26.5
// source: test.proto

package test

const (
	TestService_abc = "/example.TestService/abc"
	TestService_Abc = "/example.TestService/Abc"
)

var (
	FullMethods = []string{
		TestService_abc,
		TestService_Abc,
	}
)
`
	if res := cmp.Diff(string(raw), wantStr); res != "" {
		t.Errorf("(+want/-got) %s", res)
	}
}

func TestGetVersion(t *testing.T) {
	tc := []struct {
		name        string
		buildInfo   BuildInfoReaderFunc
		wantVersion string
	}{
		{
			"get version from build info",
			func() (*debug.BuildInfo, bool) {
				return &debug.BuildInfo{
					Main: debug.Module{
						Version: "v1.0.1",
					},
				}, false
			},
			"v1.0.1",
		},
		{
			"get (unknown) when BuildInfo main version is empty",
			func() (*debug.BuildInfo, bool) {
				return &debug.BuildInfo{
					Main: debug.Module{
						Version: "",
					},
				}, true
			},
			"(unknown)",
		},
		{
			"get (unknown) when failed to get BuildInfo",
			func() (*debug.BuildInfo, bool) {
				return nil, false
			},
			"(unknown()",
		},
	}
	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			res := GetVersion(c.buildInfo)
			if c.wantVersion != res {
				fmt.Printf("expect build version: %s, got: %s", c.wantVersion, res)
			}
		})
	}
}
