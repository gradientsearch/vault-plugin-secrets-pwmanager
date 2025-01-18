import {
	pubkeyEncrypt,
	exportJwkKey,
	generateSymmetricKey,
	prikeyDecrypt,
	symmetricEncrypt,
	importJWKkey,
	symmetricDecrypt
} from '$lib/helper';
import type { EncryptedEntry } from '../models/bundle/vault/entry';
import type { BundleMetadata as BundleMetadata } from '../models/bundle/vault/metadata';
import type { Entry, Metadata } from '../models/entry';
import type { Zarf } from '../models/zarf';
import { userService } from './user.service';

/**
 * Interface the BundleView uses to interface with different types of bundles e.g a Vault bundle or
 * a Category bundle.
 */
export interface BundleService {
	putEntry(pi: Entry): Promise<Error | undefined>;
	getMetadata(): Promise<[BundleMetadata | undefined, Error | undefined]>;
	init(): Promise<Error | undefined>;
}

/**
 * A KVBundleService is responsible for interfacing with a HashiCorp Vault KV2 secret mount. w.r.t.
 * pwmanager. A KV2 secret mount is a bundle of password entries. A bundle has the following path convention
 * `bundles/{{ identity.entity.id }}/<bundle guid>
 *
 * The KV2 secret mount contains the following paths:
 * - keys/{{ identity.entity.id }}: each bundle has a symmetric key used to encrypt all secrets. That
 * symmetric keys is encrypted with users public key
 * - metadata: bundle metadata for each entry
 * - entries/<entry name>: user entires
 */
export class KVBundleService implements BundleService {
	onEntriesChanged: Function;
	zarf: Zarf;
	bundle: Bundle;
	symmetricKey: CryptoKey | undefined;

	constructor(zarf: Zarf, bundle: Bundle, onEntriesChanged: Function) {
		this.zarf = zarf;
		this.bundle = bundle;
		this.onEntriesChanged = onEntriesChanged;
	}

	async init(): Promise<Error | undefined> {
		// get decryption key for this vault
		//whats the entity id

		let entityID = userService.getEntityID();
		let [key, err] = await this.zarf.Api.getBundleSymmetricKey(this.bundle, entityID);
		if (err?.toString().includes('404 not found')) {
			let [key, err] = await this.createBundleEncryptionKey(entityID);

			if (err !== undefined) {
				return err;
			}

			err = await this.createBundleMetadata();
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
			return Error('error retrieving bundle symmetric key');
		}
	}

	async createBundleEncryptionKey(
		entityID: string
	): Promise<[CryptoKey | undefined, Error | undefined]> {
		// TODO make this a version CAS version 0 only operation. Never want to overwrite a bundle symmetric key
		let key = await generateSymmetricKey();
		let jwk = await exportJwkKey(key);
		let encrypted = await pubkeyEncrypt(
			new TextEncoder().encode(JSON.stringify(jwk)),
			this.zarf.Keypair.PubKey
		);
		let err = await this.zarf.Api.PutUserKey(this.bundle, entityID, encrypted);
		if (err !== undefined) {
			console.log('err: ', err);
			return [undefined, Error('error retrieving bundle symmetric key : ', err)];
		}
		this.symmetricKey = key;
		return [key, undefined];
	}

	async createBundleMetadata(): Promise<Error | undefined> {
		let metadata: BundleMetadata = {
			entries: []
		};
		let [data, err] = await this.encryptPayload(metadata);
		if (err !== undefined) {
			return Error('error encrypted bundle metadata');
		}

		err = await this.zarf.Api.PutMetadata(this.bundle, data);
		if (err !== undefined) {
			return Error('error creating encrypted bundle metadata');
		}
	}

	async encryptPayload(payload: any): Promise<[string | undefined, Error | undefined]> {
		// TODO refactor this class so this isn't necessary
		if (this.symmetricKey === undefined) {
			return [undefined, Error('bundle encryption key is undefined')];
		}

		let [iv, encrypted] = await symmetricEncrypt(
			new TextEncoder().encode(JSON.stringify(payload)),
			this.symmetricKey
		);

		let ed: EncryptedEntry = {
			entry: encrypted,
			iv: iv
		};

		let data = {
			data: ed
		};

		return [JSON.stringify(data), undefined];
	}

