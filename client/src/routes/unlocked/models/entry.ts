import type { Input } from './input';

/**
 * Metadata is the info stored in the metadata path. This is encrypted data that's
 * not considered secret. For example, a username would be a Value but a users password
 * would not be.
 */
export interface Metadata {
	Name: string;
	Type: string;
	Value: string;
	ID: string; 
}

/**
 * Each Entry has a Core set of inputs. For example, a Login has as username, password, and link. A
 * Password just has a username and and password
 */
export interface Core {
	Items: Input[];
}

/**
 * A More holds the additional inputs a user may want to associate with the entry. Security questions,
 * birthdays etc...
 */
export interface More {
	Items: Input[];
}

/**
 * An Entry is the generic concept of `password`. Entries can be logins, passwords, secure note, ssh keys etc...
 */
export interface Entry {
	Name: string;
	Type: string;
	Metadata: Metadata;
	Core: Core;
	More: More;
	Tags: string[];
}

export enum MODE {
    VIEW,
    EDIT,
}


export function newPasswordEntry(): Entry {
	return {
		Tags: [],
		Name: 'Password',
		Type: 'password',
		Metadata: {
			Name: 'Password',
			Type: 'password',
			Value: '',
			ID: ''
		},
		Core: {
			Items: [
				{
					Type: 'text',
					Label: 'username',
					Placeholder: 'username',
					Value: ''
				},
				{
					Type: 'password',
					Label: 'Password',
					Placeholder: 'Password',
					Value: ''
				}
			]
		},
		More: {
			Items: []
		}
	};
}
