package secretsengine

import (
	"context"
	"encoding/json"
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

func (c *Userpass) Login(mount, path string, userInfo UserInfo) (LoginResponse, error) {
	r := c.c.NewRequest("POST", fmt.Sprintf("/v1/auth/%s/login/%s", mount, path))
	if err := r.SetJSONBody(userInfo); err != nil {
		return LoginResponse{}, err
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return LoginResponse{}, err
	}
	defer resp.Body.Close()

	var result LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return LoginResponse{}, err
	}

	return result, nil
}

type UserInfo struct {
	Password      string   `json:"password"`
	TokenPolicies []string `json:"token_policies"`
}

type LoginResponse struct {
	RequestID     string `json:"request_id"`
	LeaseID       string `json:"lease_id"`
	Renewable     bool   `json:"renewable"`
	LeaseDuration int    `json:"lease_duration"`
	Data          any    `json:"data"`
	WrapInfo      any    `json:"wrap_info"`
	Warnings      any    `json:"warnings"`
	Auth          Auth   `json:"auth"`
}
type Metadata struct {
	Username string `json:"username"`
}
type Auth struct {
	ClientToken    string   `json:"client_token"`
	Accessor       string   `json:"accessor"`
	Policies       []string `json:"policies"`
	TokenPolicies  []string `json:"token_policies"`
	Metadata       Metadata `json:"metadata"`
	LeaseDuration  int      `json:"lease_duration"`
	Renewable      bool     `json:"renewable"`
	EntityID       string   `json:"entity_id"`
	TokenType      string   `json:"token_type"`
	Orphan         bool     `json:"orphan"`
	MfaRequirement any      `json:"mfa_requirement"`
	NumUses        int      `json:"num_uses"`
}
