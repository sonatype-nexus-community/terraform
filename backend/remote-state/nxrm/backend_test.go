package nxrm

import (
	"testing"

	"github.com/hashicorp/terraform/backend"
)

func TestBackend_impl(t *testing.T) {
	var _ backend.Backend = new(Backend)
}

func TestBackendConfig(t *testing.T) {
	cfg := InitTestConfig()
	b := backend.TestBackendConfig(t, New(), backend.TestWrapConfig(cfg)).(*Backend)

	if b.client.userName != cfg["username"] {
		t.Fatalf(mismatchError(cfg, "username", b.client.userName))
	}
	if b.client.password != cfg["password"] {
		t.Fatalf(mismatchError(cfg, "password", b.client.password))
	}
	if b.client.url != cfg["url"] {
		t.Fatalf(mismatchError(cfg, "url", b.client.url))
	}

	if b.client.subpath != cfg["subpath"] {
		t.Fatalf(mismatchError(cfg, "subpath", b.client.subpath))
	}

	if b.client.stateName != cfg["stateName"] {
		t.Fatalf(mismatchError(cfg, "stateName", b.client.stateName))
	}

	if b.client.timeout != cfg["timeout"] {
		t.Fatalf(mismatchError(cfg, "timeout", b.client.timeout))
	}
}
