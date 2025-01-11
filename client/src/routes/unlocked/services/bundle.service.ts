// responsible for interacting with Vault

import type { Entry } from '../models/entry';
import type { Zarf } from '../models/zarf';
import { userService } from './user.service';

export interface BundleService {
	addEntry(pi: Entry): Promise<Error | undefined>;
	getEntries(): Promise<Entry[]>;
	init(): Promise<any	>;
}

// Vault is a HashiCorp Vault KV2 secret mount
export class VaultBundleService implements BundleService {
	
	onAddFn: Function;
	zarf: Zarf | undefined;
	bundle: Bundle;

	constructor(zarf: Zarf | undefined, bundle: Bundle, onAddFn: Function) {
		this.zarf = zarf;
		this.bundle = bundle;
		this.onAddFn = onAddFn;
	}
	
	async init() {
		// get decryption key for this vault
		//whats the entity id 

		let entityID = userService.getEntityID()
		let key, err =  await this.zarf?.Api?.getVaultKey(this.bundle, entityID)
		if (err?.toString().includes('404 not found')) {
			console.log('creating vault symmetric encryption key')
		}
	}
	
	async getEntries(): Promise<Entry[]> {
		let entries, err = await this.zarf?.Api?.getVaultMetadata(this.bundle);
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
