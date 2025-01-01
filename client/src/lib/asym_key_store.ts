import { get, writable } from 'svelte/store';
import { newUUK, type UUK } from './uuk';

export function isDevelopment() {
	return document.URL.toString().startsWith('http://localhost');
}

export interface KeyPair {
	PubKey: CryptoKey;
	PriKey: CryptoKey;
}

let storedKeyPairData;
if (isDevelopment()) {
	let KeyPair = localStorage.getItem('KeyPair');
	if (KeyPair) {
		storedKeyPairData = JSON.parse(KeyPair) as KeyPair;
	}
}

let keyPair: KeyPair = {};
const createStore = (initialState: KeyPair) => {
	const { set, subscribe, update } = writable(initialState);

	return {
		set: async (value: KeyPair) => {
			if (isDevelopment()) {
				let priKey = await crypto.subtle.exportKey('jwk', value.PriKey);
				let pubKey = await crypto.subtle.exportKey('jwk', value.PubKey);
				localStorage.setItem('priKey', JSON.stringify(priKey));
				localStorage.setItem('pubKey', JSON.stringify(pubKey));
			}
			keyPair = value;
		},
		get: async () => {
			if (isDevelopment()) {
				let priKeyJson = localStorage.getItem('priKey');
				let pubKeyJson = localStorage.getItem('pubKey');

				if (priKeyJson === null || pubKeyJson === null) {
					return;
				}

				let pirKey = await crypto.subtle.importKey(
					'jwk',
					JSON.parse(priKeyJson),
					{
						name: 'RSA-OAEP',
						hash: 'SHA-256'
					},
					false,
					['encrypt']
				);

				let pubKey = await crypto.subtle.importKey(
					'jwk',
					JSON.parse(pubKeyJson),
					{
						name: 'RSA-OAEP',
						hash: 'SHA-256'
					},
					false,
					['encrypt']
				);

				let kp: KeyPair = {
					PriKey: pirKey,
					PubKey: pubKey
				};

				return kp;
			}
			return keyPair;
		},
		subscribe,
		update
	};
};

export const storedKeyPair = createStore(keyPair);
export const getKeyPair = get(storedKeyPair);
