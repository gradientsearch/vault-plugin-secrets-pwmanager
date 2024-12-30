import { convertCase } from "./jsonKey";
import type { UUK } from "./uuk";

export class Api {
	vaultToken: string;
	url: string;
	mount: string;

	constructor(vaultToken: string, url: string, mount: string) {
		this.vaultToken = vaultToken;
		this.url = url;
		this.mount = mount;
	}

	async post(path: string, content: any) {
		return await fetch(`${this.url}/${path}`, {
			method: 'POST',
			body: content,
			headers: {
				'Content-Type': 'application/json; charset=UTF-8',
				'X-Vault-Token': this.vaultToken,
				'X-Vault-Request': 'true'
			}
		});
	}

	async get(path: string) {
		return await fetch(`${this.url}/${path}`, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json; charset=UTF-8',
				'X-Vault-Request': 'true',
				'X-Vault-Token': this.vaultToken
			}
		});
	}

	async tokenLookup() {
		let response = await this.get('auth/token/lookup-self');

		if (!response.ok) {
			return;
		}

		// If you care about a response:
		if (response.body !== null) {
			return response.json();
		}
	}

	async register(uuk: UUK): Promise<Error|undefined> {
        let data = JSON.stringify(convertCase(JSON.parse(JSON.stringify(uuk))))
		let response = await this.post(`${this.mount}/register`, data);

        if (response.status != 204){
            let err  = await response.text()
            return new Error(`error registers: ${err}`)
        }
	}
}