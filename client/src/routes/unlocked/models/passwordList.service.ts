// responsible for managing what is displayed

import type { PasswordItem } from "./input";

// in the passwordlist column
export interface PasswordListService {
    add(pi: PasswordItem): void
}

export class VaultPasswordListService implements PasswordListService{
    onAddFn: Function
    constructor(onAddFn: Function){
        this.onAddFn = onAddFn
    }

    
    add(pi: PasswordItem): void {
        throw new Error("Method not implemented.");
    }

    
        
}

export class CategoryPasswordList{}
