package secretsengine

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"sort"
	"strings"
	"time"

	mapstructure "github.com/go-viper/mapstructure/v2"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	BUNDLE_SCHEMA = "bundles"
)

var bundleMapOfMu = NewMapOfMu()

type pwmgrUser struct {
	EntityID        string `json:"entity_id" mapstructure:"entity_id"`
	EntityName      string `json:"entity_name" mapstructure:"entity_name"`
	IsAdmin         bool   `json:"is_admin" mapstructure:"is_admin"`
	SharedTimestamp int64  `json:"shared_timestamp" mapstructure:"shared_timestamp"`
	// comma separated string of capabilities
	Capabilities string `json:"capabilities" mapstructure:"capabilities"`
}

type pwmgrUsers struct {
	Users []pwmgrUser `json:"users" mapstructure:"users"`
}

// pwmgrBundle includes the info related to
// user bundles.
type pwmgrBundle struct {
	ID   string `json:"id"`
	Path string `json:"path"`
	// for order. we could do alphabetically
	Created int64 `json:"created"`

	OwnerEntityID string `json:"owner_entity_id"`

	Users []pwmgrUser `json:"users"`

	WALEntry bool `json:"wal_entry"`
}

type pwmgrSharedBundle struct {
	ID   string `json:"id"`
	Path string `json:"path"`
	// for order. we could do alphabetically
	Created       int64  `json:"created"`
	OwnerEntityID string `json:"owner_entity_id"`
	HasAccepted   bool   `json:"has_accepted"`
	IsAdmin       bool   `json:"is_admin"`
	// comma separated string of capabilities
	Capabilities string `json:"capabilities" mapstructure:"capabilities"`
}

type pwmgrSharedBundles map[string]pwmgrSharedBundle

// admin template
var adminTmpl = `
{{range $index, $bundle := . }}
path "bundles/data/{{$bundle.Path}}/*" {
    capabilities = [{{$first := true}}{{range $bundle.Capabilities}}{{if $first}}{{$first = false}}{{else}}, {{end}}"{{.}}"{{end}} ]
}

path "bundles/metadata/{{$bundle.Path}}/*" {
    capabilities = [ {{$first := true}}{{range $bundle.Capabilities}}{{if $first}}{{$first = false}}{{else}}, {{end}}"{{.}}"{{end}} ]
}
{{end}}`

