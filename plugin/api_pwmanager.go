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

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/api"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwe"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/pbkdf2"
)

// PwManager is used to perform PwManager operations on Vault.
type PwManager struct {
	c *api.Client
}

// PwManager is used to return the client for pwmanger API calls.
func (c *pwmanagerClient) PwManager() *PwManager {
	return &PwManager{c: c.c}
}

func (c *PwManager) Register(mount string, uuk UUK) error {
	r := c.c.NewRequest("POST", fmt.Sprintf("/v1/%s/register", mount))
	if err := r.SetJSONBody(uuk); err != nil {
		return err
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

type EncPriKey struct {
	// uuid
	Kid string `json:"kid"`
	// encoding of data e.g. A256GCM
	Enc string `json:"enc"`
	// initialization vector used to encrypt the priv key
	Iv string `json:"iv"`
	// the encrypted priv key
	Data string `json:"data"`
	// format used for encrypted data e.g JWK format
	Cty string `json:"cty"`
}

// EncSymKey contains the data required to unlock the
// users private key.
type EncSymKey struct {
	// uuid of the private key
	Kid string `json:"kid"`
	// encoding used to encrypt the data e.g. A256GCM
	Enc string `json:"enc"`
	// initialization
	Iv string `json:"iv"`
	// encrypted symmetric key
	Data string `json:"data"`
	// content type
	Cty string `json:"cty"`
	// the algorithm used to encrypt the EncSymKey e.g. 2SKD PBDKF2-HKDF
	Alg string `json:"alg"`
	// PBDKF2 iterations e.g. 650000
	P2c int `json:"p2c"`
	// initial 16 byte random sequence for secret key derivation.
	// used in the first hkdf function call
	P2s string `json:"p2s"`
}

// user unlock key
// The secret key encrypts the EncSymKey, the EncSymKey
// encrypts the users PrivateKey
type UUK struct {
	// uuid of priv key
	UUID string `json:"uuid"`
	// symmetric key used to encrypt the EncPriKey
	EncSymKey EncSymKey `json:"enc_sym_key"`
	// mp a.k.a secret key
	EncryptedBy string `json:"encrypted_by"`
	// priv key used to encrypt `Safe` data
	EncPriKey EncPriKey `json:"enc_pri_key"`
	// pub key of the private key
	PubKey interface{} `json:"pubKey"`
}

// withInitializationSalt generates a random 16 byte salt and stores the result in UUK.EncSymKey.P2s
// this salt is required in the 2SKD func
func (uuk *UUK) withInitializationSalt() error {
	initialSalt := make([]byte, 16)
	_, err := rand.Read(initialSalt)
	if err != nil {
		return fmt.Errorf("error generating initialization salt: %s", err)
	}

	// need this initial salt for creating 2SKD and decrypting symkey in the future
	uuk.EncSymKey.P2s = hex.EncodeToString(initialSalt)
	return nil
}

// withPasswordIterations sets the iteration count used in the 2SKD func
func (uuk *UUK) withPasswordIterations(iterations int) {
	uuk.EncSymKey.P2c = iterations
}

// withEncSymKey creates a symmetric key, encrypts it using the 2SKD key and stores the encrypted
// value in UUK.EncSymKey.Data, sets non secret attributes in UUK.EncSymKey.
// Returns the plaintext symmetric key to be used for encrypting the UUK.EncPriKey
func (uuk *UUK) withEncSymKey(twoSKD []byte) ([]byte, error) {
	symmetricKey := make([]byte, 32)
	_, err := rand.Read(symmetricKey)
	if err != nil {
		return nil, err
	}

	//16, 24, or 32 bytes to select
	// AES-128, AES-192, or AES-256.
	// Since symmetric key is 32 bytes this is AES-256
	c, err := aes.NewCipher(twoSKD)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return nil, err
	}

	encSymKeyIvPrefix := gcm.Seal(iv, iv, symmetricKey, nil)
	encSymKey := encSymKeyIvPrefix[gcm.NonceSize():]

	// add info to uuk - needed for users to decrypt this symmetric key later
	uuk.EncSymKey.Data = hex.EncodeToString(encSymKey)
	uuk.EncSymKey.Iv = hex.EncodeToString(iv)
	uuk.EncSymKey.Enc = "A256GCM"
	uuk.EncSymKey.Kid = uuk.UUID
	uuk.EncSymKey.Alg = "pbkdf2-hkdf"

	return symmetricKey, nil
}

// withEncPriKey encrypts creates a private key and encrypts it using the symmetric
// key stored in UUK.EncSymKey, stores the encrypted value in UUK.EncPriKey.Data, sets
// non secret attributes in UUK.EncPriKey, and returns the rsa private key
func (uuk *UUK) withEncPriKey(symmetricKey []byte) (*rsa.PrivateKey, error) {
	//------------------------------------------------------------
	// Private Key
	// Generate an RSA private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	// encrypt private key with symmetric key
	c, err := aes.NewCipher(symmetricKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return nil, err
	}

	// Convert the private key to JWK format
	jwkKey, err := jwk.Import(privateKey)
	if err != nil {
		return nil, err
	}

	// Marshal the JWK key to JSON
	jwkJson, err := json.Marshal(jwkKey)
	if err != nil {
		return nil, err
	}

	encPrivateKeyIvPrefix := gcm.Seal(iv, iv, jwkJson, nil)
	encPrivKey := encPrivateKeyIvPrefix[gcm.NonceSize():]

	// add info to the EncPrivKey
	uuk.EncPriKey.Data = hex.EncodeToString(encPrivKey)
	uuk.EncPriKey.Kid = uuk.UUID
	uuk.EncPriKey.Iv = hex.EncodeToString(iv)
	uuk.EncPriKey.Enc = "A256GCM"

	return privateKey, nil
}

// withPubKey extracts the pubKey from the private key and assigns it to UUK.PubKey
func (uuk *UUK) withPubKey(prikey *rsa.PrivateKey) error {
	jwkPriKey, err := jwk.Import(prikey)
	if err != nil {
		return err
	}

	pubkey, err := jwkPriKey.PublicKey()
	if err != nil {
		return err
	}

	pubkey.Set("kid", uuk.UUID)

	uuk.PubKey = pubkey
	return nil
}

func (uuk *UUK) withEncryptedBy(eb string) {
	uuk.EncryptedBy = eb
}

// derives the 2SKD from provided parameters
func (uuk *UUK) twoSkd(password, mount, secretKey, entityID []byte) ([]byte, error) {
	initialSalt, err := hex.DecodeString(uuk.EncSymKey.P2s)
	if err != nil {
		return nil, err
	}
	// hkdf 1
	saltHash := hkdf.New(sha256.New, initialSalt, entityID, nil)
	saltDerivedKey := make([]byte, 32)
	if _, err := io.ReadFull(saltHash, saltDerivedKey); err != nil {
		return nil, err
	}

	// pbkdf2
	keyLen := 32
	passwordDerivedKey := pbkdf2.Key(password, saltDerivedKey, uuk.EncSymKey.P2c, keyLen, sha1.New)

	// hkdf 2
	secretKeyHash := hkdf.New(sha256.New, secretKey, mount, nil)
	secretDerivedKey := make([]byte, 32)
	if _, err := io.ReadFull(secretKeyHash, secretDerivedKey); err != nil {
		return nil, err
	}

	// XOR passwordDk and secret
	twoSKD := make([]byte, 32)
	for i := range passwordDerivedKey {
		twoSKD[i] = passwordDerivedKey[i] ^ secretDerivedKey[i]
	}
	return twoSKD, nil
}

// Build fills in uuk from the derived 2SKD
// facade pattern: The order of the with funcs called is important
// since building modifies the required UUK attributes to properly
// build the UUK data.
// example json output of UUK struct
//
//	{
//	  "uuid": "0bbf993d-8e10-6dd0-1aa3-80019b69e332",
//	  "enc_sym_key":
//	    {
//	      "kid": "0bbf993d-8e10-6dd0-1aa3-80019b69e332",
//	      "enc": "A256GCM",
//	      "iv": "05e89fa122ecae403feda8dd",
//	      "data": "95af1a39e798f9f4bcd401c778f2f4ed40eb9f56da8c8f9a1d7b0601777b37bfd2c369af0076f5b94e4be1622003fa8b",
//	      "cty": "",
//	      "alg": "pbkdf2-hkdf",
//	      "p2c": 650000,
//	      "p2s": "26d6ee5149c95425a0251651f8b07ac0",
//	    },
//	  "encrypted_by": "mp",
//	  "enc_pri_key":
//	    {
//	      "kid": "0bbf993d-8e10-6dd0-1aa3-80019b69e332",
//	      "enc": "A256GCM",
//	      "iv": "647ffb7da1e70cd48370101b",
//	      "data": "809311c7ec73c746fb5c0bfa78a70a010b81853...",
//	      "cty": "",
//	    },
//	  "pubkey":
//	    {
//	      "e": "AQAB",
//	      "kid": "0bbf993d-8e10-6dd0-1aa3-80019b69e332",
//	      "kty": "RSA",
//	      "n": "3gB8w0CbpMnZCiA6QSUeCXyAsx9v...",
//	    },
//	}
func (uuk *UUK) Build(password, mount, secretKey, entityID []byte) error {
	if id, err := uuid.GenerateUUID(); err != nil {
		return fmt.Errorf("failed to create uuid for UUK: %s", err)
	} else {
		uuk.UUID = id
	}

	if err := uuk.withInitializationSalt(); err != nil {
		return err
	}

	uuk.withPasswordIterations(650000)

	twoSKD, err := uuk.twoSkd(password, mount, secretKey, entityID)
	if err != nil {
		return err
	}

	symmetricKey, err := uuk.withEncSymKey(twoSKD)
	if err != nil {
		return err
	}

	priKey, err := uuk.withEncPriKey(symmetricKey)
	if err != nil {
		return err
	}

	if err := uuk.withPubKey(priKey); err != nil {
		return err
	}

	uuk.withEncryptedBy("mp")

	return nil
}

// DecryptEncPriKey decrypts the UUK EncPriKey using the EncSymKey and the the users SecretKey
// returns priv key used to encrypt payloads
func (uuk *UUK) DecryptEncPriKey(password, mount, secretKey, entityID []byte) (jwk.Key, error) {
	twoSKD, err := uuk.twoSkd(password, mount, secretKey, entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to create 2SKD %s", err)
	}
	twoSkdCipher, err := aes.NewCipher(twoSKD)
	if err != nil {
		return nil, err
	}

	twoSkdGcm, err := cipher.NewGCM(twoSkdCipher)
	if err != nil {
		return nil, err
	}

	symIv, err := hex.DecodeString(uuk.EncSymKey.Iv)
	if err != nil {
		return nil, fmt.Errorf("error decoding symmetric iv: %s", err)
	}

	encSymKey, err := hex.DecodeString(uuk.EncSymKey.Data)
	if err != nil {
		return nil, fmt.Errorf("error decoding symmetric data: %s", err)
	}

	symmetricKey, err := twoSkdGcm.Open(nil, symIv, encSymKey, nil)
	if err != nil {
		return nil, err
	}

	// retrieve private key

	symmetricKeyCipher, err := aes.NewCipher(symmetricKey)
	if err != nil {
		return nil, err
	}

	symmetricKeyGcm, err := cipher.NewGCM(symmetricKeyCipher)
	if err != nil {
		return nil, err
	}

	encPrivateKey, err := hex.DecodeString(uuk.EncPriKey.Data)
	if err != nil {
		return nil, fmt.Errorf("error decoding private key encrypted data: %s", err)
	}

	privKeyIV, err := hex.DecodeString(uuk.EncPriKey.Iv)
	if err != nil {
		return nil, fmt.Errorf("error decoding encPrivKey Iv: %s", err)
	}

	jwkJson, err := symmetricKeyGcm.Open(nil, privKeyIV, encPrivateKey, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt private key: %s", err)
	}

	priKey, err := jwk.ParseKey(jwkJson)
	if err != nil {
		return nil, fmt.Errorf("error parsing private key: %s", err)
	}

	return priKey, nil
}

// Convenience func to encrypt data encrypted with users pubkey
func (uuk *UUK) Encrypt(payload string) ([]byte, error) {
	encrypted, err := jwe.Encrypt([]byte(payload), jwe.WithKey(jwa.RSA_OAEP(), uuk.PubKey))
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt payload: %s\n", err)
	}
	return encrypted, nil
}

// Convenience func to decrypt data encrypted with users prikey
func (uuk *UUK) Decrypt(encrypted []byte, priKey jwk.Key) ([]byte, error) {
	decrypted, err := jwe.Decrypt(encrypted, jwe.WithKey(jwa.RSA_OAEP(), priKey))
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt payload: %s\n", err)

	}
	return decrypted, nil
}

// returns json string of the UUK struct
func (uuk *UUK) Json() (string, error) {
	if j, err := json.Marshal(uuk); err != nil {
		return "", err
	} else {
		return string(j), nil
	}
}

// returns json string of the UUK struct
func (uuk *UUK) Map() (map[string]interface{}, error) {

	j, err := json.Marshal(uuk)
	if err != nil {
		return nil, err
	}
	m := make(map[string]any)
	if err := json.Unmarshal(j, &m); err != nil {
		return nil, err
	}

	return m, nil
}
