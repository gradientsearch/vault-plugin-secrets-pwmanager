package secretsengine

import (
	"errors"
)

type Clinet struct {
	Token string
}

type SignInResponse struct {
	UserID int
	Token  string
}

// pwmgrClient creates an object storing
// the client.
type pwmgrClient struct {
	Client *Clinet
}

func (p *pwmgrClient) SignOut() error                  { return nil }
func (p *pwmgrClient) SignIn() (SignInResponse, error) { return SignInResponse{}, nil }

// newClient creates a new client to access Pwmgr
// and exposes it for any secrets or roles to use.
func newClient(config *pwmgrConfig) (*pwmgrClient, error) {
	if config == nil {
		return nil, errors.New("client configuration was nil")
	}

	if config.RoleID == "" {
		return nil, errors.New("client username was not defined")
	}

	if config.SecretID == "" {
		return nil, errors.New("client password was not defined")
	}

	if config.URL == "" {
		return nil, errors.New("client URL was not defined")
	}

	return &pwmgrClient{&Clinet{}}, nil
}
