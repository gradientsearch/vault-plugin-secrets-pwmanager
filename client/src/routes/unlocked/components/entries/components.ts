/**
 * Components functions return the svelte component needed when dynamically
 * displaying a component.
 */
import PasswordInput from './inputs/passwordInput.svelte';
import ItemInput from './inputs/itemInput.svelte';
import Password from './password.svelte';


export function getInputComponent(type: string) {
    switch (type) {
        case 'password':
            return PasswordInput
        case 'text':
            return ItemInput
    }
}

export function getPasswordComponent(type: string) {
    switch (type) {
        case 'password':
            return Password
    }
}