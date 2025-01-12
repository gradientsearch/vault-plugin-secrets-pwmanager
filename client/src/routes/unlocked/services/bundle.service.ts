import {
	pubkeyEncrypt,
	exportJwkKey,
	generateSymmetricKey,
	prikeyDecrypt,
	hexToBytes,
	symmetricEncrypt,
	importJWKkey
} from '$lib/helper';
import type { Metadata } from '../models/bundle/vault/metadata';
import type { Entry } from '../models/entry';
import type { Zarf } from '../models/zarf';
import { userService } from './user.service';

/**
 * Interface the BundleView uses to interface with different types of bundles e.g a Vault bundle or
 * a Category bundle.
 */
export interface BundleService {
	addEntry(pi: Entry): Promise<Error | undefined>;
	getEntries(): Promise<[Entry[] | undefined, Error | undefined]>;
	init(): Promise<Error | undefined>;
}

/**
 * A VaultBundleService is responsible for interfacing with a HashiCorp Vault KV2 secret mount. w.r.t.
 * pwmanager a a KV2 secret mount is a vault. A vault has the following path convention
 * `vaults/{{ identity.entity.id }}/<vault name>
 *
 * The KV2 secret mount contains the following paths:
 * - keys/{{ identity.entity.id }}: each vault has a symmetric key used to encrypt all secrets. That
 * symmetric keys is encrypted with users public key
 * - metadata: vault metadata for each entry
 * - entries/<entry name>: user entires
 */
export class VaultBundleService implements BundleService {
	onAddFn: Function;
	zarf: Zarf;
	bundle: Bundle;
	symmetricKey: CryptoKey | undefined;

	constructor(zarf: Zarf, bundle: Bundle, onAddFn: Function) {
		this.zarf = zarf;
		this.bundle = bundle;
		this.onAddFn = onAddFn;
	}

	async init(): Promise<Error | undefined> {
		// get decryption key for this vault
		//whats the entity id

		let entityID = userService.getEntityID();
		let [key, err] = await this.zarf.Api.getVaultSymmetricKey(this.bundle, entityID);
		if (err?.toString().includes('404 not found')) {
			let [key, err] = await this.createVaultEncryptionKey(entityID);

			if (err !== undefined) {
				return err;
			}

			err = await this.createVaultMetadata();
			if (err != undefined) {
				return err;
			}
		}

		// TODO this is a little messy. Let's try and clean this up.
		if (key !== undefined) {
			let encryptedSymmetricKey = key.data.data.key;
			let jwk = await prikeyDecrypt(encryptedSymmetricKey, this.zarf.Keypair.PriKey);
			let jwkObj = JSON.parse(jwk);
			let ck = await importJWKkey(jwkObj);
			this.symmetricKey = ck;
		} else if (this.symmetricKey === undefined){
			return Error('error retrieving vault symmetric key')
		}
	}

	async createVaultEncryptionKey(
		entityID: string
	): Promise<[CryptoKey | undefined, Error | undefined]> {
		// TODO make this a version CAS version 0 only operation. Never want to overwrite a vault symmetric key
		let key = await generateSymmetricKey();
		let jwk = await exportJwkKey(key);
		let encrypted = await pubkeyEncrypt(
			new TextEncoder().encode(JSON.stringify(jwk)),
			this.zarf.Keypair.PubKey
		);
		let err = await this.zarf.Api.PutUserKey(this.bundle, entityID, encrypted);
		if (err !== undefined) {
			console.log('err: ', err);
			return [undefined, Error('error retrieving vault symmetric key : ', err)];
		}
		this.symmetricKey = key;
		return [key, undefined];
	}

	async createVaultMetadata(): Promise<Error | undefined> {
		let metadata: Metadata[] = [];
		let [data, err] = await this.encryptPayload(metadata);
		if (err !== undefined) {
			return Error('error encrypted vault metadata');
		}

		err = await this.zarf.Api.PutMetadata(this.bundle, data);
		if (err !== undefined) {
			return Error('error creating encrypted vault metadata');
		}
	}

	async encryptPayload(payload: any): Promise<[string | undefined, Error | undefined]> {
		// TODO refactor this class so this isn't necessary
		if (this.symmetricKey === undefined) {
			return [undefined, Error('vault encryption key is undefined')];
		}

		let [iv, encrypted] = await symmetricEncrypt(
			new TextEncoder().encode(JSON.stringify(payload)),
			this.symmetricKey
		);

		let data = {
			data: {
				entry: encrypted,
				iv: iv
			}
		};
		return [JSON.stringify(data), undefined];
	}

	async getEntries(): Promise<[Entry[] | undefined, Error | undefined]> {
		let [entries, err] = await this.zarf.Api.getVaultEntries(this.bundle);
		if (err !== undefined) {
			return [undefined, err];
		}
		return [entries, undefined];
	}

	async addEntry(e: Entry): Promise<Error | undefined> {
		//store data in vault
		// encrypt
		let entryName = crypto.randomUUID();

		let [data, err] = await this.encryptPayload(e);
		if (err != undefined) {
			return err;
		}

		this.zarf.Api.PutEntry(this.bundle, data, 'metadata');

		// update metadata
		// if error delete entry

		this.onAddFn([e]);
		return new Promise((resolve) => {
			console.log('slice add password');
			resolve(undefined);
		});
	}
}

export class CategoryBundleService {}
