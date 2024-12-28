
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
    P2c: number
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

export function toHex(plain: string) {
    return plain.split("")
        .map(c => c.charCodeAt(0).toString(16).padStart(2, "0"))
        .join("");
}

function toString(hex: string): string {
    return hex.split(/(\w\w)/g)
        .filter(p => !!p)
        .map(c => String.fromCharCode(parseInt(c, 16)))
        .join("")
}

function toString2(hex: string): Uint8Array {
    let sa = toString(hex).split(",")
    let a = new Uint8Array(sa.length)
    for (let i = 0; i < sa.length; i++) {
        a[i] = Number(sa[i])
    }
    return a
}


function withInitializationSalt(uuk: UUK): UUK {
    uuk.EncSymKey.P2s = toHex(crypto.getRandomValues(new Uint8Array(16)).toString())
    return uuk;
}

function withPasswordIterations(uuk: UUK, iterations: number): UUK {
    uuk.EncSymKey.P2c = iterations

    return uuk
}

export async function twoSkd(uuk: UUK, password: Uint8Array, mount: Uint8Array, secretKey: Uint8Array, entityID: Uint8Array): Promise<[UUK, Uint8Array]> {

    let rawKey = toString2(uuk.EncSymKey.P2s)
    let initialSalt = await crypto.subtle.importKey("raw", rawKey, "HKDF", false, [
        'deriveBits'
    ]);

    // HKDF 1
    const salt = entityID
    const info = new TextEncoder().encode('2SKD HKDF 1');
    const hkdf_params = { name: 'HKDF', hash: 'SHA-256', salt, info };
    const saltDerivedKey = await crypto.subtle.deriveBits(
        hkdf_params,
        initialSalt,
        32 * 8
    );


    const enc = new TextEncoder();

    let passwordKey = await crypto.subtle.importKey(
        "raw",
        password,
        "PBKDF2",
        false,
        ["deriveBits"],
    );


    let passwordDerivedKey = await crypto.subtle.deriveBits(
        {
            name: "PBKDF2",
            salt: saltDerivedKey,
            iterations: 100000,
            hash: "SHA-256",
        },
        passwordKey,
        32 * 8
    );

    // hkdf 2
    let secretKeyCrypto = await crypto.subtle.importKey("raw", secretKey, "HKDF", false, [
        'deriveBits'
    ]);

    // HKDF 2
    const hkdf2Salt = mount
    const hkdf2Info = new TextEncoder().encode('2SKD HKDF 2');
    const hkdf2_params = { name: 'HKDF', hash: 'SHA-256', salt: hkdf2Salt, info: hkdf2Info };
    const secretDerivedKey = await crypto.subtle.deriveBits(
        hkdf2_params,
        secretKeyCrypto,
        32 * 8
    );

    let passwordDerivedKeyBytes = new Uint8Array(passwordDerivedKey);
    let secretDerivedKeyBytes = new Uint8Array(secretDerivedKey);

    // XOR passwordDk and secret
    let twoSKD = new Uint8Array(32)
    for (let i = 0; i < 32; i++) {
        twoSKD[i] = passwordDerivedKeyBytes[i] ^ secretDerivedKeyBytes[i]
    }

    return [uuk, twoSKD]
}


function newUUK(): UUK {
    let uuk: UUK = {
        EncPriKey: {
            Kid: "",
            Enc: "",
            Iv: "",
            Data: "",
            Cty: ""
        },
        EncSymKey: {
            Kid: "",
            Enc: "",
            Iv: "",
            Data: "",
            Cty: "",
            Alg: "",
            P2c: 0,
            P2s: ""
        },
        PubKey: {
            E: "",
            Kid: "",
            Kty: "",
            N: ""
        },
        UUID: "",
        EncryptedBy: ""
    }
    uuk.UUID = crypto.randomUUID();
    return uuk
}

export async function buildUUK(password: Uint8Array, mount: Uint8Array, secretKey: Uint8Array, entityID: Uint8Array) {
    let uuk = newUUK()
    uuk = withInitializationSalt(uuk)
    uuk = withPasswordIterations(uuk, 650000)
    const textEncoder = new TextEncoder();

    let twoSkdHash: Uint8Array
    [uuk, twoSkdHash] = await twoSkd(uuk, password, mount, secretKey, entityID)

    return await crypto.subtle.generateKey('X25519', false /* extractable */, ['deriveKey']);
}
