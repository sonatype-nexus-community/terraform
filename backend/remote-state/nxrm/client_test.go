package nxrm

import (
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/terraform/backend"
)

func TestGetNXRMURL(t *testing.T) {
	b := backend.TestBackendConfig(t, New(), backend.TestWrapConfig(config)).(*Backend)

	url := b.client.getNXRMURL(b.client.stateName)

	if url != `http://localhost:8081/repository/tf-backend/this/here/demo.tfstate` {
		t.Fatalf("getNXRMURL mismatch: %s", url)
	}
}

func TestGetNXRMURLTrimUrl(t *testing.T) {
	cfg := config
	cfg["url"] = "http://localhost:8081/repository/tf-backend/"
	b := backend.TestBackendConfig(t, New(), backend.TestWrapConfig(cfg)).(*Backend)

	got := b.client.getNXRMURL(b.client.stateName)
	if got != `http://localhost:8081/repository/tf-backend/this/here/demo.tfstate` {
		t.Fatalf("getNXRMURL mismatch: %s", got)
	}
}

func TestGetNXRMURLTrimSubpathSuffix(t *testing.T) {
	cfg := config
	cfg["subpath"] = "this/here/"
	b := backend.TestBackendConfig(t, New(), backend.TestWrapConfig(cfg)).(*Backend)

	got := b.client.getNXRMURL(b.client.stateName)
	if got != `http://localhost:8081/repository/tf-backend/this/here/demo.tfstate` {
		t.Fatalf("getNXRMURL mismatch: %s", got)
	}
}

func TestGetNXRMURLTrimSubpathPrefix(t *testing.T) {
	cfg := config
	cfg["subpath"] = "/this/here"
	b := backend.TestBackendConfig(t, New(), backend.TestWrapConfig(cfg)).(*Backend)

	got := b.client.getNXRMURL(b.client.stateName)
	if got != `http://localhost:8081/repository/tf-backend/this/here/demo.tfstate` {
		t.Fatalf("getNXRMURL mismatch: %s", got)
	}
}

func TestGetHTTPClient(t *testing.T) {
	b := backend.TestBackendConfig(t, New(), backend.TestWrapConfig(config)).(*Backend)
	expectedTimeout := time.Second * time.Duration(config["timeout"].(int))

	got := b.client.getHTTPClient()
	if got.Timeout != expectedTimeout {
		t.Fatalf("getHTTPClient returned strange timeout")
	}
}

func TestGetRequest(t *testing.T) {
	b := backend.TestBackendConfig(t, New(), backend.TestWrapConfig(config)).(*Backend)

	req, err := b.client.getRequest(http.MethodGet, config["stateName"].(string), nil)
	if err != nil {
		t.Fatalf("getRequest error: %s", err)
	}
	u, p, ok := req.BasicAuth()
	if !ok {
		t.Fatalf("req.BasicAuth() not ok!")
	}

	if u != config["username"].(string) {
		t.Fatalf(mismatchError("username", u))
	}

	if p != config["password"].(string) {
		t.Fatalf(mismatchError("password", p))
	}
}
