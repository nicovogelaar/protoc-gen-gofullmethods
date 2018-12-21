# protoc-gen-gofullmethods [![Build Status](https://travis-ci.org/nicovogelaar/protoc-gen-gofullmethods.svg?branch=master)](https://travis-ci.org/nicovogelaar/protoc-gen-gofullmethods)

This is a protoc go plugin to generate constants for all gRPC server methods.

## Use case

A use case could be to add some logic for a certain rpc method. For example, adding a condition for a certain rpc method in a gRPC server middleware. See here an example on line 33: [example/server/server.go](./example/server/server.go)

## Usage

### Install
```
$ go get -u github.com/nicovogelaar/protoc-gen-gofullmethods
$ go install github.com/nicovogelaar/protoc-gen-gofullmethods
```

### Generate

Include the `--gofullmethods_out=` flag to generate the `.fullmethods.pb.go` file.

```
$ protoc -I example service.proto --go_out=plugins=grpc:example --gofullmethods_out=example
```

## Example

See here an example of the generated file: [service.fullmethods.pb.go](./example/service.fullmethods.pb.go)

```go
package example

const (
	Greeter_SayHello = "/helloworld.Greeter/SayHello"
	Greeter_SayBye   = "/helloworld.Greeter/SayBye"
)

var (
	FullMethods = []string{
		Greeter_SayHello,
		Greeter_SayBye,
	}
)
```
