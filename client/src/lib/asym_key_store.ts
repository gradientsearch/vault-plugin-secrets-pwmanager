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

const createStore = (initialState: KeyPair | undefined) => {
	const { set, subscribe, update } = writable(initialState);

	return {
		set,
		subscribe,
		update,
		set value(newValue: KeyPair) {
			if (isDevelopment()) {
				localStorage.setItem('KeyPair', JSON.stringify(newValue));
			}
			set(newValue);
		},
		get value() {
			let value: KeyPair|undefined;
			subscribe((v) => (value = v))();
			return value;
		}
	};
};

export const storedKeyPair = createStore(undefined);
export const getKeyPair = get(storedKeyPair);

