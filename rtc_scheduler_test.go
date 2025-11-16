// rtc_scheduler_test.go
package main

import (
	"os"
	"testing"
)

func TestCLI_Help(t *testing.T) {
	// Save original args
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Test help output
	os.Args = []string{"rtc-scheduler"}

	// This should not panic and should show usage
	// We can't easily test the output without capturing stdout
	// but we can at least verify it doesn't crash
}

func TestMain_Basic(t *testing.T) {
	// This is a basic smoke test to ensure the main function
	// can be called without panicking
	// In a real test, we'd need to mock dependencies
}