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
$ protoc -I example service.proto --gofullmethods_out=:example --gofullmethods_opt=paths=source_relative
```

## Example

See here an example of the generated file: [service_methods.pb.go](./example/service_methods.pb.go)

```go
package example

const (
	Method_helloworld_Greeter__SayHello = "/helloworld.Greeter/SayHello"
	Method_helloworld_Greeter__SayBye   = "/helloworld.Greeter/SayBye"
)
```
