package secretsengine

import (
	"context"
	"fmt"
	"time"

	mapstructure "github.com/go-viper/mapstructure/v2"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	USER_SCHEMA string = "users"
)

// pwmgrUserEntry defines the data required
// for a Vault user to access and call the Pwmgr
// token endpoints
type pwmgrUserEntry struct {
	Username string        `json:"username"`
	UserID   int           `json:"user_id"`
	EntityID string        `json:"entity_id"`
	Token    string        `json:"token"`
	TokenID  string        `json:"token_id"`
	TTL      time.Duration `json:"ttl"`
	MaxTTL   time.Duration `json:"max_ttl"`
	UUK      pwmgrUUKEntry `json:"uuk"`
}

// pwmgrUUKEntry defines the data required
// for a Vault register to access and call the Pwmgr
// token endpoints
type pwmgrUUKEntry struct {
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
		{
			Pattern: fmt.Sprintf("%s/", USER_SCHEMA) + uuidRegex("entity_id"),
			Fields: map[string]*framework.FieldSchema{
				"entity_id": {
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
			Pattern: fmt.Sprintf("%s/?$", USER_SCHEMA),
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
	entries, err := req.Storage.List(ctx, fmt.Sprintf("%s/", USER_SCHEMA))
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(entries), nil
}

// pathUsersRead makes a request to Vault storage to read a user and return response data
func (b *pwmgrBackend) pathUsersRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entry, err := b.getUser(ctx, req.Storage, d.Get("entity_id").(string))
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
	name, ok := d.GetOk("entity_id")
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
	err := req.Storage.Delete(ctx, fmt.Sprintf("%s/%s", USER_SCHEMA, d.Get("entity_id").(string)))
	if err != nil {
		return nil, fmt.Errorf("error deleting pwmgr user: %w", err)
	}

	return nil, nil
}

// pathRegistersWrite makes a request to Vault storage to register a users UUK.
func (b *pwmgrBackend) pathRegistersWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	registerUser, err := b.getUser(ctx, req.Storage, req.EntityID)
	if err != nil {
		return nil, err
	}

	if registerUser != nil {
		return logical.ErrorResponse("user already registered"), nil
	}

	registerEntry := pwmgrUUKEntry{}

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

	userEntry := pwmgrUserEntry{}
	userEntry.EntityID = req.EntityID
	userEntry.UUK = registerEntry

	if err := setUser(ctx, req.Storage, req.EntityID, &userEntry); err != nil {
		return nil, err
	}

	return nil, nil
}

// pathConfigExistenceCheck verifies if the configuration exists.
func (b *pwmgrBackend) pathUserExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	out, err := req.Storage.Get(ctx, fmt.Sprintf("%s/%s", USER_SCHEMA, req.EntityID))
	if err != nil {
		return false, fmt.Errorf("existence check failed: %w", err)
	}

	return out != nil, nil
}

// setRegister adds the register to the Vault storage API
func setUser(ctx context.Context, s logical.Storage, entityID string, registerEntry *pwmgrUserEntry) error {
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

// getUser gets the register from the Vault storage API
func (b *pwmgrBackend) getUser(ctx context.Context, s logical.Storage, entityID string) (*pwmgrUserEntry, error) {
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

	var register pwmgrUserEntry

	if err := entry.DecodeJSON(&register); err != nil {
		return nil, err
	}
	return &register, nil
}

const (
	pathUserHelpSynopsis    = `Manages the Vault user for generating Pwmgr tokens.`
	pathUserHelpDescription = `
This path allows you to read and write users used to generate Pwmgr tokens.
You can configure a user to manage a user's token by setting the username field.
`

	pathUserListHelpSynopsis    = `List the existing users in Pwmgr backend`
	pathUserListHelpDescription = `Users will be listed by the user name.`

	pathRegisterHelpSynopsis    = `Manages the Vault register endpoint for users to store their UUK.`
	pathRegisterHelpDescription = `
This path allows you a user to register with the pwmanager. Upon successful 
registration, the user (i.e. entityID and UUK) is added to the users schema.
`
)
