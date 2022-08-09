package main

import (
	"fmt"
	"strings"

	kubewarden "github.com/kubewarden/policy-sdk-go"
	easyjson "github.com/mailru/easyjson"
)

// IsNameSafe function looks for and match the podname with a slice of SafeNames
// If the podname contains a subsrting with the settings.SafeNames[n], the controller
// will allow the pod creation
func (s *Settings) IsNameSafe(podName string) bool {

	// If safename is declared, you are whitelisting by names
	if len(s.SafeNames) == 0 {
		return true
	}

	// Looks for the PodName in the Safe Pod Names
	for _, sn := range s.SafeNames {
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

func (s *Settings) validateSettings(payload []byte) ([]byte, error) {
	logger.Info("validating settings")

	settings := Settings{}
	err := easyjson.Unmarshal(payload, &settings)
	if err != nil {
		return kubewarden.RejectSettings(kubewarden.Message(fmt.Sprintf("Provided settings are not valid: %v", err)))
	}

	valid, err := settings.Valid()
	if err != nil {
		return kubewarden.RejectSettings(kubewarden.Message(fmt.Sprintf("Provided settings are not valid: %v", err)))
	}
	if valid {
		return kubewarden.AcceptSettings()
	}

	logger.Warn("rejecting settings")
	return kubewarden.RejectSettings(kubewarden.Message("Provided settings are not valid"))
}
