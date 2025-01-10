export interface Item {
	Name: string;
	Type: string;
	Metadata: Metadata;
}

export interface Metadata {
	Name: string;
	Type: string;
	Value: string;
	Path: string;
}

export interface Input {
	Type: string;
	Label: string;
	Placeholder: string;
	Value: string;
}

export interface Core {
	Items: Input[];
	Order: string[];
}

export interface Section {
	Name: string;
	Items: Input[];
}

export interface More {
	Items: Input[];
	Order: Section[];
}

export interface PasswordItem extends Item {
	Metadata: Metadata;
	Core: Core;
	More: More;
	Tags: string[];
}

export function newPasswordItem(): PasswordItem {
	return {
		Tags: [],
		Name: '',
		Type: '',
		Metadata: {
			Name: 'Password',
			Type: 'password',
			Value: '',
			Path: ''
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
			],
			Order: ['1']
		},
		More: {
			Items: [],
			Order: []
		}
	};
}

export class PasswordInput implements Input {
	Type: string = 'password';
	Label: string = 'password';
	Placeholder: string = 'password';
	Value: string = '';
	Metadata: any;
}

export function newPasswordInput(): PasswordInput {
	return {
		Type: '',
		Label: '',
		Placeholder: '',
		Value: '',
		Metadata: undefined
	};
}
