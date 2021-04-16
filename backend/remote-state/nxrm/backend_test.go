package nxrm

import (
	"testing"

	"github.com/hashicorp/terraform/backend"
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
