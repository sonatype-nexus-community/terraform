package nxrm

import (
	"testing"

	"github.com/hashicorp/terraform/backend"
	"github.com/hashicorp/terraform/configs/hcl2shim"
)

func TestBackend_impl(t *testing.T) {
	var _ backend.Backend = new(Backend)
}

func TestBackendConfig(t *testing.T) {
	b := backend.TestBackendConfig(t, New(), backend.TestWrapConfig(config)).(*Backend)

	if b.client.userName != config["username"] {
		t.Fatalf("Incorrect userName populated")
	}
	if b.client.password != config["password"] {
		t.Fatalf("Incorrect password populated")
	}
	if b.client.url != config["url"] {
		t.Fatalf("Incorrect url populated")
	}

	if b.client.subpath != config["subpath"] {
		t.Fatalf("Incorrect subpath populated")
	}

	if b.client.stateName != config["stateName"] {
		t.Fatalf("Incorrect stateName populated")
	}

	if b.client.timeout != config["timeout"] {
		t.Fatalf("Incorrect timeout populated")
	}
}

func TestBackendConfig_invalidSubpath(t *testing.T) {
	// sanitize null values for Go and break `subpath`
	cfg := hcl2shim.HCL2ValueFromConfigValue(map[string]interface{}{
		"username":  config["username"],
		"password":  config["password"],
		"url":       config["url"],
		"subpath":   "/this/here", // forward slash error
		"stateName": config["stateName"],
		"timeout":   config["timeout"],
	})

	_, diags := New().PrepareConfig(cfg)
	if !diags.HasErrors() {
		t.Fatal("expected config validation error")
	}
}
