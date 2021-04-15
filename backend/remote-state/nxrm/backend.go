package nxrm

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/backend"
	"github.com/hashicorp/terraform/internal/legacy/helper/schema"
	"github.com/hashicorp/terraform/states/remote"
	"github.com/hashicorp/terraform/states/statemgr"
)

func New() backend.Backend {
	s := &schema.Backend{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NXRM_USERNAME", nil),
				Description: "Username",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NXRM_PASSWORD", nil),
				Description: "Password",
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NXRM_URL", nil),
				Description: "NXRM Repo URL",
			},
			"subpath": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NXRM_SUBPATH", nil),
				Description: "NXRM Repo Subpath",
				ValidateFunc: func(v interface{}, s string) ([]string, []error) {
					// taken from backend/remote-state/s3
					if strings.HasPrefix(v.(string), "/") {
						return nil, []error{errors.New("subpath must not start with '/'")}
					}
					return nil, nil
				},
			},
			"state_name": {
				Type:        schema.TypeString,
				Required:    true,
				Default:     "terraform.tfstate",
				Description: "The desired name of the State file (defaults to terraform.tfstate)",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NXRM_CLIENT_TIMEOUT", 30),
				Description: "Timeout in seconds",
			},
		},
	}

	b := &Backend{Backend: s}
	b.Backend.ConfigureFunc = b.configure
	return b
}

type Backend struct {
	*schema.Backend

	client *NXRMClient
}

func (b *Backend) configure(ctx context.Context) error {
	data := schema.FromContextBackendConfig(ctx)

	userName := data.Get("username").(string)
	password := data.Get("password").(string)
	url := data.Get("url").(string)
	subpath := data.Get("subpath").(string)
	stateName := data.Get("state_name").(string)
	timeout := data.Get("timeout").(int)

	b.client = &NXRMClient{
		userName:        userName,
		password:        password,
		url:             url,
		subpath:         subpath,
		tfLockArtifact:  fmt.Sprintf("%s.lock.json", stateName),
		tfStateArtifact: stateName,
		timeoutSeconds:  timeout,
	}
	return nil
}

func (b *Backend) Workspaces() ([]string, error) {
	return nil, backend.ErrWorkspacesNotSupported
}

func (b *Backend) DeleteWorkspace(string) error {
	return backend.ErrWorkspacesNotSupported
}

func (b *Backend) StateMgr(name string) (statemgr.Full, error) {
	if name != backend.DefaultStateName {
		return nil, backend.ErrWorkspacesNotSupported
	}
	return &remote.State{
		Client: b.client,
	}, nil
}
