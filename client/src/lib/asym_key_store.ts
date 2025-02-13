import { get, writable } from 'svelte/store';
import { newUUK, type UUK } from './uuk';

export function isDevelopment() {
	return document.URL.toString().startsWith('http://localhost');
}

export interface KeyPair {
	PubKey: CryptoKey | undefined;
	PriKey: CryptoKey | undefined;
}

let storedKeyPairData;
if (isDevelopment()) {
	let KeyPair = localStorage.getItem('KeyPair');
	if (KeyPair) {
		storedKeyPairData = JSON.parse(KeyPair) as KeyPair;
	}
}

let keyPair: KeyPair = {
	PubKey: undefined,
	PriKey: undefined
};
const createStore = (initialState: KeyPair) => {
	const { set, subscribe, update } = writable(initialState);

	return {
		set: async (value: KeyPair) => {
			if (isDevelopment()) {
				if (value.PriKey !== undefined && value.PubKey !== undefined) {
					let priKey = await crypto.subtle.exportKey('jwk', value.PriKey);
					let pubKey = await crypto.subtle.exportKey('jwk', value.PubKey);

					localStorage.setItem('priKey', JSON.stringify(priKey));
					localStorage.setItem('pubKey', JSON.stringify(pubKey));
				}
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
					['decrypt']
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
				keyPair = kp;
				return kp;
			}

			return keyPair;
		},
		subscribe,
		update
	};
};

export const storedKeyPair = createStore(keyPair);
