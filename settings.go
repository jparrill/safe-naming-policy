package main

import (
	"strings"

	kubewarden "github.com/kubewarden/policy-sdk-go"
	kp "github.com/kubewarden/policy-sdk-go/protocol"
	"github.com/mailru/easyjson"
)

// Valid function validates the case when you recovers the settings.json into an struct
func (s *Settings) Valid() bool {
	if s.Namespace == "" {
		return false
	}

	if len(s.UnsafeNames) == 0 && len(s.SafeNames) == 0 {
		return false
	}

	return true
}

// Check if this is the namespace where should I look at.
// If true, continue with other evaluations
// If false, AcceptRequest
func (s *Settings) IsThisMyNamespace(podNamespace string) bool {

	if podNamespace == s.Namespace {
		return true
	}

	return false
}

// IsNameSafe function looks for and match the podname with a slice of SafeNames
// If the podname contains a subsrting with the settings.SafeNames[n], the controller
// will allow the pod creation
func (s *Settings) IsNameSafe(podName string) bool {

	var DefaultSafePrefixes = []string{
		"kube",
		"cert-manager",
		"local-path-provisioner",
		"coredns",
		"traefik",
		"metrics-server",
		"policy-server-default",
	}
	reservedBytes := len(DefaultSafePrefixes) + len(s.SafeNames)
	var allSafeNames = make([]string, reservedBytes, reservedBytes)

	allSafeNames = append(s.SafeNames, DefaultSafePrefixes...)

	// If safename is declared, you are whitelisting by names
	if len(s.SafeNames) == 0 {
		return true
	}

	// Looks for the PodName in the Safe Pod Names
	for _, sn := range allSafeNames {
		if strings.Contains(podName, sn) {
			return true
		}
	}

	return false
}

// IsNameUnsafe function looks for and match the podname with a slice of UnsafeNames
// If the podname contains a subsrting with the settings.UnsafeNames[n], the controller
// will deny the pod creation
func (s *Settings) IsNameUnsafe(podName string) bool {

	// If unsafename is declared, you are blacklisting by names
	if len(s.UnsafeNames) == 0 {
		return false
	}

	// Looks for the PodName in the Unsafe Pod Names
	for _, un := range s.UnsafeNames {
		if strings.Contains(podName, un) {
			return true
		}
	}

	return false
}

func NewSettingsFromValidationReq(validationReq *kp.ValidationRequest) (Settings, error) {
	settings := Settings{}
	err := easyjson.Unmarshal(validationReq.Settings, &settings)
	return settings, err
}

func validateSettings(payload []byte) ([]byte, error) {
	logger.Info("validating settings")

	settings := Settings{}
	err := easyjson.Unmarshal(payload, &settings)
	if err != nil {
		return kubewarden.RejectSettings(kubewarden.Message("Provided settings are not valid"))
	}

	valid := settings.Valid()
	if err != nil {
		return kubewarden.RejectSettings(kubewarden.Message("Provided settings are not valid"))
	}
	if valid {
		return kubewarden.AcceptSettings()
	}

	logger.Warn("rejecting settings")
	return kubewarden.RejectSettings(kubewarden.Message("Provided settings are not valid"))
}
