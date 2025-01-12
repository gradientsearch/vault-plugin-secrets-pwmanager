import {
	pubkeyEncrypt,
	exportJwkKey,
	generateSymmetricKey,
	prikeyDecrypt,
	hexToBytes,
	symmetricEncrypt,
	importJWKkey,
	symmetricDecrypt,
	bytesToHex
} from '$lib/helper';
import type { EncryptedEntry } from '../models/bundle/vault/entry';
import type { VaultMetadata } from '../models/bundle/vault/metadata';
import type { Entry } from '../models/entry';
import type { Zarf } from '../models/zarf';
import { userService } from './user.service';

/**
 * Interface the BundleView uses to interface with different types of bundles e.g a Vault bundle or
 * a Category bundle.
 */
export interface BundleService {
	addEntry(pi: Entry): Promise<Error | undefined>;
	getMetadata(): Promise<[VaultMetadata | undefined, Error | undefined]>;
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
		} else if (this.symmetricKey === undefined) {
			return Error('error retrieving vault symmetric key');
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
		let metadata: VaultMetadata = {
			entries: []
		};
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

		let data: EncryptedEntry = {
			data: {
				entry: encrypted,
				iv: iv
			}
		};
		return [JSON.stringify(data), undefined];
	}

	async decryptPayload(ed: EncryptedEntry): Promise<[string | undefined, Error | undefined]> {
		// TODO refactor this class so this isn't necessary
		if (this.symmetricKey === undefined) {
			return [undefined, Error('vault encryption key is undefined')];
		}

		let decrypted = await symmetricDecrypt(ed.data.entry, ed.data.iv, this.symmetricKey);

		return [decrypted, undefined];
	}

	async getMetadata(): Promise<[VaultMetadata | undefined, Error | undefined]> {
		let [md, err] = await this.zarf.Api.getMetadata(this.bundle);
		if (err !== undefined) {
			return [undefined, err];
		}

		if (md === undefined) {
			return [undefined, Error('no data returned from server')];
		}

		let [plaintext, err2] = await this.decryptPayload(md.data);
		if (err2 !== undefined) {
			return [undefined, err2];
		}

		if (plaintext === undefined) {
			return [undefined, Error('decrypted metadata entry undefined')];
		}
		let vm = JSON.parse(plaintext) as VaultMetadata;

		return [vm, undefined];
	}

	async addEntry(e: Entry): Promise<Error | undefined> {
		//store data in vault
		// encrypt
		let entryName = crypto.randomUUID();
		e.Metadata.ID = entryName;

		let [data, err2] = await this.encryptPayload(e);
		if (err2 != undefined) {
			return err2;
		}

		let err3 = await this.zarf.Api.PutEntry(this.bundle, data, e.Metadata.ID);
		if (err3 !== undefined) {
			return Error(`error putting entry:  ${err3.message}`);
		}

		let [metadata, err] = await this.getMetadata();
		if (err !== undefined) {
			return Error('error retrieving latest vault metadata');
		}

		console.log(metadata)

		metadata?.entries.push(e.Metadata);
		
		let ep = await this.encryptPayload(metadata)

		// TODO add CAS version from
		let err4 = await this.zarf.Api.PutMetadata(this.bundle, ep);
		if (err4 !== undefined) {
			return Error('error putting metadata: ', err4);
		}

		// TODO if error delete entry

		this.onAddFn([metadata]);
		return new Promise((resolve) => {
			console.log('slice add password');
			resolve(undefined);
		});
	}
}

export class CategoryBundleService {}
