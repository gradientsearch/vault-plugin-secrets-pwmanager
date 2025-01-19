
/**
 * An Input stores the HTML input attributes.
 */
export interface Input {
	Type: string;
	Label: string;
	Placeholder: string;
	Value: string;
}

/**
 * A PasswordInput sets the default for a password input.
 */
export class PasswordInput implements Input {
	Type: string = 'password';
	Label: string = 'password';
	Placeholder: string = 'password';
	Value: string = '';
	Metadata: any;
}

export class TextInput implements Input {
	Type: string = 'text';
	Label: string = 'text';
	Placeholder: string = 'text';
	Value: string = '';
	Metadata: any;
}

export class DateInput implements Input {
	Type: string = 'date';
	Label: string = 'date';
	Placeholder: string = '';
	Value: string = '';
	Metadata: any;
}