// pathBundle extends the Vault API with a `/bundle`
// endpoint for the backend. You can choose whether
// or not certain attributes should be displayed,
// required, and named. For example, password
// is marked as sensitive and will not be output
// when you read the configuration.
func pathBundle(b *pwManagerBackend) []*framework.Path {
	return []*framework.Path{
		{
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
		},
		{
			Pattern: fmt.Sprintf("bundles/%s/%s/users", uuidRegex("owner_entity_id"), uuidRegex("bundle_id")),
			Fields: map[string]*framework.FieldSchema{
				"owner_entity_id": {
					Type:        framework.TypeLowerCaseString,
					Description: "entity id of the bundle owner",
					Required:    true,
				},
				"bundle_id": {
					Type:        framework.TypeLowerCaseString,
					Description: "uuid of the bundle",
					Required:    true,
				},
				"users": {
					Type:        framework.TypeMap,
					Description: "users for this bundle",
					Required:    false,
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.pathBundleUsersWrite,
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.pathBundleUsersWrite,
				},
			},
			ExistenceCheck:  b.pathBundleExistenceCheck,
			HelpSynopsis:    pathBundleHelpSynopsis,
			HelpDescription: pathBundleHelpDescription,
		},
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

	userSharedBundlePath := fmt.Sprintf("%s/%s/sharedWithMe", BUNDLE_SCHEMA, req.EntityID)
	sbp, err := getSharedUserBundles(ctx, req.Storage, userSharedBundlePath)

	if err != nil {
		return nil, err
	}

	if sbp == nil {
		sbp = &pwmgrSharedBundles{}
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"bundles":        bundles,
			"shared_bundles": *sbp,
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
	pb.OwnerEntityID = req.EntityID
	pb.ID = newBundleUUID

	// under the bundles path we store user bundles lists under /bundles/<EntityID>/bundles/<BundleUUID>
	// we need to specify a seconds bundles in the path because later we will add in shared with me path
	// for all bundles that are shared with a user.
	// Note user bundle mounts have the naming convention /bundles/<EntityID>/<UUID>.
	newBundlePath := fmt.Sprintf("%s/%s/bundles/%s", BUNDLE_SCHEMA, req.EntityID, newBundleUUID)
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
	err := req.Storage.Delete(ctx, BUNDLE_SCHEMA)

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
		// TODO this needs to return an error!!!!
		return nil, nil
	}

	pb := new(pwmgrBundle)
	if err := entry.DecodeJSON(&pb); err != nil {
		return nil, fmt.Errorf("error reading root configuration: %w", err)
	}

	// return the bundle path, we are done
	return pb, nil
}

func setBundle(ctx context.Context, s logical.Storage, path string, pb pwmgrBundle) error {

	if s == nil {
		return nil
	}

	// TODO move this to a putBundle method
	// Only after successfully completing previous step
	// we will update the users. Even if the previous
	// steps were partially completed, for example
	newBundleEntry, err := logical.StorageEntryJSON(path, pb)

	if err != nil {
		return err
	}

	if err := s.Put(ctx, newBundleEntry); err != nil {
		return err
	}

	// return the bundle path, we are done
	return nil
}

func getSharedUserBundles(ctx context.Context, s logical.Storage, path string) (*pwmgrSharedBundles, error) {

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

	sb := new(pwmgrSharedBundles)
	if err := entry.DecodeJSON(&sb); err != nil {
		return nil, fmt.Errorf("error reading shared bundles %w", err)
	}

	// return the bundle path, we are done
	return sb, nil
}

func putSharedUserBundles(ctx context.Context, s logical.Storage, path string, sbp pwmgrSharedBundles) error {
	//	TODO versioned writes
	if s == nil {
		return fmt.Errorf("error storage cannot be nil")
	}

	entry, err := logical.StorageEntryJSON(path, sbp)

	if err != nil {
		return err
	}

	if err := s.Put(ctx, entry); err != nil {
		return err
	}

	return nil
}

// list bundles returns the list of user bundles
func (b *pwManagerBackend) listBundles(ctx context.Context, s logical.Storage, entityID string) ([]pwmgrBundle, error) {
	if s == nil {
		return nil, nil
	}

	userBundlePaths := fmt.Sprintf("%s/%s/bundles/", BUNDLE_SCHEMA, entityID)
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

	sort.Slice(userBundles, func(i, j int) bool {
		return userBundles[i].Created < userBundles[j].Created
	})

	return userBundles, nil
}

// bundles/<id>/users

// pathBundleWrite updates the configuration for the backend
func (b *pwManagerBackend) pathBundleUsersWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	ownerEntityID, ok := d.GetOk("owner_entity_id")
	if !ok {
		return logical.ErrorResponse("missing owner entity id "), nil
	}

	bundleID, ok := d.GetOk("bundle_id")
	if !ok {
		return logical.ErrorResponse("missing bundle id "), nil
	}

	newUsers := pwmgrUsers{}

	// TODO add WAL Write

	if users, ok := d.GetOk("users"); ok {
		if err := mapstructure.Decode(users, &newUsers); err != nil {
			return logical.ErrorResponse("error decoding users"), nil
		}
	} else {
		return logical.ErrorResponse("missing users"), nil
	}

	usersPubKeys := map[string]map[string]string{}

	for i := range newUsers.Users {
		nu := &newUsers.Users[i]
		// make sure entity name and entity id have integrity
		userEntityID, err := b.getUserEntityIDByName(ctx, req.Storage, nu.EntityName)
		nu.EntityID = userEntityID

		if err != nil {
			return logical.ErrorResponse("error retrieving new users entityID"), nil
		}

		// grab the public key of user
		userUUK, err := b.getUser(ctx, req.Storage, nu.EntityID)

		if err != nil || userUUK == nil {
			return logical.ErrorResponse("error retrieving new users public key"), nil
		}
		usersPubKeys[nu.EntityID] = userUUK.UUK.PubKey
	}

	//TODO parameterize bundles in config - bundles is the kv-v2 store used to store user bundles.

	// under the bundles path we store user bundles lists under /bundles/<EntityID>/bundles/<BundleUUID>
	// we need to specify a seconds bundles in the path because later we will add in shared with me path
	// for all bundles that are shared with a user.
	// Note user bundle mounts have the naming convention /bundles/<EntityID>/<UUID>.
	bundlePath := fmt.Sprintf("%s/%s/bundles/%s", BUNDLE_SCHEMA, ownerEntityID, bundleID)

	bundleLock := bundleMapOfMu.Lock(bundlePath)
	defer bundleLock.Unlock()

	pb, err := getBundle(ctx, req.Storage, bundlePath)
	if err != nil || pb == nil {
		return logical.ErrorResponse("bundle not found"), nil
	}

	// The current edge cases are around data integrity.
	// If a current bundle user was given less/more access
	// and those updates were written for the user but the
	// server crashed before the bundle users were written
	// to disk.
	//
	// Other scenarios exists but have less security concerns. e.g
	// 1. If a user was added but the server crashed the users pubkey
	// was not returned to encrypt the bundle key.
	// 2. If a user was deleted the user will still show up
	// in list of bundle users.
	//
	// In the future we can write to a WAL log and then
	// revert to the previous well known state.
	pb.WALEntry = true
	setBundle(ctx, req.Storage, bundlePath, *pb)

	// validate user has perms to manipulate vault users
	isAdmin := false
	if req.EntityID == ownerEntityID {
		isAdmin = true
	} else {
		for _, u := range pb.Users {
			if u.IsAdmin && u.EntityID == req.EntityID {
				isAdmin = true
				break
			}
		}
	}

	if !isAdmin {
		return logical.ErrorResponse("not authorized"), nil
	}

	modifiedUsers := []pwmgrUser{}
	users := []pwmgrUser{}

	for _, nu := range newUsers.Users {
		// TODO think through workflow .
		// how will a user be modified?
		// we go ahead and update the users policy to allow
		// if users declines invitation we can delete the bundle
		// and recreate the policy for that user
		var userMatch *pwmgrUser

		for _, u := range pb.Users {
			if u.EntityID == nu.EntityID {
				userMatch = &u
			}
		}

		if userMatch != nil {

			// copy over
			nu.SharedTimestamp = userMatch.SharedTimestamp

			// check if modified
			modified := false

			// If pervious request did not finish we'll update all user
			if pb.WALEntry {
				modified = true
			}

			if nu.Capabilities != userMatch.Capabilities {
				modified = true
			}

			if nu.IsAdmin != userMatch.IsAdmin {
				modified = true
			}

			if nu.EntityName != userMatch.EntityName {
				modified = true
			}

			if modified {
				modifiedUsers = append(modifiedUsers, nu)
			}

		} else {
			// grab the entity of the new user
			nu.SharedTimestamp = time.Now().Unix()
			modifiedUsers = append(modifiedUsers, nu)
		}

		users = append(users, nu)

	}

	/////////// delete users ////////////////
	deletedUsers := map[string]bool{}
	for _, u := range pb.Users {
		deletedUsers[u.EntityID] = true
	}

	for _, u := range users {
		deletedUsers[u.EntityID] = false
	}

	for k, v := range deletedUsers {
		if v {
			// add bundle to users shared bundles (duplicate data is ok)

			userSharedBundlePath := fmt.Sprintf("%s/%s/sharedWithMe", BUNDLE_SCHEMA, k)
			sharedBundleLock := bundleMapOfMu.Lock(userSharedBundlePath)
			{
				var sbp *pwmgrSharedBundles
				sbp, err := getSharedUserBundles(ctx, req.Storage, userSharedBundlePath)

				if err != nil {
					sharedBundleLock.Unlock()
					return logical.ErrorResponse("error reading users shared bundles"), err
				}

				sbpm := *sbp
				delete(sbpm, bundleID.(string))

				err = putSharedUserBundles(ctx, req.Storage, userSharedBundlePath, *sbp)

				if err != nil {
					sharedBundleLock.Unlock()
					return nil, err
				}

				sharedBundleLock.Unlock()
				// TODO update user policy
				// template new policy
				// post to vault

				// TODO remove user key from bundle
				//

			}

		}
	}

	for _, mu := range modifiedUsers {
		// add bundle to users shared bundles (duplicate data is ok)
		userSharedBundlePath := fmt.Sprintf("%s/%s/sharedWithMe", BUNDLE_SCHEMA, mu.EntityID)

		sharedBundleLock := bundleMapOfMu.Lock(userSharedBundlePath)
		{
			var sbp *pwmgrSharedBundles
			sbp, err := getSharedUserBundles(ctx, req.Storage, userSharedBundlePath)

			if err != nil {
				sharedBundleLock.Unlock()
				return logical.ErrorResponse("error reading users shared bundles"), err
			}

			if sbp == nil {
				sbm := make(pwmgrSharedBundles)

				sbm[bundleID.(string)] = pwmgrSharedBundle{
					ID:            pb.ID,
					OwnerEntityID: pb.OwnerEntityID,
					Path:          pb.Path,
					Created:       time.Now().Unix(),
					HasAccepted:   false,
					IsAdmin:       mu.IsAdmin,
					Capabilities:  mu.Capabilities,
				}

				sbp = &sbm
			} else {
				sbsm := *sbp
				sb, ok := sbsm[bundleID.(string)]

				if !ok {
					sb = pwmgrSharedBundle{
						ID:            pb.ID,
						OwnerEntityID: pb.OwnerEntityID,
						Path:          pb.Path,
						Created:       time.Now().Unix(),
						HasAccepted:   false,
						IsAdmin:       mu.IsAdmin,
						Capabilities:  mu.Capabilities,
					}
					sbsm[bundleID.(string)] = sb

				} else {
					sb.Capabilities = mu.Capabilities
					sb.IsAdmin = mu.IsAdmin
					sb.Created = time.Now().Unix()
					sbsm[bundleID.(string)] = sb
				}
			}

			err = putSharedUserBundles(ctx, req.Storage, userSharedBundlePath, *sbp)

			if err != nil {
				sharedBundleLock.Unlock()
				return nil, err
			}

			//TODO move this out of loop
			tmpl, err := template.New("test").Parse(adminTmpl)
			if err != nil {
				sharedBundleLock.Unlock()
				return nil, err
			}

			sharedBundles := []interface{}{}
			for _, v := range *sbp {
				// update users policy
				paths := strings.Split(v.Path, `/data/`)
				if len(paths) != 2 {
					sharedBundleLock.Unlock()
					return nil, fmt.Errorf("bundle path is invalid: %s", v.Path)
				}

				caps := strings.Split(v.Capabilities, `,`)

				b := struct {
					Path         string
					Capabilities []string
				}{Path: paths[1], Capabilities: caps}

				sharedBundles = append(sharedBundles, b)

			}

			var tpl bytes.Buffer
			err = tmpl.Execute(&tpl, sharedBundles)
			if err != nil {
				sharedBundleLock.Unlock()
				return nil, err
			}

			// TODO find out if backend knows the mount we currently are in. if not we can
			// add it to the config
			policyName := fmt.Sprintf("%s/entity/%s", "pwmanager", mu.EntityName)
			err = b.c.c.Sys().PutPolicy(policyName, tpl.String())

			if err != nil {
				sharedBundleLock.Unlock()
				return logical.ErrorResponse(fmt.Sprintf("error updating user policy: %s", err)), nil
			}
			sharedBundleLock.Unlock()
		}

		//// pull in all user shared bundles
		// error handling to undo if not successful
	}

	// TODO move this to a putBundle method
	// Only after successfully completing previous step
	// we will update the users. Even if the previous
	// steps were partially completed, for example
	pb.Users = users
	pb.WALEntry = false

	err = setBundle(ctx, req.Storage, bundlePath, *pb)
	if err != nil {
		return logical.ErrorResponse("error storing bundle with new user"), err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"pubkey": usersPubKeys,
		},
	}, nil

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
