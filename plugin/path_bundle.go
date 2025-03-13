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
					Callback: b.pathBundleCreate,
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.pathBundleCreate,
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.pathBundleDelete,
				},
			},
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
					Type:        framework.TypeSlice,
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

///////////////////////// bundle read /////////////////////////

// pathBundleRead returns the users owned and shared with me bundles
func (b *pwManagerBackend) pathBundleRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	bundles, err := b.listBundles(ctx, req.Storage, req.EntityID)
	if err != nil {
		return nil, err
	}

	sharedBundles, err := b.listSharedBundles(ctx, req.Storage, req.EntityID)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"bundles":        bundles,
			"shared_bundles": sharedBundles,
		},
	}, nil
}

// pathBundleCreate updates the configuration for the backend
func (b *pwManagerBackend) pathBundleCreate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	d, err := b.bundleCreate(ctx, req.Storage, req.EntityID)
	if err != nil {
		return logical.ErrorResponse("error creating bundle"), nil
	}
	return &logical.Response{
		Data: d,
	}, nil
}

// list bundles returns the list of user bundles
// TODO make this a map instead of a list
func (b *pwManagerBackend) listBundles(ctx context.Context, s logical.Storage, entityID string) ([]pwmgrBundle, error) {
	if s == nil {
		return nil, nil
	}

	userBundleListPaths := fmt.Sprintf("%s/%s/bundles/", BUNDLE_SCHEMA, entityID)
	userBundlesUUIDs, err := s.List(ctx, userBundleListPaths)
	if err != nil {
		return nil, err
	}

	if userBundlesUUIDs == nil {
		userBundlesUUIDs = []string{}
	}

	userBundles := []pwmgrBundle{}

	// Note no slash at the end
	userBundlePaths := fmt.Sprintf("%s/%s/bundles", BUNDLE_SCHEMA, entityID)
	for _, ub := range userBundlesUUIDs {
		pb, err := getBundle(ctx, s, fmt.Sprintf("%s/%s", userBundlePaths, ub))
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

func (b *pwManagerBackend) listSharedBundles(ctx context.Context, s logical.Storage, entityID string) ([]pwmgrSharedBundle, error) {
	userSharedBundlePath := fmt.Sprintf("%s/%s/sharedWithMe", BUNDLE_SCHEMA, entityID)
	sb, err := getSharedUserBundles(ctx, s, userSharedBundlePath)

	if err != nil {
		return nil, err
	}

	if sb == nil {
		sb = pwmgrSharedBundles{}
	}

	sharedBundles := []pwmgrSharedBundle{}

	for _, v := range sb {
		sharedBundles = append(sharedBundles, v)
	}

	sort.Slice(sharedBundles, func(i, j int) bool {
		return sharedBundles[i].Created < sharedBundles[j].Created
	})

	return sharedBundles, nil
}

///////////////////////// bundle create /////////////////////////

func (b *pwManagerBackend) bundleCreate(ctx context.Context, s logical.Storage, entityID string) (map[string]interface{}, error) {
	newBundleUUID, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}

	newBundleName := fmt.Sprintf("%s/%s", entityID, newBundleUUID)
	//TODO parameterize bundles in config - bundles is the kv-v2 store used to store user bundles.
	newBundleSecretPath := fmt.Sprintf("bundles/data/%s", newBundleName)

	pb := new(pwmgrBundle)
	pb.Path = newBundleSecretPath
	pb.Created = time.Now().Unix()
	pb.OwnerEntityID = entityID
	pb.ID = newBundleUUID

	// under the bundles path we store user bundles lists under /bundles/<EntityID>/bundles/<BundleUUID>
	// we need to specify a seconds bundles in the path because later we will add in shared with me path
	// for all bundles that are shared with a user.
	// Note user bundle mounts have the naming convention /bundles/<EntityID>/<UUID>.
	newBundlePath := fmt.Sprintf("%s/%s/bundles/%s", BUNDLE_SCHEMA, entityID, newBundleUUID)
	entry, err := logical.StorageEntryJSON(newBundlePath, pb)

	if err != nil {
		return nil, err
	}

	if err := s.Put(ctx, entry); err != nil {

		return nil, err
	}

	bundles, err := b.listBundles(ctx, s, entityID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"bundles": bundles,
		"path":    newBundleSecretPath,
	}, nil

}

///////////////////////// bundle delete /////////////////////////

// pathBundleDelete removes the configuration for the backend
func (b *pwManagerBackend) pathBundleDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return logical.ErrorResponse("not implemented"), nil
}

// /////////////////////// bundle create/update users /////////////////////////

