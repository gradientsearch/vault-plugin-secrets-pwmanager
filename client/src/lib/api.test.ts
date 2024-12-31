import { expect, test } from 'vitest';
import { Api } from './api';
import { buildUUK, bytesToHex } from './uuk';
import { VITE_VAULT_TOKEN } from '$env/static/private';

test('api', async () => {
	//test.skip()
	let api = new Api(VITE_VAULT_TOKEN, 'http://localhost:8200', 'pwmanager');
	let json = await api.tokenLookup();

	const textEncoder = new TextEncoder();
	let password = textEncoder.encode('typingcats');
	let mount = textEncoder.encode('pwmanager');

	let secretKey = crypto.getRandomValues(new Uint8Array(16));

	let entityID = json != undefined ? json['data']['entity_id'] : crypto.randomUUID();

	let uuk = await buildUUK(password, mount, secretKey, new TextEncoder().encode(entityID));

	let err = await api.register(uuk);
});