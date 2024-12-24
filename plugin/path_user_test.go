package secretsengine

import (
	"context"
	"strconv"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

const (
	userName   = "testpwmgr"
	testTTL    = int64(120)
	testMaxTTL = int64(3600)
)

// TestUserUser uses a mock backend to check
// user create, read, update, and delete.
func TestUserUser(t *testing.T) {
	b, s := getTestBackend(t)

	t.Run("List All Users", func(t *testing.T) {
		for i := 1; i <= 10; i++ {
			_, err := testTokenUserCreate(t, b, s,
				userName+strconv.Itoa(i),
				map[string]interface{}{
					"username": userID,
					"ttl":      testTTL,
					"max_ttl":  testMaxTTL,
				})
			require.NoError(t, err)
		}

		resp, err := testTokenUserList(t, b, s)
		require.NoError(t, err)
		require.Len(t, resp.Data["keys"].([]string), 10)
	})

	t.Run("Create User User - pass", func(t *testing.T) {
		resp, err := testTokenUserCreate(t, b, s, userName, map[string]interface{}{
			"username": userID,
			"ttl":      testTTL,
			"max_ttl":  testMaxTTL,
		})

		require.Nil(t, err)
		require.Nil(t, resp.Error())
		require.Nil(t, resp)
	})

	t.Run("Read User User", func(t *testing.T) {
		resp, err := testTokenUserRead(t, b, s)

		require.Nil(t, err)
		require.Nil(t, resp.Error())
		require.NotNil(t, resp)
		require.Equal(t, resp.Data["username"], userID)
	})
	t.Run("Update User User", func(t *testing.T) {
		resp, err := testTokenUserUpdate(t, b, s, map[string]interface{}{
			"ttl":     "1m",
			"max_ttl": "5h",
		})

		require.Nil(t, err)
		require.Nil(t, resp.Error())
		require.Nil(t, resp)
	})

	t.Run("Re-read User User", func(t *testing.T) {
		resp, err := testTokenUserRead(t, b, s)

		require.Nil(t, err)
		require.Nil(t, resp.Error())
		require.NotNil(t, resp)
		require.Equal(t, resp.Data["username"], userID)
	})

	t.Run("Delete User User", func(t *testing.T) {
		_, err := testTokenUserDelete(t, b, s)

		require.NoError(t, err)
	})
}

// Utility function to create a user while, returning any response (including errors)
func testTokenUserCreate(t *testing.T, b *pwmgrBackend, s logical.Storage, name string, d map[string]interface{}) (*logical.Response, error) {
	t.Helper()
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "user/" + name,
		Data:      d,
		Storage:   s,
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Utility function to update a user while, returning any response (including errors)
func testTokenUserUpdate(t *testing.T, b *pwmgrBackend, s logical.Storage, d map[string]interface{}) (*logical.Response, error) {
	t.Helper()
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "user/" + userName,
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

// Utility function to read a user and return any errors
func testTokenUserRead(t *testing.T, b *pwmgrBackend, s logical.Storage) (*logical.Response, error) {
	t.Helper()
	return b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "user/" + userName,
		Storage:   s,
	})
}

// Utility function to list users and return any errors
func testTokenUserList(t *testing.T, b *pwmgrBackend, s logical.Storage) (*logical.Response, error) {
	t.Helper()
	return b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ListOperation,
		Path:      "user/",
		Storage:   s,
	})
}

// Utility function to delete a user and return any errors
func testTokenUserDelete(t *testing.T, b *pwmgrBackend, s logical.Storage) (*logical.Response, error) {
	t.Helper()
	return b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.DeleteOperation,
		Path:      "user/" + userName,
		Storage:   s,
	})
}
