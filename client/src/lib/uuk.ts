interface EncPriKey {
	Kid: string;
	// encoding of data e.g. A256GCM
	Enc: string;
	// initialization vector used to encrypt the priv key
	Iv: string;
	// the encrypted priv key
	Data: string;
	// format used for encrypted data e.g JWK format
	Cty: string;
}

// EncSymKey contains the data required to unlock the
// users private key.
interface EncSymKey {
	// uuid of the private key
	Kid: string;
	// encoding used to encrypt the data e.g. A256GCM
	Enc: string;
	// initialization
	Iv: string;
	// encrypted symmetric key
	Data: string;
	// content type
	Cty: string;
	// the algorithm used to encrypt the EncSymKey e.g. 2SKD PBDKF2-HKDF
	Alg: string;
	// PBDKF2 iterations e.g. 650000
	P2c: number;
	// initial 16 byte random sequence for secret key derivation.
	// used in the first hkdf function call
	P2s: string;
}

interface PubKey {
	E?: string;
	N?: string;
	Kid: string;
	Kty: string;
}

// user unlock key
// The secret key encrypts the EncSymKey, the EncSymKey
// encrypts the users PrivateKey
export interface UUK {
	// uuid of priv key
	UUID: string;
	// symmetric key used to encrypt the EncPriKey
	EncSymKey: EncSymKey;
	// mp a.k.a secret key
	EncryptedBy: string;
	// priv key used to encrypt
	EncPriKey: EncPriKey;
	// pub key of the private key
	PubKey: PubKey;
}

function newUUK(): UUK {
	let uuk: UUK = {
		EncPriKey: {
			Kid: '',
			Enc: '',
			Iv: '',
			Data: '',
			Cty: ''
		},
		EncSymKey: {
			Kid: '',
			Enc: '',
			Iv: '',
			Data: '',
			Cty: '',
			Alg: '',
			P2c: 0,
			P2s: ''
		},
		PubKey: {
			E: '',
			Kid: '',
			Kty: '',
			N: ''
		},
		UUID: '',
		EncryptedBy: ''
	};
	uuk.UUID = crypto.randomUUID();
	return uuk;
}

export async function buildUUK(
	password: Uint8Array,
	mount: Uint8Array,
	secretKey: Uint8Array,
	entityID: Uint8Array
): Promise<UUK> {
	let uuk = newUUK();
	uuk = withInitializationSalt(uuk);
	uuk = withPasswordIterations(uuk, 650000);

	let twoSKD = await twoSkd(uuk, password, mount, secretKey, entityID);

	let symmetricKey: Uint8Array;
	[uuk, symmetricKey] = await withEncSymKey(uuk, twoSKD);

	let pubkey: JsonWebKey;
	[uuk, pubkey] = await withEncPriKey(uuk, symmetricKey);

	uuk = withPubkey(uuk, pubkey);

	return uuk;
}

function withInitializationSalt(uuk: UUK): UUK {
	uuk.EncSymKey.P2s = bytesToHex(crypto.getRandomValues(new Uint8Array(16)));
	return uuk;
}

function withPasswordIterations(uuk: UUK, iterations: number): UUK {
	uuk.EncSymKey.P2c = iterations;
	return uuk;
}

export async function twoSkd(
	uuk: UUK,
	password: Uint8Array,
	mount: Uint8Array,
	secretKey: Uint8Array,
	entityID: Uint8Array
): Promise<Uint8Array> {
	let rawKey = hexToBytes(uuk.EncSymKey.P2s);
	let initialSalt = await crypto.subtle.importKey('raw', rawKey, 'HKDF', false, ['deriveBits']);

	// HKDF 1
	const salt = entityID;
	const info = new TextEncoder().encode('2SKD HKDF 1');
	const hkdf_params = { name: 'HKDF', hash: 'SHA-256', salt, info };
	const saltDerivedKey = await crypto.subtle.deriveBits(hkdf_params, initialSalt, 32 * 8);

	let passwordKey = await crypto.subtle.importKey('raw', password, 'PBKDF2', false, ['deriveBits']);

	let passwordDerivedKey = await crypto.subtle.deriveBits(
		{
			name: 'PBKDF2',
			salt: saltDerivedKey,
			iterations: uuk.EncSymKey.P2c,
			hash: 'SHA-256'
		},
		passwordKey,
		32 * 8
	);

	// hkdf 2
	let secretKeyCrypto = await crypto.subtle.importKey('raw', secretKey, 'HKDF', false, [
		'deriveBits'
	]);

	// HKDF 2
	const hkdf2Salt = mount;
	const hkdf2Info = new TextEncoder().encode('2SKD HKDF 2');
	const hkdf2_params = { name: 'HKDF', hash: 'SHA-256', salt: hkdf2Salt, info: hkdf2Info };
	const secretDerivedKey = await crypto.subtle.deriveBits(hkdf2_params, secretKeyCrypto, 32 * 8);

	let passwordDerivedKeyBytes = new Uint8Array(passwordDerivedKey);
	let secretDerivedKeyBytes = new Uint8Array(secretDerivedKey);

	// XOR passwordDk and secretDk
	let twoSKD = new Uint8Array(32);
	for (let i = 0; i < 32; i++) {
		twoSKD[i] = passwordDerivedKeyBytes[i] ^ secretDerivedKeyBytes[i];
	}

	return twoSKD;
}

