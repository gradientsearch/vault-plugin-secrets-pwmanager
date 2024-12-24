package secretsengine

import (
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// pathCredentials extends the Vault API with a `/creds`
// endpoint for a role. You can choose whether
// or not certain attributes should be displayed,
// required, and named.
func pathCredentials(b *pwmgrBackend) *framework.Path {
	return &framework.Path{
		Pattern: "creds/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeLowerCaseString,
				Description: "Name of the role",
				Required:    true,
			},
		},
		Callbacks:       map[logical.Operation]framework.OperationFunc{},
		HelpSynopsis:    pathCredentialsHelpSyn,
		HelpDescription: pathCredentialsHelpDesc,
	}
}

const pathCredentialsHelpSyn = `
Generate a Pwmgr API token from a specific Vault role.
`

const pathCredentialsHelpDesc = `
This path generates a Pwmgr API user tokens
based on a particular role. A role can only represent a user token,
since Pwmgr doesn't have other types of tokens.
`
