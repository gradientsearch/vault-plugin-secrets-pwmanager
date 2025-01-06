import type { Component } from 'svelte';
import PasswordInput from './components/passwordInput.svelte';
import ItemInput from './components/itemInput.svelte';


export function getComponent(type: string) {
    switch (type) {
        case 'password':
            return PasswordInput
        case 'text':
            return ItemInput
    }
}