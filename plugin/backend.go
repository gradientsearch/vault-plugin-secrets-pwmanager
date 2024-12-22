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
	lock    sync.RWMutex
	client  *pwmgrClient
	storage logical.Storage
	logger  hclog.Logger
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
			pathRole(&b),
			pathRegister(&b),
			[]*framework.Path{
				pathConfig(&b),
				pathCredentials(&b),
			},
		),
		Secrets: []*framework.Secret{
			b.pwmgrToken(),
		},
		BackendType:    logical.TypeLogical,
		Invalidate:     b.invalidate,
		InitializeFunc: b.initialize,
	}

	return &b
}

func (b *pwmgrBackend) initialize(ctx context.Context, req *logical.InitializationRequest) error {
	b.client = newClient(req.Storage, b.logger)

	b.storage = req.Storage
	go b.client.renewLoop()
	return nil
}

// reset clears any client configuration for a new
// backend to be configured
func (b *pwmgrBackend) reset() {
	b.lock.Lock()
	defer b.lock.Unlock()
	if b.client != nil {
		b.client.renew <- nil
	}
}

// invalidate clears an existing client configuration in
// the backend
func (b *pwmgrBackend) invalidate(ctx context.Context, key string) {
	if key == "config" {
		b.reset()
	}
}

// getClient locks the backend as it configures and creates a
// a new client for the target API
func (b *pwmgrBackend) getClient(ctx context.Context, s logical.Storage) (*pwmgrClient, error) {
	b.lock.RLock()
	unlockFunc := b.lock.RUnlock
	defer func() { unlockFunc() }()

	if b.client != nil {
		return b.client, nil
	}

	b.lock.RUnlock()
	b.lock.Lock()
	unlockFunc = b.lock.Unlock

	config, err := getConfig(ctx, s)
	if err != nil {
		return nil, err
	}

	if config == nil {
		config = new(pwmgrConfig)
	}

	if err != nil {
		return nil, err
	}

	return b.client, nil
}

// backendHelp should contain help information for the backend
const backendHelp = `
The Pwmgr secrets backend dynamically generates user tokens.
After mounting this backend, credentials to manage Pwmgr user tokens
must be configured with the "config/" endpoints.
`
