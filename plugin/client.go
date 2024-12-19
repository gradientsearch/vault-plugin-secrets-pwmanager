package secretsengine

type Pwmgr struct {
	Token string
}

// pwmgrClient creates an object storing
// the client.
type pwmgrClient struct {
	Client Pwmgr
}

func (p *pwmgrClient) SignOut() error {
	return nil
}

// newClient creates a new client to access HashiCups
// and exposes it for any secrets or roles to use.
func newClient(config *pwmgrConfig) (*pwmgrClient, error) {
	return &pwmgrClient{}, nil
}
