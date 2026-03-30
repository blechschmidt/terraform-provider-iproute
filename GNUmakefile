default: build

build:
	go build -o terraform-provider-iproute

install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/example/iproute/0.1.0/linux_amd64
	cp terraform-provider-iproute ~/.terraform.d/plugins/registry.terraform.io/example/iproute/0.1.0/linux_amd64/

test:
	go test ./internal/... -v -timeout 30m

testacc:
	TF_ACC=1 go test ./internal/... -v -timeout 30m

coverage:
	TF_ACC=1 go test ./internal/... -v -timeout 30m -coverprofile=coverage.out -covermode=atomic
	go tool cover -func=coverage.out

lint:
	golangci-lint run ./...

.PHONY: default build install test testacc coverage lint
