package secretsengine

import (
	"context"
	"fmt"

	mapstructure "github.com/go-viper/mapstructure/v2"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// pwmgrRegisterEntry defines the data required
// for a Vault register to access and call the Pwmgr
// token endpoints
type pwmgrRegisterEntry struct {
	// uuid of priv key
	UUID string `json:"uuid"`
	// symmetric key used to encrypt the EncPriKey
	EncSymKey EncSymKey `json:"enc_sym_key"`
	// mp a.k.a secret key
	EncryptedBy string `json:"encrypted_by"`
	// priv key used to encrypt `Safe` data
	EncPriKey EncPriKey `json:"enc_pri_key"`
	// pub key of the private key
	PubKey map[string]string `json:"pubkey"`
}

// toResponseData returns response data for a register
func (r *pwmgrRegisterEntry) toResponseData() map[string]interface{} {
	respData := map[string]interface{}{}
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
			Pattern: "register",
			Fields: map[string]*framework.FieldSchema{
				"uuid": {
					Type:        framework.TypeLowerCaseString,
					Description: "unique id",
					Required:    true,
				},
				"enc_sym_key": {
					Type:        framework.TypeMap,
					Description: "encrypted key that encrypts the private key",
					Required:    true,
				},
				"encrypted_by": {
					Type:        framework.TypeString,
					Description: "UUID of the key that encrypts the encSymmetricKey",
					Required:    true,
				},
				"pubKey": {
					Type:        framework.TypeMap,
					Description: "public part of the key pair",
					Required:    true,
				},
				"enc_pri_key": {
					Type:        framework.TypeMap,
					Description: "private part of the key pair",
					Required:    true,
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

	registerEntry, err := b.getRegister(ctx, req.Storage, req.EntityID)
	if err != nil {
		return nil, err
	}

	if registerEntry == nil {
		registerEntry = &pwmgrRegisterEntry{}
	}

	createOperation := (req.Operation == logical.CreateOperation)

	if uuid, ok := d.GetOk("uuid"); ok {
		registerEntry.UUID = uuid.(string)
	} else if !ok && createOperation {
		return nil, fmt.Errorf("missing username in register")
	}

	if encSymKey, ok := d.GetOk("enc_sym_key"); ok {
		if err := mapstructure.Decode(encSymKey, &registerEntry.EncSymKey); err != nil {
			return logical.ErrorResponse("error decoding encSymKey"), nil
		}
	} else if createOperation {
		return logical.ErrorResponse("must have encSymKey"), nil
	}

	if encryptedBy, ok := d.GetOk("encrypted_by"); ok {
		registerEntry.EncryptedBy = encryptedBy.(string)
	} else if createOperation {
		return logical.ErrorResponse("must have encryptedBy"), nil
	}

	if encPriKey, ok := d.GetOk("enc_pri_key"); ok {
		if err := mapstructure.Decode(encPriKey, &registerEntry.EncPriKey); err != nil {
			return logical.ErrorResponse("error decoding encPriKey"), nil
		}
	} else if createOperation {
		return logical.ErrorResponse("must have encPriKey"), nil
	}

	if pubkey, ok := d.GetOk("pubKey"); ok {
		if err := mapstructure.Decode(pubkey, &registerEntry.PubKey); err != nil {
			return logical.ErrorResponse("error decoding encPriKey"), nil
		}

	} else if createOperation {
		return logical.ErrorResponse("must have pubkey"), nil
	}

	if err := setRegister(ctx, req.Storage, req.EntityID, registerEntry); err != nil {
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
func setRegister(ctx context.Context, s logical.Storage, entityId string, registerEntry *pwmgrRegisterEntry) error {
	entry, err := logical.StorageEntryJSON("register/"+entityId, registerEntry)
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
func (b *pwmgrBackend) getRegister(ctx context.Context, s logical.Storage, entityID string) (*pwmgrRegisterEntry, error) {
	if entityID == "" {
		return nil, fmt.Errorf("missing register entity ID")
	}

	entry, err := s.Get(ctx, "register/"+entityID)
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
