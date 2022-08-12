module github.com/jparrill/safe-naming-policy

go 1.19

require (
	github.com/francoispqt/onelog v0.0.0-20190306043706-8c2bb31b10a4
	github.com/kubewarden/k8s-objects v1.24.0-kw3
	github.com/kubewarden/policy-sdk-go v0.2.3
	github.com/mailru/easyjson v0.7.7
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.2.2
	github.com/wapc/wapc-guest-tinygo v0.3.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/deckarep/golang-set v1.8.0 // indirect
	github.com/francoispqt/gojay v0.0.0-20181220093123-f2cc13a668ca // indirect
	github.com/go-openapi/strfmt v0.0.0-00010101000000-000000000000 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/kubewarden/gjson v1.7.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/tidwall/match v1.0.3 // indirect
	github.com/tidwall/pretty v1.0.2 // indirect
)

replace github.com/go-openapi/strfmt => github.com/kubewarden/strfmt v0.1.0
