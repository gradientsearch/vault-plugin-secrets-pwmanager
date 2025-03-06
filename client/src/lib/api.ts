import type { HvEncryptedEntry } from '../routes/unlocked/models/bundle/vault/entry';
import type { BundleSymmetricKey } from '../routes/unlocked/models/bundle/vault/keys';
import type { HvMetadata, BundleMetadata } from '../routes/unlocked/models/bundle/vault/metadata';
import { convertCase, revertCase } from './jsonKey';
import type { newUUK, UUK } from './uuk';

// TODO make method naming convention the same
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

	async delete(path: string) {
		return await fetch(`${this.url}/v1/${path}`, {
			method: 'DELETE',
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

	async getMetadata(b: Bundle): Promise<[HvMetadata | undefined, Error | undefined]> {
		let response = await this.get(`${b.Path}/metadata/entries`);

		if (response.status === 404) {
			// no passwords exist for this vault yet
			// TODO create the pwmanager metadata secret and return the created one
			return [undefined, Error('metadata does not exist')];
		}

		if (response.status != 200) {
			let err = await response.text();
			return [undefined, new Error(`error getting password items metadata: ${err}`)];
		}

		let entriesMetadata = await response.json();
		return [entriesMetadata, undefined];
	}

	async getBundleSymmetricKey(
		b: Bundle,
		entityID: string
	): Promise<[BundleSymmetricKey | undefined, Error | undefined]> {
		let response = await this.get(`${b.Path}/keys/${entityID}`);

		if (response.status === 404) {
			return [undefined, Error('404 not found')];
		}

		if (response.status != 200) {
			let err = await response.text();
			return [undefined, new Error(`error getting vault symmetric key: ${err}`)];
		}

		let vsk = (await response.json()) as BundleSymmetricKey;

		return [vsk, undefined];
	}

	async PutUserKey(b: Bundle, entityID: string, key: string): Promise<Error | undefined> {
		// TODO move this up to bundle service
		let data = {
			data: {
				key: key
			}
		};
		let response = await this.post(`${b.Path}/keys/${entityID}`, JSON.stringify(data));

		if (response.status != 200) {
			let err = await response.text();
			return new Error(`error registering: ${err}`);
		}
		return;
	}

	async PutMetadata(b: Bundle, metadata: any): Promise<Error | undefined> {
		let response = await this.post(`${b.Path}/metadata/entries`, metadata);

		if (response.status != 200) {
			let err = await response.text();
			return new Error(`error registering: ${err}`);
		}
	}

	async DestroyEntry(b: Bundle, id: string): Promise<Error | undefined> {
		let response = await this.delete(`${b.Path}/metadata/entries/${id}`);

		if (response.status != 204) {
			let err = await response.text();
			return new Error(`error deleting entry: ${err}`);
		}
	}

	async PutEntry(b: Bundle, data: any, id: string): Promise<Error | undefined> {
		let response = await this.post(`${b.Path}/entries/${id}`, data);

		if (response.status != 200) {
			let err = await response.text();
			return new Error(`error registering: ${err}`);
		}
	}

	async GetEntry(b: Bundle, id: string ): Promise<[HvEncryptedEntry | undefined, Error | undefined]> {
		let response = await this.get(`${b.Path}/entries/${id}`);

		if (response.status === 404) {
			// no passwords exist for this vault yet
			// TODO create the pwmanager metadata secret and return the created one
			return [undefined, Error('entry does not exist')];
		}

		if (response.status != 200) {
			let err = await response.text();
			return [undefined, new Error(`error getting entry: ${err}`)];
		}

		let hee = await response.json();
		return [hee, undefined];
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