func (b *pwManagerBackend) setUsersEntityID(ctx context.Context, s logical.Storage, newUsers []pwmgrUser) ([]pwmgrUser, error) {
	users := []pwmgrUser{}
	for _, nu := range newUsers {
		userEntityID, err := b.getUserEntityIDByName(ctx, s, nu.EntityName)
		if err != nil || len(userEntityID) == 0 {
			return []pwmgrUser{}, fmt.Errorf("error retrieving new users entityID")
		}

		nu.EntityID = userEntityID
		users = append(users, nu)
	}
	return users, nil
}

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

	newUsers := []pwmgrUser{}
	if usersMap, ok := d.GetOk("users"); ok {
		if err := mapstructure.Decode(usersMap, &newUsers); err != nil {
			return logical.ErrorResponse("error decoding users"), nil
		}
	} else {
		return logical.ErrorResponse("missing users"), nil
	}

	newUsers, err := b.setUsersEntityID(ctx, req.Storage, newUsers)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	usersPubKeys, err := b.getUserPubKeys(ctx, req.Storage, newUsers)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
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
	// in list of bundle users but will not show up as a shared bundle
	// for the user.
	//
	// In the future we can write to a WAL log and then
	// revert to the previous well known state on boot.
	pb.WALEntry = true
	setBundle(ctx, req.Storage, bundlePath, *pb)

	err = b.isUserBundleAdmin(req.EntityID, ownerEntityID.(string), pb.Users)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	modifiedUsers := b.getModifiedBundleUsers(*pb, newUsers)
	users := b.getUpdatedBundleUsers(*pb, newUsers)

	err = b.removeBundleUsers(ctx, req.Storage, *pb, users)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	/// modified users
	err = b.updateModifiedUsers(ctx, req.Storage, *pb, modifiedUsers)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	pb.Users = users
	pb.WALEntry = false

	err = setBundle(ctx, req.Storage, bundlePath, *pb)
	if err != nil {
		return logical.ErrorResponse("error storing bundle with new user"), err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"pubkeys": usersPubKeys,
		},
	}, nil
}

// getUserPubKeys retrieves the users public key from the KV store
func (b *pwManagerBackend) getUserPubKeys(ctx context.Context, s logical.Storage, newUsers []pwmgrUser) (map[string]PubKey, error) {
	usersPubKeys := map[string]PubKey{}

	for i := range newUsers {
		nu := &newUsers[i]
		userEntityID, err := b.getUserEntityIDByName(ctx, s, nu.EntityName)
		nu.EntityID = userEntityID

		if err != nil {
			return map[string]PubKey{}, fmt.Errorf("error retrieving new users entityID")
		}

		userUUK, err := b.getUser(ctx, s, nu.EntityID)
		if err != nil || userUUK == nil {
			return map[string]PubKey{}, fmt.Errorf("error retrieving new users public key")
		}
		usersPubKeys[nu.EntityID] = userUUK.UUK.PubKey
	}

	return usersPubKeys, nil
}

// isUserBundleAdmin checks to see if the user updating the bundle users is the owner or is an admin
// for this bundle.
func (b *pwManagerBackend) isUserBundleAdmin(entityID string, ownerEntityID string, bundleUsers []pwmgrUser) error {

	// validate user has perms to manipulate vault users
	isAdmin := false
	if entityID == ownerEntityID {
		isAdmin = true
	} else {
		for _, u := range bundleUsers {
			if u.IsAdmin && u.EntityID == entityID {
				isAdmin = true
				break
			}
		}
	}

	if !isAdmin {
		return fmt.Errorf("not authorized")
	}

	return nil
}

// getModifiedBundleUsers will return the list of users that were modified or new users. The modified
// users will have their shared bundles document updated as well as their policy.
func (b *pwManagerBackend) getModifiedBundleUsers(pb pwmgrBundle, newUsers []pwmgrUser) []pwmgrUser {
	modifiedUsers := []pwmgrUser{}

	for _, nu := range newUsers {
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

	}
	return modifiedUsers
}

// getUpdatedBundleUsers will return the updated list of bundle users.
func (b *pwManagerBackend) getUpdatedBundleUsers(pb pwmgrBundle, newUsers []pwmgrUser) []pwmgrUser {
	users := []pwmgrUser{}

	for _, nu := range newUsers {
		var userMatch *pwmgrUser
		for _, u := range pb.Users {
			if u.EntityID == nu.EntityID {
				userMatch = &u
			}
		}

		if userMatch != nil {
			nu.SharedTimestamp = userMatch.SharedTimestamp
		} else {
			nu.SharedTimestamp = time.Now().Unix()
		}
		users = append(users, nu)
	}

	return users
}

