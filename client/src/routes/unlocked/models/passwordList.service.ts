// responsible for managing what is displayed

import type { Api } from '$lib/api';
import type { PasswordItem } from './input';
import type { Zarf } from './zarf';

// in the passwordlist column
export interface PasswordListService {
	add(pi: PasswordItem):  Promise<Error | undefined> ;
    get(): PasswordItem[]
}

export class VaultPasswordListService implements PasswordListService {
	onAddFn: Function;
	zarf: Zarf | undefined;
    selectedVault: PasswordList

	constructor(zarf: Zarf | undefined, selectedVault: PasswordList, onAddFn: Function) {
        
        this.zarf = zarf
        this.selectedVault = selectedVault
		this.onAddFn = onAddFn;
	}
    get(): PasswordItem[] {
        this.zarf?.Api?.getPasswordListMetadata(this.selectedVault)
        return []
    }

	async add(pi: PasswordItem): Promise<Error | undefined> {
        //store data in vault
        // update metadata
        this.onAddFn([pi])
		return new Promise((resolve) => {
            console.log('slice add password')
            resolve(undefined);
        });
	}
}

export class CategoryPasswordList {}
