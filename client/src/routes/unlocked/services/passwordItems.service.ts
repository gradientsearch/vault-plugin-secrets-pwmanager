// responsible for interacting with Vault

import type { Entry } from '../models/input';
import type { Zarf } from '../models/zarf';
import { userService } from './user.service';

// in the passwordlist column
export interface PasswordItemsService {
	add(pi: Entry): Promise<Error | undefined>;
	get(): Promise<Entry[]>;

}

export class VaultPasswordItemsService implements PasswordItemsService {
	
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
		await this.zarf?.Api?.getVaultKey(entityID)
		this.bundle.Path
		
		//await this.zarf?.Api?.get()
	}
	
	async get(): Promise<Entry[]> {
		let passwordItems, err = await this.zarf?.Api?.getVaultMetadata(this.bundle);
		return [];
	}

	async add(pi: Entry): Promise<Error | undefined> {
		//store data in vault
		// update metadata
		this.onAddFn([pi]);
		return new Promise((resolve) => {
			console.log('slice add password');
			resolve(undefined);
		});
	}
}

export class CategoryItemsList {}
