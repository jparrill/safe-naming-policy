package main

import (
	"github.com/francoispqt/onelog"
	kubewarden "github.com/kubewarden/policy-sdk-go"
	wapc "github.com/wapc/wapc-guest-tinygo"
)

var logWriter = kubewarden.KubewardenLogWriter{}
var logger = onelog.New(&logWriter, onelog.ALL)

func main() {
	wapc.RegisterFunctions(wapc.Functions{
		"validate":          validate,
		"validate_settings": validateSettings,
	})
}
