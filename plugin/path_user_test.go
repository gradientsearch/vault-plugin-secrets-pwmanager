package secretsengine

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

const (
	USER_ENTITY_ID = "928e91c7-db18-9673-4342-6f731c7f561a"
)

// TestUserUser uses a mock backend to check
// user create, read, update, and delete.
func TestUserUser(t *testing.T) {
	b, s := getTestBackend(t)

	t.Run("List All Users", func(t *testing.T) {
		for i := 1; i <= 10; i++ {
			id, err := uuid.GenerateUUID()
			if err != nil {
				t.Fatalf("failed creating uuid for List All Users: %s", err)
			}

			_, err = testTokenRegisterCreate(t, b, s, id)
			require.NoError(t, err)
		}

		resp, err := testTokenUserList(t, b, s)
		require.NoError(t, err)
		require.Len(t, resp.Data["keys"].([]string), 10)
	})

	t.Run("Create User User - pass", func(t *testing.T) {
		resp, err := testTokenRegisterCreate(t, b, s, USER_ENTITY_ID)

		require.Nil(t, err)
		require.Nil(t, resp.Error())
		require.Nil(t, resp)
	})

	t.Run("Read User User", func(t *testing.T) {
		resp, err := testTokenUserRead(t, b, s)

		require.Nil(t, err)
		require.Nil(t, resp.Error())
		require.NotNil(t, resp)
		require.Equal(t, resp.Data["entity_id"], USER_ENTITY_ID)
	})
	t.Run("Update User User", func(t *testing.T) {
		resp, err := testTokenUserUpdate(t, b, s, USER_ENTITY_ID, map[string]interface{}{

			"entity_id": USER_ENTITY_ID,
			"uuk":       pwmgrUUKEntry{},
		})

		require.Nil(t, err)
		require.Nil(t, resp.Error())
		require.Nil(t, resp)

		otherUserEntityID, err := uuid.GenerateUUID()
		if err != nil {
			t.Fatalf("failed generating uuid: %s", err)
		}
		resp, err = testTokenUserUpdate(t, b, s, otherUserEntityID, map[string]interface{}{
			"entity_id": otherUserEntityID,
			"uuk":       pwmgrUUKEntry{},
		})

		if err != nil {
			t.Fatalf("error making update request %s", err)
		}

		require.NotNil(t, resp)
		require.NotNil(t, resp.IsError())
		require.Equal(t, resp.Error().Error(), "users can only modify their own user information", "error response changed")
	})

	t.Run("Re-read User User", func(t *testing.T) {
		resp, err := testTokenUserRead(t, b, s)

		require.Nil(t, err)
		require.Nil(t, resp.Error())
		require.NotNil(t, resp)
		require.Equal(t, resp.Data["entity_id"], USER_ENTITY_ID)
	})

	t.Run("Delete User User", func(t *testing.T) {
		_, err := testTokenUserDelete(t, b, s)

		require.NoError(t, err)
	})
}

// Utility function to update a user while, returning any response (including errors)
func testTokenUserUpdate(t *testing.T, b *pwmgrBackend, s logical.Storage, entityIDToUpdate string, d map[string]interface{}) (*logical.Response, error) {
	t.Helper()
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      fmt.Sprintf("%s/%s", USER_SCHEMA, entityIDToUpdate),
		Data:      d,
		Storage:   s,
		EntityID:  USER_ENTITY_ID, // !important users can only modify their own user information
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Utility function to read a user and return any errors
func testTokenUserRead(t *testing.T, b *pwmgrBackend, s logical.Storage) (*logical.Response, error) {
	t.Helper()
	return b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ReadOperation,
		Path:      fmt.Sprintf("%s/%s", USER_SCHEMA, USER_ENTITY_ID),
		Storage:   s,
	})
}

