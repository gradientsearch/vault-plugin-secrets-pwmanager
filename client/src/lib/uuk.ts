
interface EncPriKey {
	Kid: string
	// encoding of data e.g. A256GCM
	Enc: string
	// initialization vector used to encrypt the priv key
	Iv: string
	// the encrypted priv key
	Data: string
	// format used for encrypted data e.g JWK format
	Cty: string

}

// EncSymKey contains the data required to unlock the
// users private key.
interface EncSymKey {
    // uuid of the private key
    Kid: string
    // encoding used to encrypt the data e.g. A256GCM
    Enc: string
    // initialization
    Iv: string
    // encrypted symmetric key
    Data: string
    // content type
    Cty: string
    // the algorithm used to encrypt the EncSymKey e.g. 2SKD PBDKF2-HKDF
    Alg: string
    // PBDKF2 iterations e.g. 650000
    P2c: string
    // initial 16 byte random sequence for secret key derivation.
    // used in the first hkdf function call
    P2s: string
}


// user unlock key
// The secret key encrypts the EncSymKey, the EncSymKey
// encrypts the users PrivateKey
export interface UUK {
    // uuid of priv key
    UUID: string
    // symmetric key used to encrypt the EncPriKey
    EncSymKey: EncSymKey
    // mp a.k.a secret key
    EncryptedBy: string
    // priv key used to encrypt
    EncPriKey: EncPriKey
    // pub key of the private key
    PubKey: any
}


async function uuk() {
    await crypto.subtle.generateKey('X25519', false /* extractable */, ['deriveKey']);
}
