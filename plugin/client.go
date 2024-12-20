package secretsengine

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/hashicorp/go-hclog"
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
	URL      string
	Token    string
	client   *http.Client
	renew    chan (interface{})
	done     chan (interface{})
	roleID   string
	secretID string
	store    map[string]string
	logger   hclog.Logger
}

type SignInResponse struct {
	UserID int
	Token  string
}

type Auth struct {
	Renewable     bool              `json:"renewable"`
	LeaseDuration int               `json:"lease_duration"`
	Metadata      map[string]string `json:"metadata"`
	TokenPolicies []string          `json:"token_policies"`
	Accessor      string            `json:"accessor"`
	ClientToken   string            `json:"client_token"`
}
type AuthRoleLoginResponse struct {
	Auth          Auth        `json:"auth"`
	Warnings      interface{} `json:"warnings"`
	WrapInfo      interface{} `json:wrap_info""`
	Data          interface{} `json:"data"`
	LeaseDuration int         `json:"lease_duration"`
	Renewable     bool        `json:"renewable"`
	LeaseID       string      `json:"lease_id"`
}

func (p *pwmgrClient) renewLoop() {
	// TODO parameterize
	t := time.NewTicker(45 * time.Minute)
	for {
		select {
		case <-p.renew:
			if err := p.Login(); err != nil {
				p.logger.Error(err.Error())
			}
			t.Reset(45 * time.Minute)
		case <-t.C:
			p.renew <- nil

		case <-p.done:
			return
		}
	}
}

func (p *pwmgrClient) SignOut() error { return nil }

func (p *pwmgrClient) Login() error {
	url := fmt.Sprintf("%s/v1/auth/approle/login", p.URL)

	cfg := struct {
		RoleID   string `json:"role_id"`
		SecretID string `json:"secret_id"`
	}{
		RoleID:   p.roleID,
		SecretID: p.secretID,
	}

	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(cfg); err != nil {
		return fmt.Errorf("encode data: %w", err)
	}

	req, err := http.NewRequestWithContext(context.TODO(), http.MethodPost, url, &b)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("do: %w", err)
	}
	p.logger.Debug("client approle login successful")

	var response AuthRoleLoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("json decode: %w", err)
	}
	p.logger.Debug("client AuthRoleLoginResponse", "auth", fmt.Sprintf("%+v", response))

	p.Token = response.Auth.ClientToken

	return nil
}

// newClient creates a new client to access Pwmgr
// and exposes it for any secrets or roles to use.
func newClient(config *pwmgrConfig, logger hclog.Logger) (*pwmgrClient, error) {
	if config == nil {
		return nil, errors.New("client configuration was nil")
	}

	if config.RoleID == "" {
		return nil, errors.New("client role_id was not defined")
	}

	if config.SecretID == "" {
		return nil, errors.New("client secret_id was not defined")
	}

	if config.URL == "" {
		return nil, errors.New("client url was not defined")
	}

	pc := pwmgrClient{
		URL:      config.URL,
		client:   &defaultClient,
		store:    make(map[string]string),
		renew:    make(chan (interface{})),
		done:     make(chan (interface{})),
		roleID:   config.RoleID,
		secretID: config.SecretID,
		logger:   logger,
	}

	/* TODO add clean up logic and signaling*/
	go pc.renewLoop()
	return &pc, nil
}
