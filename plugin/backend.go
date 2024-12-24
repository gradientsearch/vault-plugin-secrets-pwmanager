package secretsengine

import (
	"context"
	"strings"
	"sync"

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

// pwmgrBackend defines an object that
// extends the Vault backend and stores the
// target API's client.
type pwmgrBackend struct {
	*framework.Backend
	lock   sync.RWMutex
	logger hclog.Logger
}

// backend defines the target API backend
// for Vault. It must include each path
// and the secrets it will store.
func backend() *pwmgrBackend {
	var b = pwmgrBackend{}
	appLogger := hclog.New(&hclog.LoggerOptions{
		Name:  "pwmgr",
		Level: hclog.LevelFromString("DEBUG"),
	})

	b.logger = appLogger
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),
		PathsSpecial: &logical.Paths{
			LocalStorage: []string{},
			SealWrapStorage: []string{
				"config",
				"role/*",
			},
		},
		Paths: framework.PathAppend(
			pathUser(&b),
		),
		BackendType: logical.TypeLogical,
		Invalidate:  b.invalidate,
	}

	return &b
}

// reset clears any client configuration for a new
// backend to be configured
func (b *pwmgrBackend) reset() {
	b.lock.Lock()
	defer b.lock.Unlock()
}

// invalidate clears an existing client configuration in
// the backend
func (b *pwmgrBackend) invalidate(ctx context.Context, key string) {
	if key == "config" {
		b.reset()
	}
}

// backendHelp should contain help information for the backend
const backendHelp = `
The Pwmgr secrets backend dynamically generates user tokens.
After mounting this backend, credentials to manage Pwmgr user tokens
must be configured with the "config/" endpoints.
`