async function withEncSymKey(uuk: UUK, twoSKD: Uint8Array): Promise<[UUK, Uint8Array]> {
	let twoSKDKey = await crypto.subtle.importKey('raw', twoSKD, 'AES-GCM', false, ['encrypt']);

	let symmetricKey = crypto.getRandomValues(new Uint8Array(32));

	const iv = crypto.getRandomValues(new Uint8Array(12));

	let encSymKey = await crypto.subtle.encrypt({ name: 'AES-GCM', iv: iv }, twoSKDKey, symmetricKey);

	uuk.EncSymKey.Data = bytesToHex(new Uint8Array(encSymKey));
	uuk.EncSymKey.Iv = bytesToHex(iv);
	uuk.EncSymKey.Enc = 'A256GCM';
	uuk.EncSymKey.Kid = uuk.UUID;
	uuk.EncSymKey.Alg = 'pbkdf2-hkdf';

	return [uuk, symmetricKey];
}

async function withEncPriKey(uuk: UUK, symmetricKey: Uint8Array): Promise<[UUK, JsonWebKey]> {
	let cryptokey = await crypto.subtle.generateKey(
		{
			name: 'RSA-OAEP',
			modulusLength: 4096,
			publicExponent: new Uint8Array([1, 0, 1]),
			hash: 'SHA-256'
		},
		true,
		['encrypt', 'decrypt']
	);
	let key = cryptokey as CryptoKeyPair;
	let prikey = await crypto.subtle.exportKey('jwk', key.privateKey);
	let pubkey = await crypto.subtle.exportKey('jwk', key.publicKey);

	let encKey = await crypto.subtle.importKey('raw', symmetricKey, 'AES-GCM', false, ['encrypt']);

	const iv = crypto.getRandomValues(new Uint8Array(12));
	let encjwk = await crypto.subtle.encrypt(
		{ name: 'AES-GCM', iv: iv },
		encKey,
		new TextEncoder().encode(JSON.stringify(prikey))
	);

	uuk.EncPriKey.Data = bytesToHex(new Uint8Array(encjwk));
	uuk.EncPriKey.Iv = bytesToHex(iv);
	uuk.EncPriKey.Enc = 'A256GCM';
	uuk.EncPriKey.Kid = uuk.UUID;

	return [uuk, pubkey];
}

function withPubkey(uuk: UUK, pubkey: JsonWebKey): UUK {
	uuk.PubKey.E = pubkey.e;
	uuk.PubKey.N = pubkey.n;
	uuk.PubKey.Kid = uuk.UUID;
	uuk.PubKey.Kty = 'RSA';

	return uuk;
}
// Convert a hex string to a byte array
function hexToBytes(hex: string): Uint8Array {
	let bytes = [];
	for (let c = 0; c < hex.length; c += 2) bytes.push(parseInt(hex.substr(c, 2), 16));
	return new Uint8Array(bytes);
}

// Convert a byte array to a hex string
function bytesToHex(bytes: Uint8Array): string {
	let hex = [];
	for (let i = 0; i < bytes.length; i++) {
		let current = bytes[i] < 0 ? bytes[i] + 256 : bytes[i];
		hex.push((current >>> 4).toString(16));
		hex.push((current & 0xf).toString(16));
	}
	return hex.join('');
}