// removeBundleUsers will remove the bundle from the users shared bundle document and update the user policy
// to remove access to the bundle.
func (b *pwManagerBackend) removeBundleUsers(ctx context.Context, s logical.Storage, pb pwmgrBundle, users []pwmgrUser) error {
	usersToRemove := map[string]pwmgrUser{}
	for _, u := range pb.Users {
		usersToRemove[u.EntityID] = u
	}

	for _, u := range users {
		delete(usersToRemove, u.EntityID)
	}

	for _, u := range usersToRemove {
		userSharedBundlePath := fmt.Sprintf("%s/%s/sharedWithMe", BUNDLE_SCHEMA, u.EntityID)
		sharedBundleLock := bundleMapOfMu.Lock(userSharedBundlePath)
		{

			sbs, err := getSharedUserBundles(ctx, s, userSharedBundlePath)

			if err != nil {
				sharedBundleLock.Unlock()
				return fmt.Errorf("error reading users shared bundles")
			}

			delete(sbs, pb.ID)

			err = setSharedUserBundles(ctx, s, userSharedBundlePath, sbs)
			if err != nil {
				sharedBundleLock.Unlock()
				return err
			}

			err = b.UpdateUserPolicy(sbs, u.EntityName)
			if err != nil {
				sharedBundleLock.Unlock()
				return err
			}
		}
		sharedBundleLock.Unlock()
	}
	return nil
}

// updateModifiedUsers will add the bundles to a users shared bundles document or update the existing
// document as well as updating the user policy to access the bundle.
func (b *pwManagerBackend) updateModifiedUsers(ctx context.Context, s logical.Storage, pb pwmgrBundle, modifiedUsers []pwmgrUser) error {
	for _, mu := range modifiedUsers {
		// add bundle to users shared bundles (duplicate data is ok)
		userSharedBundlePath := fmt.Sprintf("%s/%s/sharedWithMe", BUNDLE_SCHEMA, mu.EntityID)

		sharedBundleLock := bundleMapOfMu.Lock(userSharedBundlePath)
		{
			sbs, err := getSharedUserBundles(ctx, s, userSharedBundlePath)

			if err != nil {
				sharedBundleLock.Unlock()
				return fmt.Errorf("error reading users shared bundles")
			}

			if sbs == nil {
				sbs = pwmgrSharedBundles{}
				sbs[pb.ID] = pwmgrSharedBundle{
					ID:            pb.ID,
					OwnerEntityID: pb.OwnerEntityID,
					Path:          pb.Path,
					Created:       time.Now().Unix(),
					HasAccepted:   false,
					IsAdmin:       mu.IsAdmin,
					Capabilities:  mu.Capabilities,
				}

			} else {
				sb, ok := sbs[pb.ID]

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
					sbs[pb.ID] = sb
				} else {
					sb.Capabilities = mu.Capabilities
					sb.IsAdmin = mu.IsAdmin
					sbs[pb.ID] = sb
				}
			}

			err = setSharedUserBundles(ctx, s, userSharedBundlePath, sbs)
			if err != nil {
				sharedBundleLock.Unlock()
				return err
			}

			err = b.UpdateUserPolicy(sbs, mu.EntityName)
			if err != nil {
				sharedBundleLock.Unlock()
				return fmt.Errorf("error updating user policy: %s", err)
			}
			sharedBundleLock.Unlock()
		}
	}
	return nil
}

func (b *pwManagerBackend) UpdateUserPolicy(sbs pwmgrSharedBundles, entityName string) error {
	tmpl, err := template.New("test").Parse(adminTmpl)
	if err != nil {
		return err
	}

	sharedBundles := []interface{}{}
	for _, v := range sbs {
		paths := strings.Split(v.Path, `/data/`)
		if len(paths) != 2 {
			return fmt.Errorf("bundle path is invalid: %s", v.Path)
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
		return err
	}

	// TODO find out if backend knows the mount we currently are in. if not we can
	backendMount := "pwmanager"
	policyName := fmt.Sprintf("%s/entity/%s", backendMount, entityName)
	err = b.policyService.PutPolicy(policyName, tpl.String())

	return err
}

///////////////////////// bundle kv helper /////////////////////////

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

func getSharedUserBundles(ctx context.Context, s logical.Storage, path string) (pwmgrSharedBundles, error) {

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

	sb := pwmgrSharedBundles{}
	if err := entry.DecodeJSON(&sb); err != nil {
		return nil, fmt.Errorf("error reading shared bundles %w", err)
	}

	// return the bundle path, we are done
	return sb, nil
}

func setSharedUserBundles(ctx context.Context, s logical.Storage, path string, sbp pwmgrSharedBundles) error {
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

// bundles/<id>/users

// pathBundleExistenceCheck verifies if the configuration exists.
func (b *pwManagerBackend) pathBundleExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	out, err := req.Storage.Get(ctx, req.Path)
	if err != nil {
		return false, fmt.Errorf("existence check failed: %w", err)
	}

	return out != nil, nil
}

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

// pathBundleHelpSynopsis summarizes the help text for the bundles
const pathBundleHelpSynopsis = `bundles endpoints allow users to create and share bundles.`

// pathBundleHelpDescription describes the help text for the bundles
const pathBundleHelpDescription = `bundles endpoints allow users to create and share bundles.`
