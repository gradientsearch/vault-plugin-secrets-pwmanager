package secretsengine

import (
	"context"
	"crypto/rand"
	"os"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

// TestUserRegister uses a mock backend to check
// register create
func TestUserRegister(t *testing.T) {
	b, s := getTestBackend(t)

	t.Run("Create User Register - pass", func(t *testing.T) {
		resp, err := testTokenRegisterCreate(t, b, s)
		require.Nil(t, err)
		require.Nil(t, resp.Error())
		require.Nil(t, resp)
	})
}

// Utility function to create a register while, returning any response (including errors)
func testTokenRegisterCreate(t *testing.T, b *pwmgrBackend, s logical.Storage) (*logical.Response, error) {
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
		stephen := users["stephen"]
		stephen.WithUUK(th)

		t.Logf("Register User")
		{
			if err := stephen.Client.PwManager().Register(mount, stephen.UUK); err != nil {
				th.Testing.Fatalf("\t%s error registering user: %s", FAILURE, err)
			}
			t.Logf("\t %s should be able to register user\n", SUCCESS)

			if err := stephen.Client.PwManager().Register(mount, stephen.UUK); err == nil {
				th.Testing.Fatalf("\t%sshould not be allowed to register more than once: %s", FAILURE, err)
			}
			t.Logf("\t %s should not be able to register twice\n", SUCCESS)
		}
	}
}
