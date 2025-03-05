package secretsengine

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	bundleStoragePath = "bundle"
)

// pwmgrBundle includes the info related to
// user bundles.
type pwmgrBundle struct {
	BundleID string `json:"bundle_id"`
}

// pathBundle extends the Vault API with a `/bundle`
// endpoint for the backend. You can choose whether
// or not certain attributes should be displayed,
// required, and named. For example, password
// is marked as sensitive and will not be output
// when you read the configuration.
func pathBundle(b *pwManagerBackend) *framework.Path {
	return &framework.Path{
		Pattern: "bundle",
		Fields:  map[string]*framework.FieldSchema{},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathBundleRead,
			},
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.pathBundleWrite,
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathBundleWrite,
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.pathBundleDelete,
			},
		},
		ExistenceCheck:  b.pathBundleExistenceCheck,
		HelpSynopsis:    pathBundleHelpSynopsis,
		HelpDescription: pathBundleHelpDescription,
	}
}

// pathBundleExistenceCheck verifies if the configuration exists.
func (b *pwManagerBackend) pathBundleExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	out, err := req.Storage.Get(ctx, req.Path)
	if err != nil {
		return false, fmt.Errorf("existence check failed: %w", err)
	}

	return out != nil, nil
}

// pathBundleRead reads the configuration and outputs non-sensitive information.
func (b *pwManagerBackend) pathBundleRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	bundles, err := getBundles(ctx, req.Storage, req.EntityID)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"bundles": bundles,
		},
	}, nil
}

// pathBundleWrite updates the configuration for the backend
func (b *pwManagerBackend) pathBundleWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := getBundle(ctx, req.Storage, "TODO")
	if err != nil {
		return nil, err
	}

	createOperation := (req.Operation == logical.CreateOperation)

	if config == nil {
		if !createOperation {
			return nil, errors.New("config not found during update operation")
		}
		config = new(pwmgrBundle)
	}

	if roleID, ok := data.GetOk("bundle_id"); ok {
		config.BundleID = roleID.(string)
	} else if !ok && createOperation {
		return nil, fmt.Errorf("missing bundle_id in configuration")
	}

	entry, err := logical.StorageEntryJSON(bundleStoragePath, config)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	// reset the client so the next invocation will pick up the new configuration
	b.renew <- nil

	return nil, nil
}

// pathBundleDelete removes the configuration for the backend
func (b *pwManagerBackend) pathBundleDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete(ctx, bundleStoragePath)

	if err == nil {
		// reset the client so the next invocation will pick up the new configuration
		b.renew <- nil
	}

	return nil, err
}

func getBundle(ctx context.Context, s logical.Storage, path string) (*pwmgrBundle, error) {

	if s == nil {
		return nil, nil
	}

	entry, err := s.Get(ctx, bundleStoragePath)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	config := new(pwmgrBundle)
	if err := entry.DecodeJSON(&config); err != nil {
		return nil, fmt.Errorf("error reading root configuration: %w", err)
	}

	// return the config, we are done
	return config, nil
}

// list bundles returns the list of user bundles
func getBundles(ctx context.Context, s logical.Storage, entityId string) ([]string, error) {
	if s == nil {
		return nil, nil
	}

	userBundlePath := fmt.Sprintf("%s/%s/bundles", bundleStoragePath, entityId)
	userBundles, err := s.List(ctx, userBundlePath)
	if err != nil {
		return nil, err
	}

	if userBundles == nil {
		return nil, nil
	}

	sharedUserBundlePath := fmt.Sprintf("%s/%s/shared", bundleStoragePath, entityId)
	sharedUserBundles, err := s.List(ctx, sharedUserBundlePath)
	if err != nil {
		return nil, err
	}

	if sharedUserBundles == nil {
		return nil, nil
	}

	bundles := append(userBundles, sharedUserBundles...)

	return bundles, nil
}

// pathBundleHelpSynopsis summarizes the help text for the bundles
const pathBundleHelpSynopsis = `Configure the Pwmgr backend.`

// pathBundleHelpDescription describes the help text for the bundles
const pathBundleHelpDescription = `
The Pwmgr secret backend requires credentials for managing
JWTs issued to users working with the products API.

You must sign up with a bundle_id and secret_id and
specify the Vault AppRole address
before using this secrets backend.
`
