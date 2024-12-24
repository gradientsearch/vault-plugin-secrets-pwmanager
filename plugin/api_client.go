// TODO refactor wit vault/api client
package secretsengine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	vault "github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/logical"
)

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

// Client provides support to access Hashicorp's Client product for keys.
type pwmgrClient struct {
	URL     string
	Token   string
	client  *http.Client
	renew   chan (interface{})
	done    chan (interface{})
	store   map[string]string
	logger  hclog.Logger
	storage logical.Storage
}

type SignInResponse struct {
	UserID int
	Token  string
}

func (p *pwmgrClient) renewLoop() {

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

func (p *pwmgrClient) SignOut() error { return nil }

func (p *pwmgrClient) Login() error {
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

	p.Token = response.Auth.ClientToken
	p.URL = config.URL

	return nil
}

// newClient creates a new client to access Pwmgr
// and exposes it for any secrets or roles to use.
func newClient(storage logical.Storage, logger hclog.Logger) *pwmgrClient {
	pc := pwmgrClient{
		client:  &defaultClient,
		store:   make(map[string]string),
		renew:   make(chan (interface{})),
		done:    make(chan (interface{})),
		logger:  logger,
		storage: storage,
	}
	return &pc
}

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
