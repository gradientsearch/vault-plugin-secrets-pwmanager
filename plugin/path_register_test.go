package secretsengine

import (
	"context"
	"crypto/rand"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

const (
	registerName = "testpwmgr"
	registerID   = "75e46130-bf36-11ef-8781-abab9f01266a"
)

// TestUserRegister uses a mock backend to check
// register create, read, update, and delete.
func TestUserRegister(t *testing.T) {
	b, s := getTestBackend(t)

	t.Run("Create User Register - pass", func(t *testing.T) {
		resp, err := testTokenRegisterCreate(t, b, s, registerName, map[string]interface{}{
			"username": registerID,
			"ttl":      testTTL,
			"max_ttl":  testMaxTTL,
		})

		require.Nil(t, err)
		require.Nil(t, resp.Error())
		require.Nil(t, resp)
	})
	return
	t.Run("List All Registers", func(t *testing.T) {
		for i := 1; i <= 10; i++ {
			_, err := testTokenRegisterCreate(t, b, s,
				registerName+strconv.Itoa(i),
				map[string]interface{}{
					"username": registerID,
					"ttl":      testTTL,
					"max_ttl":  testMaxTTL,
				})
			require.NoError(t, err)
		}

		resp, err := testTokenRegisterList(t, b, s)
		require.NoError(t, err)
		require.Len(t, resp.Data["keys"].([]string), 10)
	})

	t.Run("Read User Register", func(t *testing.T) {
		resp, err := testTokenRegisterRead(t, b, s)

		require.Nil(t, err)
		require.Nil(t, resp.Error())
		require.NotNil(t, resp)
		require.Equal(t, resp.Data["username"], registerID)
	})
	t.Run("Update User Register", func(t *testing.T) {
		resp, err := testTokenRegisterUpdate(t, b, s, map[string]interface{}{
			"ttl":     "1m",
			"max_ttl": "5h",
		})

		require.Nil(t, err)
		require.Nil(t, resp.Error())
		require.Nil(t, resp)
	})

	t.Run("Re-read User Register", func(t *testing.T) {
		resp, err := testTokenRegisterRead(t, b, s)

		require.Nil(t, err)
		require.Nil(t, resp.Error())
		require.NotNil(t, resp)
		require.Equal(t, resp.Data["username"], registerID)
	})

	t.Run("Delete User Register", func(t *testing.T) {
		_, err := testTokenRegisterDelete(t, b, s)

		require.NoError(t, err)
	})
}

// Utility function to create a register while, returning any response (including errors)
func testTokenRegisterCreate(t *testing.T, b *pwmgrBackend, s logical.Storage, name string, d map[string]interface{}) (*logical.Response, error) {
	t.Helper()

	uuk := UUK{}
	password := "gophers"
	secretKey := make([]byte, 32)

	if _, err := rand.Read(secretKey); err != nil {
		t.Fatalf("error creating secret key: %s; ", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatalf("error generating UUID")
	}
	uuk.Build([]byte(password), []byte("pwmanager"), secretKey, []byte(id))
	m, err := uuk.Map()
	if err != nil {
		t.Fatal(err)
	}

	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "register",
		Data:      m,
		Storage:   s,
		EntityID:  id,
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Utility function to update a register while, returning any response (including errors)
func testTokenRegisterUpdate(t *testing.T, b *pwmgrBackend, s logical.Storage, d map[string]interface{}) (*logical.Response, error) {
	t.Helper()
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "register/" + registerName,
		Data:      d,
		Storage:   s,
	})

	if err != nil {
		return nil, err
	}

	if resp != nil && resp.IsError() {
		t.Fatal(resp.Error())
	}
	return resp, nil
}

// Utility function to read a register and return any errors
func testTokenRegisterRead(t *testing.T, b *pwmgrBackend, s logical.Storage) (*logical.Response, error) {
	t.Helper()
	return b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "register/" + registerName,
		Storage:   s,
	})
}

// Utility function to list registers and return any errors
func testTokenRegisterList(t *testing.T, b *pwmgrBackend, s logical.Storage) (*logical.Response, error) {
	t.Helper()
	return b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ListOperation,
		Path:      "register/",
		Storage:   s,
	})
}

// Utility function to delete a register and return any errors
func testTokenRegisterDelete(t *testing.T, b *pwmgrBackend, s logical.Storage) (*logical.Response, error) {
	t.Helper()
	return b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.DeleteOperation,
		Path:      "register/" + registerName,
		Storage:   s,
	})
}
func TestRegisterUser(t *testing.T) {

	th, err := NewTestHarness(t, "TestRegisterUser", false)
	if err != nil {
		t.Fatalf("failed to start Vault container")
	}
	defer th.Teardown()

	mount := "pwmanager"
	th.WithPwManagerMount()

	userPolicy, err := os.ReadFile("policies/pwmanager_user_default.hcl")
	if err != nil {
		th.Testing.Fatalf("error reading user policy: %s", err)
	}

	policies := map[string]string{
		"pwmanager/user/default": string(userPolicy),
	}

	th.WithPolicies(policies)
	users := th.WithUserpassAuth("pwmanager", []string{"stephen", "frank", "bob", "alice"})
	stephen := users["stephen"]

	stephen.WithUUK(th)

	if err := stephen.Client.PwManager().Register(mount, stephen.UUK); err != nil {
		th.Testing.Fatalf("error registering user: %s", err)
	}
}
