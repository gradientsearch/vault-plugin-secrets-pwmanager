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
	mrand "math/rand"
	"strconv"
	"testing"
	"time"

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

func init() {
	mrand.Seed(time.Now().UnixNano())
}

var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[mrand.Intn(len(letters))]
	}
	return string(b)
}

func TestKeyDerivation(t *testing.T) {
	password := "super-secret"                                                                      // user password secret
	version := "H1"                                                                                 // version of pwmgr - not secret
	mount := "pwmgr"                                                                                // secret mount - not secret
	secret := "BNIFFWPTTPOKHRYEHEPVJHFFJOFMGDRDMCBMRW"                                              // random secret user creates locally - secret
	secretID := fmt.Sprintf("%s%s%s", version, mount, secret)                                       // combination Secret ID - secret
	entityID := "52638ce9-c2a1-6a28-85ed-e61f3e9a697e"                                              // vault entity id (use token to look this up) - not secret
	initialSalt := []byte{154, 130, 13, 129, 242, 173, 81, 82, 69, 126, 236, 43, 235, 86, 104, 240} // 16 byte random salt (testing we hardcode it) not secret

	saltKdf := hkdf.New(sha256.New, initialSalt, []byte(entityID), []byte("HKDF"))
	salt := make([]byte, 32)
	if _, err := io.ReadFull(saltKdf, salt); err != nil {
		panic(err)
	}

	iter := 650000
	keyLen := 32

	passwordDk := pbkdf2.Key([]byte(password), salt, iter, keyLen, sha1.New)

	secretIdKdf := hkdf.New(sha256.New, []byte(secretID), []byte(mount), []byte("HKDF"))
	secretIdDk := make([]byte, 32)
	if _, err := io.ReadFull(secretIdKdf, secretIdDk); err != nil {
		panic(err)
	}

	// XOR the decoded byte slices
	twoSKD := make([]byte, 32)
	for i := range passwordDk {
		// xor dk and key2
		twoSKD[i] = passwordDk[i] ^ secretIdDk[i]
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

	symmetricKey, _ := hex.DecodeString("59a8e993838166a888bd0940bd30f640ba41807b5e2df42e15a67eb057f73d7e")

	// symmetricKey := make([]byte, 32)
	// _, err = rand.Read(symmetricKey)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	c, err := aes.NewCipher(twoSKD)
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

	encSymmetricKey := gcm.Seal(nonce, nonce, []byte(symmetricKey), nil)

	//Get the nonce size
	nonceSize := gcm.NonceSize()

	plaintext, err := gcm.Open(nil, nonce, encSymmetricKey[nonceSize:], nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	c2, err := aes.NewCipher(twoSKD)
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

	fmt.Printf("secretID: %s\n", secretID)
	fmt.Printf("%s\n", hex.EncodeToString(plaintext))
	fmt.Printf("%s\n", (plaintext2))
	fmt.Printf("%s\n", jwkJson)
	fmt.Printf("%s\n", pubJson)
}
