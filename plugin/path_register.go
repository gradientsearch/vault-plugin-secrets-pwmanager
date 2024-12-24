package secretsengine

import (
	"context"
	"fmt"

	mapstructure "github.com/go-viper/mapstructure/v2"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const USER_SCHEMA string = "users"

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

// pathRegister extends the Vault API with a `/register`
// endpoint for the backend.
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
			ExistenceCheck: b.pathUserExistenceCheck,
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.pathRegistersWrite,
				},
			},
			HelpSynopsis:    pathRegisterHelpSynopsis,
			HelpDescription: pathRegisterHelpDescription,
		},
	}
}

// pathRegistersWrite makes a request to Vault storage to register a users UUK.
func (b *pwmgrBackend) pathRegistersWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	registerEntry, err := b.getRegister(ctx, req.Storage, req.EntityID)
	if err != nil {
		return nil, err
	}

	if registerEntry != nil {
		return logical.ErrorResponse("user already registered"), nil
	}

	registerEntry = &pwmgrRegisterEntry{}

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

// pathConfigExistenceCheck verifies if the configuration exists.
func (b *pwmgrBackend) pathUserExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	out, err := req.Storage.Get(ctx, fmt.Sprintf("%s/%s", "users", req.EntityID))
	if err != nil {
		return false, fmt.Errorf("existence check failed: %w", err)
	}

	return out != nil, nil
}

// setRegister adds the register to the Vault storage API
func setRegister(ctx context.Context, s logical.Storage, entityID string, registerEntry *pwmgrRegisterEntry) error {
	entry, err := logical.StorageEntryJSON(fmt.Sprintf("%s/%s", USER_SCHEMA, entityID), registerEntry)
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

	entry, err := s.Get(ctx, fmt.Sprintf("%s/%s", USER_SCHEMA, entityID))
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
	pathRegisterHelpSynopsis    = `Manages the Vault register endpoint for users to store their UUK.`
	pathRegisterHelpDescription = `
This path allows you a user to register with the pwmanager. Upon successful 
registration, the user (i.e. entityID and UUK) is added to the users schema.
`
)
