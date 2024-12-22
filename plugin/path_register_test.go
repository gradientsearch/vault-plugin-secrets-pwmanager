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

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwe"
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

// compute fills in uuk from the derived 2SKD
func (uuk *UUK) compute(password []byte, mount []byte, secretKey []byte, entityID []byte) error {
	if id, err := uuid.GenerateUUID(); err != nil {
		return fmt.Errorf("failed to create uuid for UUK: %s", err)
	} else {
		uuk.UUID = id
	}

	initialSalt := []byte{154, 130, 13, 129, 242, 173, 81, 82, 69, 126, 236, 43, 235, 86, 104, 240} // 16 byte random salt (testing we hardcode it) not secret
	_, err := rand.Read(initialSalt)
	if err != nil {
		return err
	}

	// need this initial salt for decrypting UUK
	uuk.EncSymKey.P2s = hex.EncodeToString(initialSalt)

	//-----------------------------------------------------------------------
	// 2skd functions

	// // hkdf 1
	saltHash := hkdf.New(sha256.New, initialSalt, entityID, nil)
	saltDerivedKey := make([]byte, 32)
	if _, err := io.ReadFull(saltHash, saltDerivedKey); err != nil {
		return err
	}

	// // password hash
	iter := 650000
	uuk.EncSymKey.P2c = iter
	keyLen := 32

	passwordDerivedKey := pbkdf2.Key(password, saltDerivedKey, iter, keyLen, sha1.New)

	// // hkdf 2
	secretKeyHash := hkdf.New(sha256.New, secretKey, mount, nil)
	secretDerivedKey := make([]byte, 32)
	if _, err := io.ReadFull(secretKeyHash, secretDerivedKey); err != nil {
		return err
	}

	// // XOR passwordDk and secret
	twoSKD := make([]byte, 32)
	for i := range passwordDerivedKey {
		twoSKD[i] = passwordDerivedKey[i] ^ secretDerivedKey[i]
	}

	//-------------------------------------------------------
	// // create symmetric key

	symmetricKey := make([]byte, 32)
	_, err = rand.Read(symmetricKey)
	if err != nil {
		return err
	}

	// // encrypt symetric key with 2skd key

	//16, 24, or 32 bytes to select
	// AES-128, AES-192, or AES-256.
	// Since symmetric key is 32 bytes this is AES-256
	twoSkdCipher, err := aes.NewCipher(twoSKD)
	if err != nil {
		return err
	}

	twoSkdGcm, err := cipher.NewGCM(twoSkdCipher)
	if err != nil {
		return err
	}

	symIv := make([]byte, twoSkdGcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, symIv)
	if err != nil {
		return err
	}

	encSymKeyIvPrefix := twoSkdGcm.Seal(symIv, symIv, symmetricKey, nil)
	encSymKey := encSymKeyIvPrefix[twoSkdGcm.NonceSize():]

	// // add info to uuk - needed for users to decrypt this symmetric key later
	uuk.EncSymKey.Data = hex.EncodeToString(encSymKey)
	uuk.EncSymKey.Iv = hex.EncodeToString(symIv)
	uuk.EncSymKey.Enc = "A256GCM"
	uuk.EncSymKey.Kid = uuk.UUID
	uuk.EncSymKey.Alg = "pbkdf2-hkdf"

	//------------------------------------------------------------
	// Private Key
	// Generate an RSA private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// // encrypt private key with symmetric key
	symmetricKeyCipher, err := aes.NewCipher(symmetricKey)
	if err != nil {
		return err
	}

	symmetricKeyGcm, err := cipher.NewGCM(symmetricKeyCipher)
	if err != nil {
		return err
	}

	privKeyIV := make([]byte, symmetricKeyGcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, privKeyIV)
	if err != nil {
		return err
	}

	// Convert the private key to JWK format
	jwkKey, err := jwk.Import(privateKey)
	if err != nil {
		return err
	}

	// Marshal the JWK key to JSON
	jwkJson, err := json.Marshal(jwkKey)
	if err != nil {
		return err
	}

	encPrivateKeyIvPrefix := symmetricKeyGcm.Seal(privKeyIV, privKeyIV, jwkJson, nil)
	encPrivKey := encPrivateKeyIvPrefix[symmetricKeyGcm.NonceSize():]
	// add info to the encPrivKey
	uuk.EncPriKey.Data = hex.EncodeToString(encPrivKey)
	uuk.EncPriKey.Kid = uuk.UUID
	uuk.EncPriKey.Iv = hex.EncodeToString(privKeyIV)
	uuk.EncPriKey.Enc = "A256GCM"

	//-------------------------------------------------------
	// pubkey
	pubkey, err := jwkKey.PublicKey()
	if err != nil {
		fmt.Println(err)
	}

	pubkey.Set("kid", uuk.UUID)

	uuk.PubKey = pubkey

	return nil
}

