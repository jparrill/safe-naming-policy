# safe-naming-policy

Kubewarden policy which validates the PodName in order to allow and disallow the Pod Deployment into the Kubernetes cluster

## Requirements

- tinygo
- kwctl

## Deployment

In order to deploy this policy in K8s, you just need a [K8s cluster](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/create-cluster-kubeadm/) + [Kubewarden](https://docs.kubewarden.io/quick-start#installation). Then you can go to `assets/deployment` folder and execute:

```bash
kubectl apply -f ClusterAdmissionPolicy-sample.yaml
```

And if you want to test if it's working, you can use the testing sample pods in the same folder:

- `kubectl apply -f invalidPod.yaml`: Will fail at deployment because of it's blacklisted
- `kubectl apply -f validPod.yaml`: Will success at deployment because it's whitelisted
- `kubectl apply -f notListedPod.yaml`: Will fail at deployment because it's not whitelisted


**Note:** You could change the `settings` field on the `ClusterAdmissionPolicy-sample.yaml` file in order to fit your needs.

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

  **UPDATE**: PR it's now merged, so in the next release from `v1.1.1` we will have Apple Silicon binaries 🎉🎉🎉
