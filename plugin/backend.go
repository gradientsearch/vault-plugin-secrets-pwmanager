package secretsengine

import (
	"context"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := backend()

	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}

	return b, nil
}

// pwManagerBackend defines an object that
// extends the Vault backend and stores the
// target API's client.
type pwManagerBackend struct {
	*framework.Backend
	logger hclog.Logger
}

// backend defines the target API backend
// for Vault. It must include each path
// and the secrets it will store.
func backend() *pwManagerBackend {
	var b = pwManagerBackend{}
	appLogger := hclog.New(&hclog.LoggerOptions{
		Name:  "pwManager",
		Level: hclog.LevelFromString("DEBUG"),
	})

	b.logger = appLogger
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),
		PathsSpecial: &logical.Paths{
			LocalStorage: []string{},
		},
		Paths: framework.PathAppend(
			pathUser(&b),
		),
		BackendType: logical.TypeLogical,
	}

	return &b
}

// backendHelp should contain help information for the backend
const backendHelp = `
The PwManager secrets backend manages access to encrypted user password vaults.
After mounting this backend, users can register with this mount and start creating
password vaults and adding secrets to those vaults via end-to-end encryption using 2KSD.
`
