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

let keyPair: KeyPair = {}
const createStore = (initialState: KeyPair) => {
	const { set, subscribe, update } = writable(initialState);

	return {
		set: (value: KeyPair) => {
			if (isDevelopment()) {
				localStorage.setItem('KeyPair', JSON.stringify(value));
			}
			keyPair = value
		},
		get:() => {
			return keyPair
		},
		subscribe,
		update,
	};
};

export const storedKeyPair = createStore(keyPair);
export const getKeyPair = get(storedKeyPair);

