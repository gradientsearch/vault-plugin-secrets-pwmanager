import { expect, test } from 'vitest';
import { convertCase, type JSONObject } from './jsonKey';
import { type UUK } from './uuk';

test('convertCase', () => {
	let uuk: UUK = {
		EncPriKey: {
			Kid: 'kid',
			Enc: 'enc',
			Iv: 'iv',
			Data: 'data',
			Cty: 'cty'
		},
		EncryptedBy: 'mp'
	};
	let uukJson = JSON.parse(JSON.stringify(uuk));
	let snakeCaseJson = JSON.stringify(convertCase(uukJson, true));
	expect(snakeCaseJson).toContain('enc_pri_key');

	let json = '{"users":{"users":[{"EntityName":"stephen","EntityID":"","Capabilities":"","IsAdmin":false},{"EntityName":"stephen","EntityID":"","Capabilities":"","IsAdmin":false},{"EntityName":"stephen","EntityID":"","Capabilities":"","IsAdmin":false},{"EntityName":"stephen","EntityID":"","Capabilities":"","IsAdmin":false},{"EntityName":"stephen","EntityID":"","Capabilities":"","IsAdmin":false},{"EntityName":"stephen","EntityID":"","Capabilities":"","IsAdmin":false}]}}'
	
	let obj = convertCase(JSON.parse(json), true)
	console.log('print')
	

});
