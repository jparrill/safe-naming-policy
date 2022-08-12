package main

import (
	"strings"

	"github.com/francoispqt/onelog"
	"github.com/kubewarden/gjson"
	kubewarden "github.com/kubewarden/policy-sdk-go"
	kp "github.com/kubewarden/policy-sdk-go/protocol"
	easyjson "github.com/mailru/easyjson"
)

func validate(payload []byte) ([]byte, error) {
	//
	validationRequest := kp.ValidationRequest{}
	err := easyjson.Unmarshal(payload, &validationRequest)
	if err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(err.Error()),
			kubewarden.Code(400))
	}

	podName := gjson.GetBytes(payload, "request.object.metadata.name").Str
	podNs := gjson.GetBytes(payload, "request.object.metadata.namespace").Str

	// Create a Settings instance from the ValidationRequest object
	settings, err := NewSettingsFromValidationReq(&validationRequest)
	if err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(err.Error()),
			kubewarden.Code(400))
	}

	logger.DebugWithFields("validating pod object", func(e onelog.Entry) {
		e.String("name", podName)
		e.String("namespace", podNs)
		//		e.String("%s", string(pod))
	})

	// From here starts the rules validations
	// Is this the right Namespace?
	if !(settings.IsThisMyNamespace(podNs)) {
		return kubewarden.AcceptRequest()
	}

	// The podName is blacklisted?
	if settings.IsNameUnsafe(podName) {
		logger.InfoWithFields("The pod creation is not allowed by a Kubewarden policy", func(e onelog.Entry) {
			e.String("name", podName)
			e.String("unsafe_names", strings.Join(settings.UnsafeNames, ","))
		})

		return kubewarden.RejectRequest(
			kubewarden.Message("The Unsafe pattern provided is blacklisted"),
			kubewarden.NoCode)
	}

	// The podName is whitelisted?
	if !(settings.IsNameSafe(podName)) {
		logger.InfoWithFields("The pod creation is not allowed by a Kubewarden policy", func(e onelog.Entry) {
			e.String("name", podName)
			e.String("safe_names", strings.Join(settings.SafeNames, ","))
		})

		return kubewarden.RejectRequest(
			kubewarden.Message("The Safe pattern provided is not whitelisted"),
			kubewarden.NoCode)
	}

	return kubewarden.AcceptRequest()
}
