import type { Component } from 'svelte';

import PasswordInput from './components/passwordInput.svelte';
import ItemInput from './components/itemInput.svelte';


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