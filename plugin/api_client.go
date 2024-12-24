// TODO refactor wit vault/api client
package secretsengine

import (
	"fmt"

	"github.com/hashicorp/vault/api"
	vault "github.com/hashicorp/vault/api"
)

// pwmanger client wrapper over the vault api client. I want to be able
// to extend the api client but can't since the api client exists in the
// vault repo.
type pwmanagerClient struct {
	c *api.Client
}

// New
func NewClient(token string, hostPort string) (*pwmanagerClient, error) {
	config := vault.DefaultConfig()
	config.Address = "http://" + hostPort

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Vault client: %v", err)
	}

	client.SetToken(token)
	if err != nil {
		return nil, fmt.Errorf("error connecting to vault: %s", err)
	}

	return &pwmanagerClient{c: client}, nil
}
