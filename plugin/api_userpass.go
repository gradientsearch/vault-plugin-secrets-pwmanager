package secretsengine

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/api"
)

// Userpass is used to perform Userpass operations on Vault.
type Userpass struct {
	c *api.Client
}

// Userpass is used to return the client for userpass API calls.
func (c *pwmanagerClient) Userpass() *Userpass {
	return &Userpass{c: c.c}
}

func (c *Userpass) User(mount, path string, userInfo UserInfo) error {
	r := c.c.NewRequest("POST", fmt.Sprintf("/v1/auth/%s/users/%s", mount, path))
	if err := r.SetJSONBody(userInfo); err != nil {
		return err
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

type UserInfo struct {
	Password      string   `json:"password"`
	TokenPolicies []string `json:"token_policies"`
}
