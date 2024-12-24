package secretsengine

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// pwmgrUserEntry defines the data required
// for a Vault user to access and call the Pwmgr
// token endpoints
type pwmgrUserEntry struct {
	Username string        `json:"username"`
	UserID   int           `json:"user_id"`
	Token    string        `json:"token"`
	TokenID  string        `json:"token_id"`
	TTL      time.Duration `json:"ttl"`
	MaxTTL   time.Duration `json:"max_ttl"`
}

// toResponseData returns response data for a user
func (r *pwmgrUserEntry) toResponseData() map[string]interface{} {
	respData := map[string]interface{}{
		"ttl":      r.TTL.Seconds(),
		"max_ttl":  r.MaxTTL.Seconds(),
		"username": r.Username,
	}
	return respData
}

// pathUser extends the Vault API with a `/user`
// endpoint for the backend. You can choose whether
// or not certain attributes should be displayed,
// required, and named. You can also define different
// path patterns to list all users.
func pathUser(b *pwmgrBackend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "user/" + framework.GenericNameRegex("name"),
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeLowerCaseString,
					Description: "Name of the user",
					Required:    true,
				},
				"username": {
					Type:        framework.TypeString,
					Description: "The username for the Pwmgr product API",
					Required:    true,
				},
				"ttl": {
					Type:        framework.TypeDurationSecond,
					Description: "Default lease for generated credentials. If not set or set to 0, will use system default.",
				},
				"max_ttl": {
					Type:        framework.TypeDurationSecond,
					Description: "Maximum time for user. If not set or set to 0, will use system default.",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathUsersRead,
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.pathUsersWrite,
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.pathUsersWrite,
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.pathUsersDelete,
				},
			},
			HelpSynopsis:    pathUserHelpSynopsis,
			HelpDescription: pathUserHelpDescription,
		},
		{
			Pattern: "user/?$",
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: b.pathUsersList,
				},
			},
			HelpSynopsis:    pathUserListHelpSynopsis,
			HelpDescription: pathUserListHelpDescription,
		},
	}
}

// pathUsersList makes a request to Vault storage to retrieve a list of users for the backend
func (b *pwmgrBackend) pathUsersList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List(ctx, "user/")
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(entries), nil
}

// pathUsersRead makes a request to Vault storage to read a user and return response data
func (b *pwmgrBackend) pathUsersRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entry, err := b.getUser(ctx, req.Storage, d.Get("name").(string))
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: entry.toResponseData(),
	}, nil
}

// pathUsersWrite makes a request to Vault storage to update a user based on the attributes passed to the user configuration
func (b *pwmgrBackend) pathUsersWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name, ok := d.GetOk("name")
	if !ok {
		return logical.ErrorResponse("missing user name"), nil
	}

	userEntry, err := b.getUser(ctx, req.Storage, name.(string))
	if err != nil {
		return nil, err
	}

	if userEntry == nil {
		userEntry = &pwmgrUserEntry{}
	}

	createOperation := (req.Operation == logical.CreateOperation)

	if username, ok := d.GetOk("username"); ok {
		userEntry.Username = username.(string)
	} else if !ok && createOperation {
		return nil, fmt.Errorf("missing username in user")
	}

	if ttlRaw, ok := d.GetOk("ttl"); ok {
		userEntry.TTL = time.Duration(ttlRaw.(int)) * time.Second
	} else if createOperation {
		userEntry.TTL = time.Duration(d.Get("ttl").(int)) * time.Second
	}

	if maxTTLRaw, ok := d.GetOk("max_ttl"); ok {
		userEntry.MaxTTL = time.Duration(maxTTLRaw.(int)) * time.Second
	} else if createOperation {
		userEntry.MaxTTL = time.Duration(d.Get("max_ttl").(int)) * time.Second
	}

	if userEntry.MaxTTL != 0 && userEntry.TTL > userEntry.MaxTTL {
		return logical.ErrorResponse("ttl cannot be greater than max_ttl"), nil
	}

	if err := setUser(ctx, req.Storage, name.(string), userEntry); err != nil {
		return nil, err
	}

	return nil, nil
}

// pathUsersDelete makes a request to Vault storage to delete a user
func (b *pwmgrBackend) pathUsersDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete(ctx, "user/"+d.Get("name").(string))
	if err != nil {
		return nil, fmt.Errorf("error deleting pwmgr user: %w", err)
	}

	return nil, nil
}

// setUser adds the user to the Vault storage API
func setUser(ctx context.Context, s logical.Storage, name string, userEntry *pwmgrUserEntry) error {
	entry, err := logical.StorageEntryJSON("user/"+name, userEntry)
	if err != nil {
		return err
	}

	if entry == nil {
		return fmt.Errorf("failed to create storage entry for user")
	}

	if err := s.Put(ctx, entry); err != nil {
		return err
	}

	return nil
}

// getUser gets the user from the Vault storage API
func (b *pwmgrBackend) getUser(ctx context.Context, s logical.Storage, name string) (*pwmgrUserEntry, error) {
	if name == "" {
		return nil, fmt.Errorf("missing user name")
	}

	entry, err := s.Get(ctx, "user/"+name)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	var user pwmgrUserEntry

	if err := entry.DecodeJSON(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

const (
	pathUserHelpSynopsis    = `Manages the Vault user for generating Pwmgr tokens.`
	pathUserHelpDescription = `
This path allows you to read and write users used to generate Pwmgr tokens.
You can configure a user to manage a user's token by setting the username field.
`

	pathUserListHelpSynopsis    = `List the existing users in Pwmgr backend`
	pathUserListHelpDescription = `Users will be listed by the user name.`
)
