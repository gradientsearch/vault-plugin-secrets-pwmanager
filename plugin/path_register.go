package secretsengine

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// pwmgrRegisterEntry defines the data required
// for a Vault register to access and call the Pwmgr
// token endpoints
type pwmgrRegisterEntry struct {
	Username string        `json:"username"`
	UserID   int           `json:"user_id"`
	Token    string        `json:"token"`
	TokenID  string        `json:"token_id"`
	TTL      time.Duration `json:"ttl"`
	MaxTTL   time.Duration `json:"max_ttl"`
}

// toResponseData returns response data for a register
func (r *pwmgrRegisterEntry) toResponseData() map[string]interface{} {
	respData := map[string]interface{}{
		"ttl":      r.TTL.Seconds(),
		"max_ttl":  r.MaxTTL.Seconds(),
		"username": r.Username,
	}
	return respData
}

// pathRegister extends the Vault API with a `/register`
// endpoint for the backend. You can choose whether
// or not certain attributes should be displayed,
// required, and named. You can also define different
// path patterns to list all registers.
func pathRegister(b *pwmgrBackend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "register/" + framework.GenericNameRegex("name"),
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeLowerCaseString,
					Description: "Name of the register",
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
					Description: "Maximum time for register. If not set or set to 0, will use system default.",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathRegistersRead,
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.pathRegistersWrite,
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.pathRegistersWrite,
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.pathRegistersDelete,
				},
			},
			HelpSynopsis:    pathRegisterHelpSynopsis,
			HelpDescription: pathRegisterHelpDescription,
		},
		{
			Pattern: "register/?$",
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: b.pathRegistersList,
				},
			},
			HelpSynopsis:    pathRegisterListHelpSynopsis,
			HelpDescription: pathRegisterListHelpDescription,
		},
	}
}

// pathRegistersList makes a request to Vault storage to retrieve a list of registers for the backend
func (b *pwmgrBackend) pathRegistersList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List(ctx, "register/")
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(entries), nil
}

// pathRegistersRead makes a request to Vault storage to read a register and return response data
func (b *pwmgrBackend) pathRegistersRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entry, err := b.getRegister(ctx, req.Storage, d.Get("name").(string))
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

// pathRegistersWrite makes a request to Vault storage to update a register based on the attributes passed to the register configuration
func (b *pwmgrBackend) pathRegistersWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name, ok := d.GetOk("name")
	if !ok {
		return logical.ErrorResponse("missing register name"), nil
	}

	registerEntry, err := b.getRegister(ctx, req.Storage, name.(string))
	if err != nil {
		return nil, err
	}

	if registerEntry == nil {
		registerEntry = &pwmgrRegisterEntry{}
	}

	createOperation := (req.Operation == logical.CreateOperation)

	if username, ok := d.GetOk("username"); ok {
		registerEntry.Username = username.(string)
	} else if !ok && createOperation {
		return nil, fmt.Errorf("missing username in register")
	}

	if ttlRaw, ok := d.GetOk("ttl"); ok {
		registerEntry.TTL = time.Duration(ttlRaw.(int)) * time.Second
	} else if createOperation {
		registerEntry.TTL = time.Duration(d.Get("ttl").(int)) * time.Second
	}

	if maxTTLRaw, ok := d.GetOk("max_ttl"); ok {
		registerEntry.MaxTTL = time.Duration(maxTTLRaw.(int)) * time.Second
	} else if createOperation {
		registerEntry.MaxTTL = time.Duration(d.Get("max_ttl").(int)) * time.Second
	}

	if registerEntry.MaxTTL != 0 && registerEntry.TTL > registerEntry.MaxTTL {
		return logical.ErrorResponse("ttl cannot be greater than max_ttl"), nil
	}

	if err := setRegister(ctx, req.Storage, name.(string), registerEntry); err != nil {
		return nil, err
	}

	return nil, nil
}

// pathRegistersDelete makes a request to Vault storage to delete a register
func (b *pwmgrBackend) pathRegistersDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete(ctx, "register/"+d.Get("name").(string))
	if err != nil {
		return nil, fmt.Errorf("error deleting pwmgr register: %w", err)
	}

	return nil, nil
}

// setRegister adds the register to the Vault storage API
func setRegister(ctx context.Context, s logical.Storage, name string, registerEntry *pwmgrRegisterEntry) error {
	entry, err := logical.StorageEntryJSON("register/"+name, registerEntry)
	if err != nil {
		return err
	}

	if entry == nil {
		return fmt.Errorf("failed to create storage entry for register")
	}

	if err := s.Put(ctx, entry); err != nil {
		return err
	}

	return nil
}

// getRegister gets the register from the Vault storage API
func (b *pwmgrBackend) getRegister(ctx context.Context, s logical.Storage, name string) (*pwmgrRegisterEntry, error) {
	if name == "" {
		return nil, fmt.Errorf("missing register name")
	}

	entry, err := s.Get(ctx, "register/"+name)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	var register pwmgrRegisterEntry

	if err := entry.DecodeJSON(&register); err != nil {
		return nil, err
	}
	return &register, nil
}

const (
	pathRegisterHelpSynopsis    = `Manages the Vault register for generating Pwmgr tokens.`
	pathRegisterHelpDescription = `
This path allows you to read and write registers used to generate Pwmgr tokens.
You can configure a register to manage a user's token by setting the username field.
`

	pathRegisterListHelpSynopsis    = `List the existing registers in Pwmgr backend`
	pathRegisterListHelpDescription = `Registers will be listed by the register name.`
)