	async decryptPayload(ee: EncryptedEntry): Promise<[string | undefined, Error | undefined]> {
		// TODO refactor this class so this isn't necessary
		if (this.symmetricKey === undefined) {
			return [undefined, Error('bundle encryption key is undefined')];
		}

		let decrypted = await symmetricDecrypt(ee.entry, ee.iv, this.symmetricKey);

		return [decrypted, undefined];
	}

	async getMetadata(): Promise<[BundleMetadata | undefined, Error | undefined]> {
		let [md, err] = await this.zarf.Api.getMetadata(this.bundle);
		if (err !== undefined) {
			return [undefined, err];
		}

		if (md === undefined) {
			return [undefined, Error('no data returned from server')];
		}

		let [plaintext, err2] = await this.decryptPayload(md.data.data);
		if (err2 !== undefined) {
			return [undefined, err2];
		}

		if (plaintext === undefined) {
			return [undefined, Error('decrypted metadata entry undefined')];
		}
		let vm = JSON.parse(plaintext) as BundleMetadata;

		return [vm, undefined];
	}

	// TODO: rename this to put entry
	async putEntry(e: Entry): Promise<Error | undefined> {
		//store data in vault
		// encrypt
		let newEntry = e.Metadata.ID.length === 0;
		if (newEntry) {
			let entryName = crypto.randomUUID();
			e.Metadata.ID = entryName;
		}

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
			return Error('error retrieving latest bundle metadata');
		}

		if (metadata === undefined) {
			return Error('error metadata should be defined but was undefined');
		}

		if (newEntry) {
			metadata?.entries.push(e.Metadata);
		} else {
			// loop through and update metadata if it exists
			let updatedMetadata: Metadata[] = [];
			metadata?.entries.forEach((me) => {
				if (me.ID === e.Metadata.ID) {
					updatedMetadata.push(e.Metadata);
				} else {
					updatedMetadata.push(me);
				}
			});

			metadata.entries = updatedMetadata;
		}

		let ep = await this.encryptPayload(metadata);

		// TODO add CAS version from
		let err4 = await this.zarf.Api.PutMetadata(this.bundle, ep);
		if (err4 !== undefined) {
			return Error('error putting metadata: ', err4);
		}

		// TODO if error delete entry
		this.onEntriesChanged(metadata);
		return undefined;
	}

	async getEntry(m: Metadata): Promise<[Entry | undefined, Error | undefined]> {
		let [hee, err] = await this.zarf.Api.GetEntry(this.bundle, m.ID);

		if (err !== undefined) {
			return [undefined, Error(`error getting entry: ${err}`)];
		}

		if (hee === undefined) {
			return [undefined, Error(`error - encrypted entry is undefined: ${err}`)];
		}

		let [payload, err2] = await this.decryptPayload(hee.data.data);

		if (err2 !== undefined || payload === undefined) {
			return [undefined, Error(`error decrypting EncryptedEntry: ${err}`)];
		}

		let e = JSON.parse(payload);
		return [e, undefined];
	}

	async deleteEntry(id: string): Promise<Error | undefined> {
		let [metadata, err] = await this.getMetadata();
		if (err !== undefined) {
			return Error('error retrieving latest bundle metadata');
		}

		if (metadata === undefined) {
			return Error('error metadata should be defined but was undefined');
		}

		// loop through and and remove metadata matching
		let updatedMetadata: Metadata[] = [];
		metadata?.entries.forEach((me) => {
			if (me.ID !== id) {
				updatedMetadata.push(me);
			}
		});

		metadata.entries = updatedMetadata;

		let ep = await this.encryptPayload(metadata);

		// TODO add CAS version from
		let err4 = await this.zarf.Api.PutMetadata(this.bundle, ep);
		if (err4 !== undefined) {
			return Error('error putting metadata: ', err4);
		}

		return undefined;
	}
}

export class CategoryBundleService {}
