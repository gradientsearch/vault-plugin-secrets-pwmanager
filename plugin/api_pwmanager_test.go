package secretsengine

import (
	"crypto/rand"
	"fmt"
	"io"
	"testing"
)

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
	err := uuk.Build(password, mount, secretKey, entityID)
	if err != nil {
		t.Fatalf("error creating uuk: %s", err)
	}

	priKey, err := uuk.DecryptEncPriKey(password, mount, secretKey, entityID)
	if err != nil {
		t.Fatalf("error decrypting prikey: %s", err)
	}

	// encrypt data
	const payload = `gopher`
	encrypted, err := uuk.Encrypt(payload)
	if err != nil {
		t.Fatalf("error encrypting payload: %s", err)
	}

	decrypted, err := uuk.Decrypt(encrypted, priKey)
	if err != nil {
		t.Fatalf("error decrypting payload: %s", err)
	}
	if string(decrypted) != payload {
		t.Fatalf("expected %s to equal %s", decrypted, payload)
	}

	fmt.Println(uuk.Json())
}
