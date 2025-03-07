package secretsengine

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	bundleStoragePath = "bundle"
)

// pwmgrBundle includes the info related to
// user bundles.
type pwmgrBundle struct {
	Path string `json:"path"`
	// for order. we could do alphabetically
	Created int64 `json:"created"`
}

// pathBundle extends the Vault API with a `/bundle`
// endpoint for the backend. You can choose whether
// or not certain attributes should be displayed,
// required, and named. For example, password
// is marked as sensitive and will not be output
// when you read the configuration.
func pathBundle(b *pwManagerBackend) *framework.Path {
	return &framework.Path{
		Pattern: "bundles",
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
	bundles, err := b.listBundles(ctx, req.Storage, req.EntityID)
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

	newBundleUUID, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}

	newBundleName := fmt.Sprintf("%s/%s", req.EntityID, newBundleUUID)
	//TODO parameterize bundles in config - bundles is the kv-v2 store used to store user bundles.
	newBundleSecretPath := fmt.Sprintf("bundles/data/%s", newBundleName)

	// There are some pros and cons to creating new secret mounts for new user vaults.
	// A pro with distinct kv-v2 stores data is separate and metadata per mount is easily
	// available. A major con is the limit of secret mounts in vault is set to 14000
	// https://developer.hashicorp.com/vault/docs/internals/limits#mount-point-limits
	// even tho vault kv-v2 stores are multiplexed there may still be overhead for the additional
	// mounts. Right now, i've pivoted to just use a single kv-v2 store and use policies to control
	// access. In the future, we can always pivot back if necessary.

	// mi := api.MountInput{
	// 	Type: "kv-v2",
	// }
	// err = b.c.c.Sys().Mount(newBundleMountPath, &mi)

	// if err != nil {
	// 	return nil, err
	// }

	pb := new(pwmgrBundle)
	pb.Path = newBundleSecretPath
	pb.Created = time.Now().Unix()

	// under the bundles path we store user bundles lists under /bundles/<EntityID>/bundles/<BundleUUID>
	// we need to specify a seconds bundles in the path because later we will add in shared with me path
	// for all bundles that are shared with a user.
	// Note user bundle mounts have the naming convention /bundles/<EntityID>/<UUID>.
	newBundlePath := fmt.Sprintf("%s/%s/bundles/%s", bundleStoragePath, req.EntityID, newBundleUUID)
	entry, err := logical.StorageEntryJSON(newBundlePath, pb)

	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {

		return nil, err
	}

	bundles, err := b.listBundles(ctx, req.Storage, req.EntityID)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"bundles": bundles,
			"path":    newBundleSecretPath,
		},
	}, nil

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

	entry, err := s.Get(ctx, path)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	pb := new(pwmgrBundle)
	if err := entry.DecodeJSON(&pb); err != nil {
		return nil, fmt.Errorf("error reading root configuration: %w", err)
	}

	// return the bundle path, we are done
	return pb, nil
}

// list bundles returns the list of user bundles
func (b *pwManagerBackend) listBundles(ctx context.Context, s logical.Storage, entityID string) ([]pwmgrBundle, error) {
	if s == nil {
		return nil, nil
	}

	userBundlePaths := fmt.Sprintf("%s/%s/bundles/", bundleStoragePath, entityID)
	userBundlesUUIDs, err := s.List(ctx, userBundlePaths)
	if err != nil {
		return nil, err
	}

	if userBundlesUUIDs == nil {
		userBundlesUUIDs = []string{}
	}

	userBundles := []pwmgrBundle{}
	for _, ub := range userBundlesUUIDs {
		// TODO list needs a / at the end. Clean this up
		pb, err := getBundle(ctx, s, fmt.Sprintf("%s%s", userBundlePaths, ub))
		if err != nil {
			return nil, err
		}
		if pb == nil {
			b.logger.Warn(fmt.Sprintf("user bundle is nil: path: %s%s", userBundlePaths, ub))
			continue
		}
		userBundles = append(userBundles, *pb)
	}

	// // to accept shared bundle path as well.
	// sharedWithUserBundlePath := fmt.Sprintf("%s/%s/shared-with-user", bundleStoragePath, entityID)
	// sharedWithUserBundles, err := s.List(ctx, sharedWithUserBundlePath)
	// if err != nil {
	// 	return nil, err
	// }

	// if sharedWithUserBundles == nil {
	// 	sharedWithUserBundles = []string{}
	// }

	sort.Slice(userBundles, func(i, j int) bool {
		return userBundles[i].Created < userBundles[j].Created
	})

	return userBundles, nil
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
