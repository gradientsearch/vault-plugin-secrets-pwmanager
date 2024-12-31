import { writable } from 'svelte/store';
import { newUUK, type UUK } from './uuk';

export function isDevelopment () {
    return document.URL.toString().startsWith('http://localhost')
}

let storedUUKData;
if (isDevelopment()) {
	let uuk = localStorage.getItem('uuk');
	if (uuk) {
		storedUUKData = JSON.parse(uuk) as UUK;
	}
}

// Initialize the store with the retrieved value
export const storedUUK = writable(storedUUKData);
