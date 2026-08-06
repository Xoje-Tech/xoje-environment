package doctor

import (
	"testing"
)

func TestGoVersionCheck(t *testing.T) {
	check := &GoVersionCheck{Required: "1.22"}
	result := check.Run()

	if result.Name != "Go Version" {
		t.Errorf("expected name 'Go Version', got %s", result.Name)
	}
}

func TestNetworkCheck(t *testing.T) {
	check := &NetworkCheck{Target: "google.com"}
	result := check.Run()

	if result.Name != "Network Connectivity" {
		t.Errorf("expected name 'Network Connectivity', got %s", result.Name)
	}
}

func TestConfigCheck(t *testing.T) {
	check := &ConfigCheck{Path: "config.yaml"}
	result := check.Run()

	if result.Name != "Configuration File" {
		t.Errorf("expected name 'Configuration File', got %s", result.Name)
	}
}
