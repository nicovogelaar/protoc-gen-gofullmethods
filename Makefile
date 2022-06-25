.PHONY: ci
ci: lint install gen test

.PHONY: install
install:
	go install .

.PHONY: gen
gen:
	protoc --go_out=./example --go_opt=paths=source_relative --go-grpc_out=./example --go-grpc_opt=paths=source_relative --gofullmethods_out=./example --gofullmethods_opt=paths=source_relative -I example service.proto

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	golangci-lint run -v