// Utility function to list users and return any errors
func testTokenUserList(t *testing.T, b *pwmgrBackend, s logical.Storage) (*logical.Response, error) {
	t.Helper()
	return b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ListOperation,
		Path:      fmt.Sprintf("%s/", USER_SCHEMA),
		Storage:   s,
	})
}

// Utility function to delete a user and return any errors
func testTokenUserDelete(t *testing.T, b *pwmgrBackend, s logical.Storage) (*logical.Response, error) {
	t.Helper()
	return b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.DeleteOperation,
		Path:      fmt.Sprintf("%s/%s", USER_SCHEMA, USER_ENTITY_ID),
		Storage:   s,
	})
}

// Utility function to create a register while, returning any response (including errors)
func testTokenRegisterCreate(t *testing.T, b *pwmgrBackend, s logical.Storage, entityID string) (*logical.Response, error) {
	t.Helper()

	uuk := UUK{}
	password := "gophers"
	secretKey := make([]byte, 32)

	if _, err := rand.Read(secretKey); err != nil {
		t.Fatalf("error creating secret key: %s; ", err)
	}

	uuk.Build([]byte(password), []byte("pwmanager"), secretKey, []byte(entityID))
	m, err := uuk.Map()
	if err != nil {
		t.Fatal(err)
	}

	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "register",
		Data:      m,
		Storage:   s,
		EntityID:  entityID,
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// TestUserRegister uses a mock backend to check
// register create
func TestUserRegister(t *testing.T) {
	b, s := getTestBackend(t)

	t.Run("Create User Register - pass", func(t *testing.T) {
		resp, err := testTokenRegisterCreate(t, b, s, USER_ENTITY_ID)
		require.Nil(t, err)
		require.Nil(t, resp.Error())
		require.Nil(t, resp)
	})
}

func TestRegisterUser(t *testing.T) {
	mount := "pwmanager"

	t.Log("Test Registering User")
	{

		th, err := NewTestHarness(t, "TestRegisterUser", false)
		if err != nil {
			t.Fatalf("\tfailed to create test harness")
		}
		defer th.Teardown()

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

		for _, v := range users {
			v.WithUUK(th)
		}

		t.Logf("Register User Register")
		{

			for k, v := range users {
				if err := v.Client.Users().Register(mount, v.UUK); err != nil {
					th.Testing.Fatalf("\t%s error registering user %s: %s", FAILURE, k, err)
				}
				t.Logf("\t%s should be able to register user %s\n", SUCCESS, k)
			}
		}

		t.Logf("Register User Register Twice")
		{
			for k, v := range users {

				if err := v.Client.Users().Register(mount, v.UUK); err == nil {
					th.Testing.Fatalf("\t%sshould not be allowed to register %s more than once: %s", FAILURE, k, err)
				}
				t.Logf("\t%s should not be able to register %s twice\n", SUCCESS, k)
			}
		}

		t.Logf("Update User")
		{
			for k, v := range users {
				if err := v.Client.Users().Update(mount, v.LoginResponse.Auth.EntityID, v.UUK); err != nil {
					th.Testing.Fatalf("\t%s error updating user %s: %s", FAILURE, k, err)
				}
				t.Logf("\t%s should be able to update user %s\n", SUCCESS, k)
			}
		}

		t.Logf("Update another User")
		{

			stephen := users["stephen"]
			frank := users["frank"]

			if err := stephen.Client.Users().Update(mount, frank.LoginResponse.Auth.EntityID, stephen.UUK); err == nil {
				th.Testing.Fatalf("\t%s stephen should not be able to update frank: %s", FAILURE, err)
			}
			t.Logf("\t%s stephen should not be able to update frank\n", SUCCESS)

			if err := frank.Client.Users().Update(mount, stephen.LoginResponse.Auth.EntityID, stephen.UUK); err == nil {
				th.Testing.Fatalf("\t%s frank should not be able to update stephen: %s", FAILURE, err)
			}
			t.Logf("\t%s frank should not be able to update stephen\n", SUCCESS)
		}
	}
}
