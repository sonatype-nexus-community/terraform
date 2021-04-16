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
		t.Fatalf(mismatchError("userName", b.client.userName))
	}
	if b.client.password != config["password"] {
		t.Fatalf(mismatchError("password", b.client.password))
	}
	if b.client.url != config["url"] {
		t.Fatalf(mismatchError("url", b.client.url))
	}

	if b.client.subpath != config["subpath"] {
		t.Fatalf(mismatchError("subpath", b.client.subpath))
	}

	if b.client.stateName != config["stateName"] {
		t.Fatalf(mismatchError("stateName", b.client.stateName))
	}

	if b.client.timeout != config["timeout"] {
		t.Fatalf(mismatchError("timeout", b.client.timeout))
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
