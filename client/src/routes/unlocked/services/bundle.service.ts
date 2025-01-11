import { pubkeyEncrypt, exportJwkKey, generateSymmetricKey, prikeyDecrypt, hexToBytes } from '$lib/helper';
import type { Entry } from '../models/entry';
import type { Zarf } from '../models/zarf';
import { userService } from './user.service';

/**
 * Interface the BundleView uses to interface with different types of bundles e.g a Vault bundle or
 * a Category bundle.
 */
export interface BundleService {
	addEntry(pi: Entry): Promise<Error | undefined>;
	getEntries(): Promise<Entry[]>;
	init(): Promise<Error | undefined>;
}

/**
 * A VaultBundleService is responsible for interfacing with a HashiCorp Vault KV2 secret mount. w.r.t.
 * pwmanager a a KV2 secret mount is a vault. A vault has the following path convention
 * `vaults/{{ identity.entity.id }}/<vault name>
 *
 * The KV2 secret mount contains the following paths:
 * - keys/{{ identity.entity.id }}: each vault has a symmetric key used to encrypt all secrets. That
 * symmetric keys is encrypted with users public key
 * - metadata: vault metadata for each entry
 * - entries/<entry name>: user entires
 */
export class VaultBundleService implements BundleService {
	onAddFn: Function;
	zarf: Zarf;
	bundle: Bundle;
	decryptionKey: CryptoKey | undefined;

	constructor(zarf: Zarf, bundle: Bundle, onAddFn: Function) {
		this.zarf = zarf;
		this.bundle = bundle;
		this.onAddFn = onAddFn;
	}

	async init(): Promise<Error | undefined> {
		// get decryption key for this vault
		//whats the entity id

		let entityID = userService.getEntityID();
		let [key, err] = await this.zarf.Api.getVaultSymmetricKey(this.bundle, entityID);
		if (err?.toString().includes('404 not found')) {
			await this.createVaultEncryptionKey(entityID)
		}

		if (key === undefined) {
			return Error('error retrieving vault symmetric key:(');
		}

		let encryptedSymmetricKey = key.data.data.key

		console.log('jwk encrypted is : ', new TextDecoder().decode(hexToBytes(encryptedSymmetricKey)))
		let jwk = await prikeyDecrypt(encryptedSymmetricKey, this.zarf.Keypair.PriKey)
		console.log('jwk is: ', jwk)
		return

	}

	async createVaultEncryptionKey(entityID: string) {
		// TODO make this a version CAS version 0 only operation. Never want to overwrite a vault encryption key
		let key = await generateSymmetricKey();
		let jwk = await exportJwkKey(key);
		let encrypted = await pubkeyEncrypt(
			new TextEncoder().encode(JSON.stringify(jwk)),
			this.zarf.Keypair.PubKey
		);
		let err = await this.zarf.Api.PutUserKey(this.bundle, entityID, encrypted);
		if (err !== undefined) {
			return Error('error retrieving vault symmetric key :(');
		}
		this.decryptionKey = key;
	}

	async getEntries(): Promise<Entry[]> {
		let entries,
			err = await this.zarf.Api.getVaultMetadata(this.bundle);
		return [];
	}

	async addEntry(pi: Entry): Promise<Error | undefined> {
		//store data in vault
		// update metadata
		this.onAddFn([pi]);
		return new Promise((resolve) => {
			console.log('slice add password');
			resolve(undefined);
		});
	}
}

export class CategoryBundleService {}