func TestKeyDerivation(t *testing.T) {
	// set up required values
	entityID := []byte("52638ce9-c2a1-6a28-85ed-e61f3e9a697e") // vault entity id (use token to look this up) - not secret
	password := []byte("super-secret")                         // user password secret
	mount := []byte("pwmgr")                                   // vault mount - not secret
	version := "H1"                                            // version of pwmgr - not secret
	randomSeq := make([]byte, 32)                              // rand 32 byte sequence - secret
	if _, err := io.ReadFull(rand.Reader, randomSeq); err != nil {
		t.Fatalf("error creating random sequence of characters: %s", err)
	}

	secretKey := []byte(fmt.Sprintf("%s%s%s", version, mount, randomSeq)) // combination Secret ID - secret

	// create UUK
	uuk := UUK{}
	err := uuk.compute(password, mount, secretKey, entityID)
	if err != nil {
		t.Fatalf("error creating uuk: %s", err)
	}

	initialSalt, err := hex.DecodeString(uuk.EncSymKey.P2s)
	if err != nil {
		t.Fatalf("error decoding P2s")
	}
	// // hkdf 1
	saltHash := hkdf.New(sha256.New, initialSalt, entityID, nil)
	saltDerivedKey := make([]byte, 32)
	if _, err := io.ReadFull(saltHash, saltDerivedKey); err != nil {
		t.Fatalf("error computing salt derived key: %s", err)
	}

	// // password hash
	keyLen := 32

	passwordDerivedKey := pbkdf2.Key(password, saltDerivedKey, uuk.EncSymKey.P2c, keyLen, sha1.New)

	// // hkdf 2
	secretKeyHash := hkdf.New(sha256.New, secretKey, mount, nil)
	secretDerivedKey := make([]byte, 32)
	if _, err := io.ReadFull(secretKeyHash, secretDerivedKey); err != nil {
		t.Fatalf("error computing secret derived key: %s", err)
	}

	// // XOR passwordDk and secret
	twoSKD := make([]byte, 32)
	for i := range passwordDerivedKey {
		twoSKD[i] = passwordDerivedKey[i] ^ secretDerivedKey[i]
	}

	//-------------------------------------------------------
	// // retrieve symmetric key

	//16, 24, or 32 bytes to select
	// AES-128, AES-192, or AES-256.
	// Since symmetric key is 32 bytes this is AES-256
	twoSkdCipher, err := aes.NewCipher(twoSKD)
	if err != nil {
		t.Fatal(err)
	}

	twoSkdGcm, err := cipher.NewGCM(twoSkdCipher)
	if err != nil {
		t.Fatal(err)
	}

	symIv, err := hex.DecodeString(uuk.EncSymKey.Iv)
	if err != nil {
		t.Fatalf("error decoding symmetric iv: %s", err)
	}

	encSymKey, err := hex.DecodeString(uuk.EncSymKey.Data)
	if err != nil {
		t.Fatalf("error decoding symmetric data: %s", err)
	}

	symmetricKey, err := twoSkdGcm.Open(nil, symIv, encSymKey, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// retrieve private key

	symmetricKeyCipher, err := aes.NewCipher(symmetricKey)
	if err != nil {
		t.Fatal(err)
	}

	symmetricKeyGcm, err := cipher.NewGCM(symmetricKeyCipher)
	if err != nil {
		t.Fatal(err)
	}

	encPrivateKey, err := hex.DecodeString(uuk.EncPriKey.Data)
	if err != nil {
		t.Fatalf("error decoding private key encrypted data: %s", err)
	}

	privKeyIV, err := hex.DecodeString(uuk.EncPriKey.Iv)
	if err != nil {
		t.Fatalf("error decoding encPrivKey Iv: %s", err)
	}

	jwkJson, err := symmetricKeyGcm.Open(nil, privKeyIV, encPrivateKey, nil)
	if err != nil {
		t.Fatalf("failed to decrypt private key: %s", err)
	}

	privKey, err := jwk.ParseKey(jwkJson)
	if err != nil {
		t.Fatalf("error parsing private key: %s", err)
	}

	// encrypt data
	const payload = `gopher`
	encrypted, err := jwe.Encrypt([]byte(payload), jwe.WithKey(jwa.RSA_OAEP(), uuk.PubKey))
	if err != nil {
		fmt.Printf("failed to encrypt payload: %s\n", err)
		return
	}

	decrypted, err := jwe.Decrypt(encrypted, jwe.WithKey(jwa.RSA_OAEP(), privKey))
	if err != nil {
		fmt.Printf("failed to decrypt payload: %s\n", err)
		return
	}
	if string(decrypted) != payload {
		t.Fatalf("expected %s to equal %s", decrypted, payload)
	}
	fmt.Printf("encrypted: %s\n", encrypted)
	fmt.Printf("decrypted: %s\n", decrypted)
}

type EncryptedPrivateKey struct {
	Kid  string `json:"kid"`
	Enc  string `json:"enc"`
	Iv   string `json:"iv"`
	Data string `json:"data"`
	Cty  string `json:"cty"`
}

type EncryptedSymmetricKey struct {
	Kid  string `json:"kid"`
	Enc  string `json:"enc"`
	Iv   string `json:"iv"`
	Data string `json:"data"`
	Cty  string `json:"cty"`
	Alg  string `json:"alg"`
	P2c  int    `json:"p2c"`
	P2s  string `json:"p2s"`
}

// user unlock key
type UUK struct {
	UUID        string                `json:"uuid"`
	EncSymKey   EncryptedSymmetricKey `json:"encSymKey"`
	EncryptedBy string                `json:"encryptedBy"`
	EncPriKey   EncryptedPrivateKey   `json:"encPriKey"`
	PubKey      interface{}           `json:"pubkey"`
}
