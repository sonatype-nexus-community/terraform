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
	testConfg := map[string]interface{}{
		"url":       "http://localhost:8081/repository/tf-backend/",
		"subpath":   "this/here",
		"username":  "",
		"password":  "",
		"stateName": "demo.tfstate",
	}
	b := backend.TestBackendConfig(t, New(), backend.TestWrapConfig(testConfg)).(*Backend)

	url := b.client.getNXRMURL(b.client.stateName)

	if url != `http://localhost:8081/repository/tf-backend/this/here/demo.tfstate` {
		t.Fatalf("getNXRMURL mismatch: %s", url)
	}
}

func TestGetHTTPClient(t *testing.T) {
	b := backend.TestBackendConfig(t, New(), backend.TestWrapConfig(config)).(*Backend)
	expectedTimeout := time.Second * time.Duration(config["timeout"].(int))

	c := b.client.getHTTPClient()

	if c.Timeout != expectedTimeout {
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
