package main

import (
	"strings"

	"github.com/francoispqt/onelog"
	corev1 "github.com/kubewarden/k8s-objects/api/core/v1"
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

	// Create a Settings instance from the ValidationRequest object
	settings, err := NewSettingsFromValidationReq(&validationRequest)
	if err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(err.Error()),
			kubewarden.Code(400))
	}

	// Getting the skeleton of a request object
	podJSON := validationRequest.Request.Object

	// Pod instance inside of pod var
	pod := &corev1.Pod{}
	if err := easyjson.Unmarshal([]byte(podJSON), pod); err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message("Cannot decode Pod object"),
			kubewarden.Code(400))
	}

	logger.DebugWithFields("validating pod object", func(e onelog.Entry) {
		e.String("name", pod.Metadata.Name)
		e.String("namespace", pod.Metadata.Namespace)
		//		e.String("%s", string(pod))
	})

	// From here starts the rules validations
	// Is this the right Namespace?
	if !(settings.IsThisMyNamespace(pod.Metadata.Namespace)) {
		return kubewarden.AcceptRequest()
	}

	// The podName is blacklisted?
	if settings.IsNameUnsafe(pod.Metadata.Name) {
		logger.InfoWithFields("The pod creation is not allowed by a Kubewarden policy", func(e onelog.Entry) {
			e.String("name", pod.Metadata.Name)
			e.String("unsafe_names", strings.Join(settings.UnsafeNames, ","))
		})

		return kubewarden.RejectRequest(
			kubewarden.Message("The Unsafe pattern provided is blacklisted"),
			kubewarden.NoCode)
	}

	// The podName is whitelisted?
	if !(settings.IsNameSafe(pod.Metadata.Name)) {
		logger.InfoWithFields("The pod creation is not allowed by a Kubewarden policy", func(e onelog.Entry) {
			e.String("name", pod.Metadata.Name)
			e.String("safe_names", strings.Join(settings.SafeNames, ","))
		})

		return kubewarden.RejectRequest(
			kubewarden.Message("The Safe pattern provided is not whitelisted"),
			kubewarden.NoCode)
	}

	return kubewarden.AcceptRequest()
}
