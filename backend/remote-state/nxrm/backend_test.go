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
	config := map[string]interface{}{
		"username": "testymctestface",
		"password": "mybigsecret",
		"url":      "http://localhost:8081/repository/tf-backend",
		"subpath":  "tf",
		"state_name": "terraform.tfstate",
		"timeout":  30,
	}

	b := backend.TestBackendConfig(t, New(), backend.TestWrapConfig(config)).(*Backend)

	if b.client.userName != "testymctestface" {
		t.Fatalf("Incorrect userName was populated")
	}
	if b.client.password != "mybigsecret" {
		t.Fatalf("Incorrect password was populated")
	}
	if b.client.url != "http://localhost:8081/repository/tf-backend" {
		t.Fatalf("Incorrect url was populated")
	}

	if b.client.subpath != "tf" {
		t.Fatalf("Incorrect subpath was populated")
	}
}

func TestBackendConfig_invalidSubpath(t *testing.T) {
	cfg := hcl2shim.HCL2ValueFromConfigValue(map[string]interface{}{
		"username": "testymctestface",
		"password": "mybigsecret",
		"url":      "http://localhost:8081/repository/tf-backend",
		"subpath":  "/tf", // forward slash error
		"state_name": "terraform.tfstate",
		"timeout":  30,
	})

	_, diags := New().PrepareConfig(cfg)
	if !diags.HasErrors() {
		t.Fatal("expected config validation error")
	}
}
