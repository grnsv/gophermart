package logger

import (
	"testing"
)

func TestNew(t *testing.T) {
	logger := New()
	if logger == nil {
		t.Fatal("Expected logger to be initialized, got nil")
	}

	err := logger.Sync()
	if err != nil {
		if err.Error() == "sync /dev/stderr: invalid argument" {
			t.Skipf("Skipping test due to known issue: %v", err)
		} else {
			t.Errorf("Expected no error on logger.Sync(), got %v", err)
		}
	}
}

func TestFake(t *testing.T) {
	logger := Fake()
	if logger == nil {
		t.Fatal("Expected fake logger to be initialized, got nil")
	}

	err := logger.Sync()
	if err != nil {
		t.Errorf("Expected no error on fake logger.Sync(), got %v", err)
	}
}
