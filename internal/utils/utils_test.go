package utils

import (
	"fmt"
	"os"
	"testing"
)

func TestGetDlqFlagOrDefault(t *testing.T) {
	os.Unsetenv("LISTEN_DLQ")
	if val := GetDlqFlagOrDefault(); val != false {
		t.Errorf("Expected false when env variable is unset, got %v", val)
	}

	os.Setenv("LISTEN_DLQ", "true")
	if val := GetDlqFlagOrDefault(); val != true {
		t.Errorf("Expected true, got %v", val)
	}

	os.Setenv("LISTEN_DLQ", "false")
	if val := GetDlqFlagOrDefault(); val != false {
		t.Errorf("Expected false, got %v", val)
	}

	os.Setenv("LISTEN_DLQ", "notabool")
	if val := GetDlqFlagOrDefault(); val != false {
		t.Errorf("Expected false on invalid bool, got %v", val)
	}
}

func TestAtoiMust(t *testing.T) {
	result := AtoiMust("123")
	if result != 123 {
		t.Errorf("Expected 123, got %d", result)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("AtoiMust did not panic on invalid input")
		}
	}()
	AtoiMust("notanumber")
}

func TestValues(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	values := MapValues(m)
	if len(values) != 2 {
		t.Errorf("Expected 2 values, got %d", len(values))
	}
	foundA, foundB := false, false
	for _, v := range values {
		if v == 1 {
			foundA = true
		}
		if v == 2 {
			foundB = true
		}
	}
	if !foundA || !foundB {
		t.Errorf("MapValues do not contain expected elements: %v", values)
	}
}

func TestWrapAndLog(t *testing.T) {
	err := fmt.Errorf("some error")
	wrappedErr := WrapAndLog(err)
	if wrappedErr == nil {
		t.Errorf("Expected wrapped error, got nil")
	}
	if wrappedErr.Error() == "" {
		t.Errorf("Wrapped error has empty message")
	}
	if !contains(wrappedErr.Error(), "some error") {
		t.Errorf("Wrapped error does not contain original message, got: %s", wrappedErr.Error())
	}
}

func contains(str, substr string) bool {
	return len(str) >= len(substr) && (str == substr || len(str) > len(substr) && (contains(str[1:], substr) || contains(str[:len(str)-1], substr)))
}
