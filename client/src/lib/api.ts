import { convertCase, revertCase } from './jsonKey';
import type { newUUK, UUK } from './uuk';

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
		return await fetch(`${this.url}/v1/${path}`, {
			method: 'POST',
			body: content,
			headers: {
				'Content-Type': 'application/json; charset=UTF-8',
				'X-Vault-Token': this.vaultToken
			}
		});
	}

	async get(path: string) {
		return await fetch(`${this.url}/v1/${path}`, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json; charset=UTF-8',
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

	async register(uuk: UUK): Promise<Error | undefined> {
		let data = convertCase(uuk, true);
		let response = await this.post(`${this.mount}/register`, data);

		if (response.status != 204) {
			let err = await response.text();
			return new Error(`error registering: ${err}`);
		}
	}

	async uuk(entityID: string): Promise<[UUK | undefined, Error | undefined]> {
		let response = await this.get(`${this.mount}/users/${entityID}`);

		if (response.status != 200) {
			let err = await response.text();
			return [undefined, new Error(`error registering: ${err}`)];
		}

		let json = await response.json();
		let uuk = revertCase<UUK>(json['data']['uuk'], false) as UUK;

		return [uuk, undefined];
	}

	async getVaultMetadata(pl: Bundle) {
		let response = await this.get(`${pl.Path}/metadata`);

		if (response.status === 404) {
			// no passwords exist for this vault yet
			// TODO create the pwmanagerMetadata secret and return the created one
			return [[], 404];
		}

		if (response.status != 200) {
			let err = await response.text();
			return [undefined, new Error(`error getting password items metadata: ${err}`)];
		}

		let json = await response.json();

		return {};
	}

	async getVaultSymmetricKey(b: Bundle, entityID: string): Promise<[VaultSymmetricKey | undefined,  Error | undefined]> {
		let response = await this.get(`${b.Path}/data/keys/${entityID}`);

		if (response.status === 404) {
			return [undefined, Error('404 not found')];
		}

		if (response.status != 200) {
			let err = await response.text();
			return [undefined, new Error(`error getting vault symmetric key: ${err}`)];
		}

		let vsk = await response.json() as VaultSymmetricKey
		

		return [vsk, undefined];
	}

	async PutUserKey(b: Bundle, entityID: string, key: string): Promise<Error|undefined>{
		let data = {
			data: {
				key: key
			}
		}
		let response = await this.post(`${b.Path}/data/keys/${entityID}`, JSON.stringify(data));

		if (response.status != 204) {
			let err = await response.text();
			return new Error(`error registering: ${err}`);
		}
	}

	async PutEntry(b: Bundle, data: any, metadata: any): Promise<Error|undefined>{
		
		

		
		let response = await this.post(`${b.Path}/data/entries/`, JSON.stringify(data));

		if (response.status != 204) {
			let err = await response.text();
			return new Error(`error registering: ${err}`);
		}
	}
}

export function getAPI() {
	let info = localStorage.getItem('loginInfo');
	if (info !== null) {
		//TODO check for null
		let infoObj = JSON.parse(info);
		return new Api(infoObj['token'], infoObj['url'], infoObj['mount']);
	}
}
