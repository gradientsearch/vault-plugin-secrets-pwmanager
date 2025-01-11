import { bytesToHex, hexToBytes } from './uuk';

export async function generateSymmetricKey() {
	return await crypto.subtle.generateKey({ name: 'AES-GCM', length: 256 }, true, [
		'encrypt',
		'decrypt'
	]);
}

export async function exportJwkKey(key: CryptoKey): Promise<string> {
	const exportedKey = await crypto.subtle.exportKey('jwk', key);
    let json = JSON.stringify(exportedKey)
	return bytesToHex(new TextEncoder().encode(json))
}

export async function symmetricEncrypt(
	payload: Uint8Array,
	symmetricKey: CryptoKey
): Promise<string> {
	const iv = window.crypto.getRandomValues(new Uint8Array(12));
	const encrypted = await crypto.subtle.encrypt({ name: 'AES-GCM', iv: iv }, symmetricKey, payload);

	return bytesToHex(new Uint8Array(encrypted));
}

export async function symmetricDecrypt(
	payload: string,
	iv: Uint8Array,
	symmetricKey: CryptoKey
): Promise<string> {
	const plaintext = await crypto.subtle.decrypt(
		{ name: 'AES-GCM', iv: iv },
		symmetricKey,
		hexToBytes(payload)
	);
	return new TextDecoder().decode(plaintext);
}
