package secretsengine

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
)

// PwManager is used to perform PwManager operations on Vault.
type PwManager struct {
	c *api.Client
}

// PwManager is used to return the client for PwManager API calls.
func (c *pwmanagerClient) PwManager() *PwManager {
	return &PwManager{c: c.c}
}

func (c *PwManager) Config(mount, jsonData string) error {
	r := c.c.NewRequest("POST", fmt.Sprintf("/v1/%s/config", mount))
	r.Body = strings.NewReader(jsonData)

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
