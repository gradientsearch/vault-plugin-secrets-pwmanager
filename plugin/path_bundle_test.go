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
		err := testBundleCreate(t, b, reqStorage)

		assert.NoError(t, err)

		/*
			err = testBundleRead(t, b, reqStorage, map[string]interface{}{
				"role_id": bundleRoleID,
				"url":     bundleUrl,
			})

			assert.NoError(t, err)

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

func testBundleDelete(t *testing.T, b logical.Backend, s logical.Storage) error {
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.DeleteOperation,
		Path:      configStoragePath,
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

func testBundleCreate(t *testing.T, b *pwManagerBackend, s logical.Storage) error {
	entityID, _ := uuid.GenerateUUID()
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

func testBundleRead(t *testing.T, b logical.Backend, s logical.Storage, expected map[string]interface{}) error {
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ReadOperation,
		Path:      configStoragePath,
		Storage:   s,
	})

	if err != nil {
		return err
	}

	if resp == nil && expected == nil {
		return nil
	}

	if resp.IsError() {
		return resp.Error()
	}

	if len(expected) != len(resp.Data) {
		return fmt.Errorf("read data mismatch (expected %d values, got %d)", len(expected), len(resp.Data))
	}

	for k, expectedV := range expected {
		actualV, ok := resp.Data[k]

		if !ok {
			return fmt.Errorf(`expected data["%s"] = %v but was not included in read output"`, k, expectedV)
		} else if expectedV != actualV {
			return fmt.Errorf(`expected data["%s"] = %v, instead got %v"`, k, expectedV, actualV)
		}
	}

	return nil
}
