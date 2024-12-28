
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

interface PubKey {
    E: string
    Kid: string
    Kty: string
    N: string
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
    PubKey: PubKey
}

function toHex(plain: string) {
    return plain.split("")
        .map(c => c.charCodeAt(0).toString(16).padStart(2, "0"))
        .join("");
}

function toString(hex: string) {
    return hex.split(/(\w\w)/g)
        .filter(p => !!p)
        .map(c => String.fromCharCode(parseInt(c, 16)))
        .join("")
}


function withInitializationSalt(uuk: UUK): UUK {
    uuk.EncSymKey.Iv = toHex(crypto.getRandomValues(new Uint8Array(16)).toString());
    return uuk;
}

function newUUK(): UUK {
    let uuk: UUK = {
        EncPriKey: {},
        EncSymKey: {},
        PubKey: {},
    }
    uuk.UUID = crypto.randomUUID();
    return uuk
}

export async function buildUUK() {
    let uuk = newUUK()
    uuk = withInitializationSalt(uuk)

    return await crypto.subtle.generateKey('X25519', false /* extractable */, ['deriveKey']);
}
