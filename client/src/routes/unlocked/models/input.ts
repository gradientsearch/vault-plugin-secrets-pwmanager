export interface Item {
	Name: string;
	Type: string;
	Metadata: any;
}

export interface Input {
	Type: string;
	Label: string;
	Placeholder: string;
	Value: string;
	Metadata: any;
}

export interface Core {
	Items: Map<string, Input>;
	Order: string[];
}

export interface Section {
	Name: string;
	Items: Map<string, Input>;
}

export interface More {
	Items: Map<string, Input>;
	Order: Section[];
}

export interface PasswordItem extends Item {
	Core: Core;
	More: More;
	Tags: string[];
}

function newPasswordItem(): PasswordItem {
	return {
		Tags: [],
		Name: '',
		Type: '',
		Metadata: undefined,
		Core: {
			Items: new Map<string, Input>(),
			Order: []
		},
		More: {
			Items: new Map<string, Input>(),
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

export function newPasswordInput(): PasswordInput{
    return {
        Type: "",
        Label: "",
        Placeholder: "",
        Value: "",
        Metadata: undefined
    }
}
