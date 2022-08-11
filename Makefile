SOURCE_FILES := $(shell find . -type f -name '*.go')
SETTINGS ?= $(shell cat assets/settings.sample.json | jq -c)

build: types_easyjson.go policy.wasm annotated-policy.wasm
build_local: types_easyjson.go policy-arm.wasm annotated-policy.wasm

policy-arm.wasm: $(SOURCE_FILES) go.mod go.sum types_easyjson.go
	tinygo build -o policy.wasm -target=wasi -no-debug .

policy.wasm: $(SOURCE_FILES) go.mod go.sum types_easyjson.go
    docker run --rm -v ${PWD}:/src -w /src tinygo/tinygo:0.24.0 tinygo build -o policy.wasm -target=wasi -no-debug .

annotated-policy.wasm: policy.wasm assets/metadata.yml
	kwctl annotate -m assets/metadata.yml -o annotated-policy.wasm policy.wasm

.PHONY: generate-easyjson
types_easyjson.go: types.go
	docker run \
		--rm \
		-v ${PWD}:/src \
		-w /src \
		golang:1.17-alpine ./hack/generate-easyjson.sh
	go mod tidy

.PHONY: test
test: types_easyjson.go
	go test -v

.PHONY: e2e-tests
e2e-tests: annotated-policy.wasm
	bats assets/e2e.bats

.PHONY: clean
clean:
	go clean
	rm -f policy.wasm annotated-policy.wasm

install:
	@echo "Loading ClusterAdmissionPolicy into Kubernetes with these Settings:"
	@echo "$(SETTINGS)"
	@echo
	envsubst < assets/deployment/ClusterAdmissionPolicy.yaml 
