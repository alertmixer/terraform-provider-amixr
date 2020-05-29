package amixr

import (
	"testing"
)

func TestConfigEmptyToken(t *testing.T) {
	config := Config{
		Token: "",
	}
	if _, err := config.Client(); err == nil {
		t.Fatalf("Expected error, but got nil")
	}
}

func TestConfigValidation(t *testing.T) {
	config := Config{
		Token: "foo",
	}
	if _, err := config.Client(); err != nil {
		t.Fatalf("Error: expected the client not fail with \"Token required\" error: %v", err)
	}
}
