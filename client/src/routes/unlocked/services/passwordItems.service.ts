// responsible for interacting with Vault

import type { PasswordItem } from '../models/input';
import type { Zarf } from '../models/zarf';

// in the passwordlist column
export interface PasswordItemsService {
	add(pi: PasswordItem): Promise<Error | undefined>;
	get(): Promise<PasswordItem[]>;

}

export class VaultPasswordItemsService implements PasswordItemsService {
	onAddFn: Function;
	zarf: Zarf | undefined;
	selectedVault: PasswordItems;

	constructor(zarf: Zarf | undefined, selectedVault: PasswordItems, onAddFn: Function) {
		this.zarf = zarf;
		this.selectedVault = selectedVault;
		this.onAddFn = onAddFn;
	}


	async init() {
		// get decryption key for this vault
		//whats the entity id 
		
		//await this.zarf?.Api?.get()
	}
	
	async get(): Promise<PasswordItem[]> {
		let passwordItems, err = await this.zarf?.Api?.getPasswordListMetadata(this.selectedVault);
		return [];
	}

	async add(pi: PasswordItem): Promise<Error | undefined> {
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
