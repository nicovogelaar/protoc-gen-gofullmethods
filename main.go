package main

import (
	"github.com/nicovogelaar/protoc-gen-gofullmethods/internal"
	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	opt := protogen.Options{}
	internal.Run(opt)
}
