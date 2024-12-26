package secretsengine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
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

// This provides a default client configuration, but it's recommended
// this is replaced by the user with application specific settings using
// the WithClient function at the time a GraphQL is constructed.
var defaultClient = http.Client{
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          1,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
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

	client *http.Client

	storage logical.Storage

	vaultToken string
	url        string
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
	b.client = &defaultClient

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

	// TODO parameterize
	t := time.NewTicker(45 * time.Minute)
	for {
		select {
		case <-p.renew:
			p.logger.Debug("renewing approle token")
			if err := p.Login(); err != nil {
				p.logger.Error(err.Error())
			} else {
				t.Reset(45 * time.Minute)
				p.logger.Debug("renewing approle token successful")
			}
		case <-t.C:
			p.renew <- nil

		case <-p.done:
			return
		}
	}
}

func (p *pwManagerBackend) Login() error {
	p.logger.Debug("starting app role request")
	config, err := getConfig(context.TODO(), p.storage)

	if config == nil || err != nil {
		return fmt.Errorf("pwmgr mount not configured. configure at /config")
	}
	p.logger.Debug("config not nil")
	p.logger.Debug(fmt.Sprintf("config %+v", config))
	url := fmt.Sprintf("%s/v1/auth/approle/login", config.URL)

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

	req, err := http.NewRequestWithContext(context.TODO(), http.MethodPost, url, &b)
	if err != nil {
		p.logger.Debug(fmt.Sprintf("error making app role request request: %s", err))
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		p.logger.Debug("error doing app role request")
		return fmt.Errorf("do: %w", err)
	}
	p.logger.Debug("client approle login successful")

	var response LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("json decode: %w", err)
	}
	p.logger.Debug("client AuthRoleLoginResponse", "auth", fmt.Sprintf("%+v", response))

	p.vaultToken = response.Auth.ClientToken
	p.url = config.URL

	return nil
}

// backendHelp should contain help information for the backend
const backendHelp = `
The PwManager secrets backend manages access to encrypted user password vaults.
After mounting this backend, users can register with this mount and start creating
password vaults and adding secrets to those vaults via end-to-end encryption using 2KSD.
`
