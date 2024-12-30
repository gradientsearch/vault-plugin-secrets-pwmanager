import { expect, test } from 'vitest';
import { buildUUK, decrypt, decryptEncPriKey, encrypt, twoSkd } from './uuk';

test('buildUUK', async () => {
	const textEncoder = new TextEncoder();
	let password = textEncoder.encode('typingcats');
	let mount = textEncoder.encode('pwmanager');
	let password2 = textEncoder.encode('typingcats');

	

	let secretKey = crypto.getRandomValues(new Uint8Array(16));
	let entityID = textEncoder.encode(crypto.randomUUID());


	
	let uuk = await buildUUK(password, mount, secretKey, entityID);
	let bits = await twoSkd(uuk, password, mount, secretKey, entityID);
	let bits2 = await twoSkd(uuk, password2, mount, secretKey, entityID);

	expect(bits).toEqual(bits2);

	let [prikey, pubkey] = await decryptEncPriKey(uuk, password, mount, secretKey, entityID);
	let encrypted = await encrypt(new TextEncoder().encode('hello cryptography'), pubkey);
	let plaintext = await decrypt(encrypted, prikey);

	expect(plaintext).toEqual('hello cryptography');
});
