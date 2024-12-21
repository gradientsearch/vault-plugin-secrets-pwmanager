package secretsengine

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"

	"github.com/lestrrat-go/jwx/v3/jwk"
	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/pbkdf2"
)

const (
	registerName = "testpwmgr"
	registerID   = "75e46130-bf36-11ef-8781-abab9f01266a"
)

// TestUserRegister uses a mock backend to check
// register create, read, update, and delete.
func TestUserRegister(t *testing.T) {
	b, s := getTestBackend(t)

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
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "register/" + name,
		Data:      d,
		Storage:   s,
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

func TestKeyDerivationHelper(t *testing.T) {
	password := "pwmgr"
	version := "H1"

	entityID := "52638ce9-c2a1-6a28-85ed-e61f3e9a697e"
	salt := []byte{154, 130, 13, 129, 242, 173, 81, 82, 69, 126, 236, 43, 235, 86, 104, 240}

	hash := sha256.New
	secret := "5h0TE4VU+qqC3GtoVYxw3EyXJFs+VYJGBQ0="

	kdf := hkdf.New(hash, []byte(version), salt, []byte(entityID))
	key1 := make([]byte, 32)
	if _, err := io.ReadFull(kdf, key1); err != nil {
		panic(err)
	}

	// Iterations and key length
	iter := 650000
	keyLen := 32

	dk := pbkdf2.Key([]byte(password), key1, iter, keyLen, sha1.New)

	hash2 := sha256.New
	kdf2 := hkdf.New(hash2, []byte(secret), []byte(entityID), dk)
	key2 := make([]byte, 32)
	if _, err := io.ReadFull(kdf2, key2); err != nil {
		panic(err)
	}

	// XOR the decoded byte slices
	xored := make([]byte, 32)
	for i := range dk {
		xored[i] = dk[i] ^ key2[i]
	}

	fmt.Println(xored)
	expected := []byte{67, 62, 174, 252, 121, 204, 236, 43, 36, 173, 78, 74, 136, 109, 107, 122, 206, 17, 186, 130, 104, 139, 199, 134, 101, 36, 244, 19, 15, 255, 36, 43}
	for i := range xored {
		if xored[i] != expected[i] {
			t.Fatalf("xored should match")
		}
	}

	// Generate an RSA private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Convert the private key to JWK format
	jwkKey, err := jwk.Import(privateKey)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Marshal the JWK key to JSON
	jwkJson, err := json.Marshal(jwkKey)
	if err != nil {
		fmt.Println(err)
		return
	}

	pubkey, err := jwkKey.PublicKey()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Marshal the JWK key to JSON
	pubJson, err := json.Marshal(pubkey)
	if err != nil {
		fmt.Println(err)
		return
	}

	symmetricKey := make([]byte, 32)
	_, err = rand.Read(symmetricKey)
	if err != nil {
		fmt.Println(err)
		return
	}

	c, err := aes.NewCipher(xored)
	if err != nil {
		fmt.Println(err)
		return
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Println(err)
		return
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		fmt.Println(err)
		return
	}

	encSymmetricKey := gcm.Seal(nonce, nonce, symmetricKey, nil)

	//Get the nonce size
	nonceSize := gcm.NonceSize()

	plaintext, err := gcm.Open(nil, nonce, encSymmetricKey[nonceSize:], nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	c2, err := aes.NewCipher(xored)
	if err != nil {
		fmt.Println(err)
		return
	}

	gcm2, err := cipher.NewGCM(c2)
	if err != nil {
		fmt.Println(err)
		return
	}

	nonce2 := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce2)
	if err != nil {
		fmt.Println(err)
		return
	}

	encPrivateKey := gcm2.Seal(nonce2, nonce2, jwkJson, nil)

	plaintext2, err := gcm2.Open(nil, nonce2, encPrivateKey[nonceSize:], nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s\n", hex.EncodeToString(plaintext))
	fmt.Printf("%s\n", (plaintext2))
	fmt.Printf("%s\n", jwkJson)
	fmt.Printf("%s\n", pubJson)
}
