.PHONY: ci
ci: lint install gen test

.PHONY: install
install:
	go install .

.PHONY: gen
gen:
	protoc -I example service.proto --go_out=plugins=grpc:example --gofullmethods_out=example

.PHONY: test
test:
	go test -v -mod=vendor ./...

.PHONY: lint
lint:
	golangci-lint run -v
