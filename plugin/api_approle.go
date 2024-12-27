package secretsengine

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
)

// AppRole is used to perform AppRole operations on Vault.
type AppRole struct {
	c *api.Client
}

// AppRole is used to return the client for AppRole API calls.
func (c *pwmanagerClient) AppRole() *AppRole {
	return &AppRole{c: c.c}
}

func (c *AppRole) Enable(path string) error {
	r := c.c.NewRequest("POST", fmt.Sprintf("/v1/sys/auth/%s", path))
	r.Body = strings.NewReader(`{"type": "approle"}`)

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (c *AppRole) CreateRole(mount, roleName string, jsonData string) error {
	r := c.c.NewRequest("POST", fmt.Sprintf("/v1/auth/%s/role/%s", mount, roleName))
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

func (c *AppRole) RoleID(mount, roleName string) (RoleResponse, error) {
	r := c.c.NewRequest("GET", fmt.Sprintf("/v1/auth/%s/role/%s/role-id", mount, roleName))

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return RoleResponse{}, err
	}
	defer resp.Body.Close()

	var result RoleResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return RoleResponse{}, err
	}

	return result, nil
}

func (c *AppRole) SecretID(mount, roleName string, jsonData string) (SecretIDResponse, error) {
	r := c.c.NewRequest("POST", fmt.Sprintf("/v1/auth/%s/role/%s/secret-id", mount, roleName))
	r.Body = strings.NewReader(jsonData)

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return SecretIDResponse{}, err
	}
	defer resp.Body.Close()

	var result SecretIDResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return SecretIDResponse{}, err
	}

	return result, nil
}

func (c *AppRole) Login(mount, jsonData string) (LoginResponse, error) {
	r := c.c.NewRequest("POST", fmt.Sprintf("/v1/auth/%s/login", mount))
	r.Body = strings.NewReader(jsonData)

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

type RoleResponse struct {
	Data RoleData `json:"data"`
}
type RoleData struct {
	RoleID string `json:"role_id"`
}

type SecretIDResponse struct {
	Data SecretData `json:"data"`
}
type SecretData struct {
	SecretIDAccessor string `json:"secret_id_accessor"`
	SecretID         string `json:"secret_id"`
	SecretIDTTL      int    `json:"secret_id_ttl"`
	SecretIDNumUses  int    `json:"secret_id_num_uses"`
}
