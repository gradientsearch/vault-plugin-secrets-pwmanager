import { exportJwkKey, generateSymmetricKey } from '$lib/helper';
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
	init(): Promise<any>;
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

	constructor(zarf: Zarf, bundle: Bundle, onAddFn: Function) {
		this.zarf = zarf;
		this.bundle = bundle;
		this.onAddFn = onAddFn;
	}

	async init() {
		// get decryption key for this vault
		//whats the entity id

		let entityID = userService.getEntityID();
		let key,
			err = await this.zarf.Api.getVaultKey(this.bundle, entityID);
		if (err?.toString().includes('404 not found')) {
			key = await generateSymmetricKey();
			let exportedKey = await exportJwkKey(key);
			let encryptedVaultKey = 
			this.zarf.Api.PutUserKey(this.bundle, entityID,  exportedKey);
			console.log('creating vault symmetric encryption key');
		}
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
