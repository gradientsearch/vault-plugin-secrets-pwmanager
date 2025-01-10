import type { Component } from 'svelte';

import PasswordInput from './inputs/passwordInput.svelte';
import ItemInput from './inputs/itemInput.svelte';


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
            return 
    }
}