# safe-naming-policy

Kubewarden policy which validates the PodName in order to allow and disallow the Pod Deployment into the Kubernetes cluster

## Requirements

- tinygo
- kwctl

## Development

- **types_easyjson.go**: Generates the EasyJson types
```
make types_easyjson.go
```

- **policy.wasm**: Generates the WASM file to be consumed by Kubewarden
```
make policy.wasm
```

- **annotated-policy.wasm**: Annotates the `policy.wasm` generated with the `assets/metadata.yml` details
```
make annotated-policy.wasm
```

You could perform all those actions using `make build`

- **e2e-tests**: Uses [Bats](https://github.com/bats-core/bats-core) framework to execute the E2E tests located in `assets/e2e.bats`
- **test**: Executes the Go tests
- **clean**: Typical `go clean` + deletion of the WASM files generated `policy.wasm` `annotated-policy.wasm`

## Caveats

### M1/M2 Apple Silicon processors

To enable the deployment on these processors we will need to do 2 things:

- Tinygo build:
  We cannot execute the build on a docker environment without make our own Tinygo container image, because they only support `amd64` and Apple Silicon ones are based on `aach64`

- Kwctl Binary:
  For now we need to recompile the `kwctl` binary by hand because they are not publishing images based on `aarch64` ([PR On going](https://github.com/kubewarden/kwctl/pull/278)). To do this we need to execute these commands (remember, to perform the compilation you will need Rust):
  
  ```bash
  git clone https://github.com/kubewarden/kwctl.git && cd kwctl
  rustup target add aarch64-apple-darwin
  cargo build --target=aarch64-apple-darwin --release
  ```
