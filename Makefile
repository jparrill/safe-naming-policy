SOURCE_FILES := $(shell find . -type f -name '*.go')

build: types_easyjson.go policy.wasm annotated-policy.wasm

policy.wasm: $(SOURCE_FILES) go.mod go.sum types_easyjson.go
	tinygo build -o policy.wasm -target=wasi -no-debug .

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
	bats e2e.bats

.PHONY: clean
clean:
	go clean
	rm -f policy.wasm annotated-policy.wasm
