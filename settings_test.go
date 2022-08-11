package main

import (
	"testing"

	easyjson "github.com/mailru/easyjson"
	"github.com/stretchr/testify/assert"
)

func TestSettingsSafe(t *testing.T) {
	UnsafePodName := "party"
	SafePodName := "lol"
	rawSettings := []byte(`{"namespace":"default","unsafe_names":["insecure-","chanchito"],"safe_names":["lol"]}`)
	settings := &Settings{}
	if err := easyjson.Unmarshal(rawSettings, settings); err != nil {
		t.Errorf("Unexpected error %+v", err)
	}
	assert.Equal(t, settings.IsNameSafe(UnsafePodName), false)
	assert.Equal(t, settings.IsNameSafe(SafePodName), true)
}

func TestSettingsUnsafe(t *testing.T) {
	UnsafePodName := "insecure"
	NotInUnsafe := "test"
	rawSettings := []byte(`{"namespace":"default","unsafe_names":["insecure","chanchito"]}`)
	settings := &Settings{}
	if err := easyjson.Unmarshal(rawSettings, settings); err != nil {
		t.Errorf("Unexpected error %+v", err)
	}
	assert.Equal(t, settings.IsNameUnsafe(UnsafePodName), true)
	assert.Equal(t, settings.IsNameUnsafe(NotInUnsafe), false)
}

func TestSettingsValidNoLists(t *testing.T) {
	rawSettings := []byte(`{"namespace":"default"}`)
	settings := &Settings{}
	if err := easyjson.Unmarshal(rawSettings, settings); err != nil {
		t.Errorf("Unexpected error %+v", err)
	}
	result := settings.Valid()
	// Fails because we need at least one of the 2 statements (blacklist|whitelist)
	assert.Equal(t, result, false)
}

func TestSettingsValidNoNS(t *testing.T) {
	rawSettings := []byte(`{}`)
	settings := &Settings{}
	if err := easyjson.Unmarshal(rawSettings, settings); err != nil {
		t.Errorf("Unexpected error %+v", err)
	}
	result := settings.Valid()
	// Fails because we need at least the NS and 1 statement
	assert.Equal(t, result, false)
}
