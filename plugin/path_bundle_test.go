package secretsengine

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
)

const (
	bundleRoleID   = "vault-plugin-testing"
	bundleSecretID = "Testing!123"
	bundleUrl      = "http://localhost:19090"
)

// testBundle mocks the creation, read, update, and delete
// of the backend configuration for Pwmgr.
func TestBundle(t *testing.T) {
	b, reqStorage := getTestBackend(t)

	t.Run("Test Bundle", func(t *testing.T) {
		entityID, _ := uuid.GenerateUUID()
		bundleID, err := testBundleCreate(t, b, reqStorage, entityID)

		assert.NoError(t, err)

		err = testBundleRead(t, b, reqStorage, entityID)

		assert.NoError(t, err)

		err = testBundleUsersAdd(t, b, reqStorage, entityID, bundleID)
		assert.NoError(t, err)
		/*
			err = testBundleUpdate(t, b, reqStorage, map[string]interface{}{
				"role_id": bundleRoleID,
				"url":     "http://pwmgr:19090",
			})

			assert.NoError(t, err)

			err = testBundleRead(t, b, reqStorage, map[string]interface{}{
				"role_id": bundleRoleID,
				"url":     "http://pwmgr:19090",
			})

			assert.NoError(t, err)

			err = testBundleDelete(t, b, reqStorage)

			assert.NoError(t, err)
		*/
	})
}

func testBundleDelete(t *testing.T, b logical.Backend, s logical.Storage, entityID string) error {
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.DeleteOperation,
		Path:      configStoragePath,
		Storage:   s,
		EntityID:  entityID,
	})

	if err != nil {
		return err
	}

	if resp != nil && resp.IsError() {
		return resp.Error()
	}
	return nil
}

func testBundleCreate(t *testing.T, b *pwManagerBackend, s logical.Storage, entityID string) (string, error) {
	ctx := context.TODO()
	data, err := b.bundleCreate(ctx, s, entityID)
	if err != nil {
		t.Errorf("error creating bundle %s", err)
	}

	if _, ok := data["bundles"]; !ok {
		t.Errorf("bundles do not exist and should")
	}

	bundles := data["bundles"].([]pwmgrBundle)

	if len(bundles) < 1 {
		t.Errorf("bundles should have more than 0 entries")
	}

	return bundles[0].ID, nil
}

func testBundleUsersAdd(t *testing.T, b *pwManagerBackend, s logical.Storage, entityID string, bundleID string) error {

	ctx := context.TODO()
	mockPolicyService := &MockPolicyService{}
	b.policyService = mockPolicyService

	var user pwManagerUserEntry
	user.UUK.PubKey = map[string]string{"test": "test"}

	b.setUserByName(ctx, s, "stephen", entityID)
	b.setUserByEntityID(ctx, s, entityID, &user)

	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.CreateOperation,
		Path:      fmt.Sprintf("bundles/%s/%s/users", entityID, bundleID),
		Storage:   s,
		EntityID:  entityID,
		Data: map[string]interface{}{
			"users": []pwmgrUser{
				{
					EntityName:   "stephen",
					IsAdmin:      true,
					Capabilities: "create,read,update,patch,delete,list",
				},
			},
		},
	})

	if err != nil || (resp != nil && resp.IsError()) {
		return resp.Error()
	}

	if p, ok := resp.Data["pubkeys"]; ok {
		for _, v := range p.(map[string]PubKey) {
			if v["test"] != "test" {
				return fmt.Errorf("should return pubkey with value `test`")
			}
		}
	}

	if mockPolicyService.CallCount != 1 {
		return fmt.Errorf("expected call count to equal 1")
	}

	sb, err := b.listSharedBundles(ctx, s, entityID)

	if err != nil {
		return err
	}

	if len(sb) != 1 {
		return fmt.Errorf("should have 1 shared bundle")
	}

	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.CreateOperation,
		Path:      fmt.Sprintf("bundles/%s/%s/users", entityID, bundleID),
		Storage:   s,
		EntityID:  entityID,
		Data: map[string]interface{}{
			"users": []pwmgrUser{},
		},
	})

	if err != nil || (resp != nil && resp.IsError()) {
		return resp.Error()
	}

	keys := []string{}
	if p, ok := resp.Data["pubkeys"]; ok {
		for k := range p.(map[string]PubKey) {
			keys = append(keys, k)
		}
	}

	if len(keys) != 0 {
		return fmt.Errorf("expected 0 pubkeys to be returned")
	}

	if mockPolicyService.CallCount != 2 {
		return fmt.Errorf("expected call count to equal 2")
	}

	return nil
}

func testBundleUpdate(t *testing.T, b logical.Backend, s logical.Storage, d map[string]interface{}) error {
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      configStoragePath,
		Data:      d,
		Storage:   s,
	})

	if err != nil {
		return err
	}

	if resp != nil && resp.IsError() {
		return resp.Error()
	}
	return nil
}

func testBundleRead(t *testing.T, b logical.Backend, s logical.Storage, entityID string) error {
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "bundles",
		Storage:   s,
		EntityID:  entityID,
	})

	if err != nil {
		return err
	}

	if resp == nil {
		return nil
	}

	if resp.IsError() {
		return resp.Error()
	}

	if v, ok := resp.Data["bundles"]; !ok {
		return fmt.Errorf("should have bundles")
	} else {
		if len(v.([]pwmgrBundle)) != 1 {
			return fmt.Errorf("should have 1 bundles")
		}
	}

	if v, ok := resp.Data["shared_bundles"]; !ok {
		return fmt.Errorf("should have bundles")
	} else {
		if len(v.([]pwmgrSharedBundle)) != 0 {
			return fmt.Errorf("should have 0 shared_bundles")
		}
	}

	return nil
}

type MockPolicyService struct {
	CallCount int
}

func (m *MockPolicyService) PutPolicy(name, rules string) error {
	m.CallCount++
	return nil
}
