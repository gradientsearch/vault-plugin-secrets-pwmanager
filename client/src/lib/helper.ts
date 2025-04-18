export async function generateSymmetricKey() {
	return await crypto.subtle.generateKey({ name: 'AES-GCM', length: 256 }, true, [
		'encrypt',
		'decrypt'
	]);
}

export async function exportJwkKey(key: CryptoKey): Promise<JsonWebKey> {
	return await crypto.subtle.exportKey('jwk', key);
}

export async function importJWKkey(jwk: JsonWebKey): Promise<CryptoKey> {
	return await crypto.subtle.importKey('jwk', jwk, 'AES-GCM', true, ['encrypt', 'decrypt']);
}

export async function symmetricEncrypt(
	payload: Uint8Array,
	symmetricKey: CryptoKey
): Promise<[string, string]> {
	const iv = crypto.getRandomValues(new Uint8Array(12));
	const encrypted = await crypto.subtle.encrypt({ name: 'AES-GCM', iv: iv }, symmetricKey, payload);

	return [bytesToHex(iv), bytesToHex(new Uint8Array(encrypted))];
}

// TODO should this return an error?
export async function symmetricDecrypt(
	payload: string,
	iv: string,
	symmetricKey: CryptoKey
): Promise<string> {
	const plaintext = await crypto.subtle.decrypt(
		{ name: 'AES-GCM', iv: hexToBytes(iv) },
		symmetricKey,
		hexToBytes(payload)
	);
	return new TextDecoder().decode(plaintext);
}

export async function pubkeyEncrypt(payload: Uint8Array, pubkey: CryptoKey): Promise<string> {
	let encrypted = await crypto.subtle.encrypt({ name: 'RSA-OAEP' }, pubkey, payload);
	return bytesToHex(new Uint8Array(encrypted));
}

export async function prikeyDecrypt(payload: string, prikey: CryptoKey): Promise<string> {
	let plaintext = await crypto.subtle.decrypt({ name: 'RSA-OAEP' }, prikey, hexToBytes(payload));
	return new TextDecoder().decode(plaintext);
}

// Convert a hex string to a byte array
export function hexToBytes(hex: string): Uint8Array {
	let bytes = [];
	for (let c = 0; c < hex.length; c += 2) bytes.push(parseInt(hex.substr(c, 2), 16));
	return new Uint8Array(bytes);
}

// Convert a byte array to a hex string
export function bytesToHex(bytes: Uint8Array): string {
	let hex = [];
	for (let i = 0; i < bytes.length; i++) {
		let current = bytes[i] < 0 ? bytes[i] + 256 : bytes[i];
		hex.push((current >>> 4).toString(16));
		hex.push((current & 0xf).toString(16));
	}
	return hex.join('');
}

export function getSecureRandomString(length: number): string {
	const charset =
		'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}<>?/';
	const charsetLength = charset.length;
	const result: string[] = [];

	const randomValues = new Uint32Array(length);
	crypto.getRandomValues(randomValues);

	for (let i = 0; i < length; i++) {
		const index = randomValues[i] % charsetLength;
		result.push(charset.charAt(index));
	}

	return result.join('');
}
