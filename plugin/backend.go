package secretsengine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

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

	// chans to renew app role
	renew chan (interface{})
	done  chan (interface{})

	c *pwmanagerClient

	storage logical.Storage
}

// backend defines the target API backend
// for Vault. It must include each path
// and the secrets it will store.
func backend() *pwManagerBackend {
	var b = pwManagerBackend{}

	b.renew = make(chan interface{})
	b.done = make(chan interface{})

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
			[]*framework.Path{
				pathConfig(&b),
			},
		),
		BackendType:    logical.TypeLogical,
		InitializeFunc: b.initialize,
	}

	go b.renewLoop()
	return &b
}

func (b *pwManagerBackend) initialize(ctx context.Context, req *logical.InitializationRequest) error {
	b.storage = req.Storage

	return nil
}

func (p *pwManagerBackend) renewLoop() {
	t := time.NewTicker(45 * time.Minute)
	for {
		select {
		case <-p.renew:
			if err := p.Login(); err != nil {
				p.logger.Error(err.Error())
			} else {
				t.Reset(45 * time.Minute)
			}
		case <-t.C:
			if err := p.Login(); err != nil {
				p.logger.Error(err.Error())
			}
		case <-p.done:
			return
		}
	}
}

func (p *pwManagerBackend) Login() error {
	config, err := getConfig(context.TODO(), p.storage)

	if config == nil || err != nil {
		return fmt.Errorf("pwmanager mount not configured. configure at /config")
	}

	p.c, err = NewClient("", config.URL)
	if err != nil {
		return fmt.Errorf("error configuring pwmanagerClient: %s", err)
	}

	cfg := struct {
		RoleID   string `json:"role_id"`
		SecretID string `json:"secret_id"`
	}{
		RoleID:   config.RoleID,
		SecretID: config.SecretID,
	}

	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(cfg); err != nil {
		return fmt.Errorf("encode data: %w", err)
	}

	// TODO config the app role mount
	response, err := p.c.AppRole().Login("approle", b.String())
	if err != nil {
		p.logger.Debug("error doing app role request")
		return fmt.Errorf("do: %w", err)
	}
	p.logger.Debug("client approle login successful")

	p.c.c.SetToken(response.Auth.ClientToken)

	return nil
}

// backendHelp should contain help information for the backend
const backendHelp = `
The PwManager secrets backend manages access to encrypted user password bundles.
After mounting this backend, users can register with this mount and start creating
password bundles and adding password entries to those bundles via end-to-end encryption using 2KSD.
`
