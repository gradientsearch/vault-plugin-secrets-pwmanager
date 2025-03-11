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
	putEntry(pi: Entry, bm: BundleMetadata): Promise<[string | undefined, Error | undefined]>;
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
			if (this.bundle.Owner === entityID) {
				let [key, err] = await this.createBundleEncryptionKey(entityID);

				if (err !== undefined) {
					return err;
				}

				err = await this.createBundleMetadata();
				if (err != undefined) {
					return err;
				}
			} else {
				alert('no bundle key available for user');
				return;
			}
		}
		// TODO this is a little messy. Let's try and clean this up.
		if (key !== undefined && this.zarf.Keypair.PriKey) {
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
		if (!this.zarf.Keypair.PubKey) {
			//TODO handle pubkey undefined
			return [undefined, Error('public key is undefined')];
		}

		// TODO make this a version CAS version 0 only operation. Never want to overwrite a bundle symmetric key
		let key = await generateSymmetricKey();
		let jwk = await exportJwkKey(key);
		let encrypted = await pubkeyEncrypt(
			new TextEncoder().encode(JSON.stringify(jwk)),
			this.zarf.Keypair.PubKey
		);

		let data = {
			data: {
				key: encrypted
			},
			options: {
				cas: 0
			}
		};
		let err = await this.zarf.Api.PutUserKey(this.bundle, entityID, data);
		if (err !== undefined) {
			return [undefined, Error('error creating bundle encryption key: ', err)];
		}
		this.symmetricKey = key;
		return [key, undefined];
	}

	async encryptBundleKey(pubkey: CryptoKey, entityID: string): Promise<Error | undefined> {
		if (!pubkey) {
			return Error('public key is undefined');
		}

		if (this.symmetricKey === undefined) {
			return Error('bundle key is undefined');
		}
		// TODO make this a version CAS version 0 only operation. Never want to overwrite a bundle symmetric key
		let jwk = await exportJwkKey(this.symmetricKey);
		let encrypted = await pubkeyEncrypt(new TextEncoder().encode(JSON.stringify(jwk)), pubkey);
		let data = {
			data: {
				key: encrypted
			},
			options: {
				cas: 0
			}
		};
		let err = await this.zarf.Api.PutUserKey(this.bundle, entityID, data);
		if (err !== undefined) {
			return Error(`error setting bundle symmetric key for ${entityID}`);
		}
	}

	async createBundleMetadata(): Promise<Error | undefined> {
		let metadata: BundleMetadata = {
			entries: [],
			bundleName: this.bundle.Name,
			version: 0
		};
		let [data, err] = await this.encryptPayload(metadata, 0);
		if (err !== undefined || data === undefined) {
			return Error('error encrypted bundle metadata');
		}
		err = await this.zarf.Api.PutMetadata(this.bundle, data);
		if (err !== undefined) {
			return Error('error creating encrypted bundle metadata');
		}
	}

	async encryptPayload(
		payload: any,
		version: number
	): Promise<[string | undefined, Error | undefined]> {
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
			data: ed,
			options: {
				cas: version
			}
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

		if (err !== undefined) {
			console.log(`err getting bundle metadata: ${err}`);
		}

		// Write latest encrypted bundle metadata to localstorage
		if (md !== undefined) {
			localStorage.setItem(`${this.bundle.Path}/metadata`, JSON.stringify(md.data));
		}

		let [plaintext, err2] = await this.decryptPayload(md.data.data);
		if (err2 !== undefined) {
			return [undefined, err2];
		}

		if (plaintext === undefined) {
			return [undefined, Error('decrypted metadata entry undefined')];
		}
		let bm = JSON.parse(plaintext) as BundleMetadata;
		bm.version = md.data.metadata.version;

		return [bm, undefined];
	}

	async putEntry(
		e: Entry,
		metadata: BundleMetadata
	): Promise<[string | undefined, Error | undefined]> {
		//store data in vault
		// encrypt
		let pathToReplace;

		let newEntry = e.Metadata.ID.length === 0;
		if (!newEntry) {
			pathToReplace = e.Metadata.Path;
			e.Metadata.Path = crypto.randomUUID();
		} else {
			e.Metadata.ID = crypto.randomUUID();
			e.Metadata.Version = 0;
			e.Metadata.Path = crypto.randomUUID();
		}

		// CAS version is always zero since password entries are immutable.
		let [data, err2] = await this.encryptPayload(e, 0);
		if (err2 != undefined) {
			return [undefined, err2];
		}

		let err3 = await this.zarf.Api.PutEntry(this.bundle, data, e.Metadata.Path);
		if (err3 !== undefined) {
			return [undefined, Error(`error putting entry:  ${err3.message}`)];
		}

		// have to increment version
		e.Metadata.Version = e.Metadata.Version + 1;

		if (metadata === undefined) {
			return [undefined, Error('error metadata should be defined but was undefined')];
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

		let ep = await this.encryptPayload(metadata, metadata.version);

		let err4 = await this.zarf.Api.PutMetadata(this.bundle, ep);
		if (err4 !== undefined) {
			await this.deleteEntry(e.Metadata.Path);
			return [undefined, Error('error putting metadata: ', err4)];
		}

		if (pathToReplace !== undefined) {
			await this.deleteEntry(pathToReplace);
		}
		this.onEntriesChanged();
		return [e.Metadata.Path, undefined];
	}

	async getEntry(m: Metadata): Promise<[Entry | undefined, Error | undefined]> {
		let [hee, err] = await this.zarf.Api.GetEntry(this.bundle, m.Path);

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
		e.Version = hee.data.metadata.version;

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
		// TODO: do a filter function instead
		let updatedMetadata: Metadata[] = [];
		metadata?.entries.forEach((me) => {
			if (me.ID !== id) {
				updatedMetadata.push(me);
			}
		});

		metadata.entries = updatedMetadata;

		let ep = await this.encryptPayload(metadata, metadata.version);

		// TODO add CAS version from
		let err2 = await this.zarf.Api.PutMetadata(this.bundle, ep);
		if (err2 !== undefined) {
			return Error('error putting metadata: ', err2);
		}

		let err3 = await this.zarf.Api.DestroyEntry(this.bundle, id);
		if (err3 !== undefined) {
			// TODO: this will be an orphaned entry if we don't retry b/c
			// we removed the entry from the metadata and thats how we keep track
			// of password entries.
			return Error('error deleting password entry: ', err3);
		}

		this.onEntriesChanged(metadata);
		return undefined;
	}

	static async initBundles(
		zarf: Zarf,
		hvBundles: HvBundle[]
	): Promise<[Bundle[] | undefined, Error | undefined]> {
		let bs: Bundle[] = [];

		for (let i = 0; i < hvBundles.length; i++) {
			let hvb = hvBundles[i];

			let users: BundleUser[] = [];
			hvb.users?.forEach((u) => {
				let caps: any = {};
				u.capabilities.split(',').forEach((c) => {
					caps[c] = true;
				});
				let bu: BundleUser = {
					EntityName: u.entity_name,
					EntityId: u.entity_id,
					Capabilities: caps,
					IsAdmin: u.is_admin
				};
				users.push(bu);
			});

			let b: Bundle = {
				Type: 'bundle',
				Path: hvb.path,
				Name: '',
				Owner: hvb.owner_entity_id,
				IsAdmin: false,
				ID: hvb.id,
				Users: users
			};

			let bundleService = new KVBundleService(zarf, b, () => {});
			let initErr = await bundleService.init();
			if (initErr !== undefined) {
				console.log(`err initing bundle bundleService: ${initErr}`);
			}
			// TODO read from localstorage/indexedDB
			let encryptedCachedMetadata = localStorage.getItem(`${hvb.path}/metadata`);

			let m: BundleMetadata = {
				entries: [],
				bundleName: '',
				version: 0
			};
			if (encryptedCachedMetadata) {
				let json = JSON.parse(encryptedCachedMetadata);
				let [cahcedMetadata, err] = await bundleService.decryptPayload(json.data);

				if (err !== undefined || cahcedMetadata === undefined) {
					//TODO handle error
				} else {
					m = JSON.parse(cahcedMetadata);
				}
			} else {
				let [bm, err] = await bundleService.getMetadata();

				if (err !== undefined) {
					console.log(`err getting bundle metadata: ${err}`);
				}

				if (bm !== undefined) {
					m = bm;
					let [em, err] = await bundleService.encryptPayload(m, m.version);
					if (err !== undefined || em === undefined) {
						//TODO Handle error
					} else {
						localStorage.setItem(`${hvb.path}/metadata`, em);
					}
				}
			}

			b.Name = m.bundleName;
			bs.push(b);
		}

		return [bs, undefined];
	}

	// getBundles retrieves all bundles and bundle metadata.
	static async getBundles(zarf: Zarf): Promise<[any | undefined, Error | undefined]> {
		let [hvBundles, err] = await zarf.Api.GetBundles();

		if (err !== undefined) {
			return [undefined, Error(`error getting bundles: ${err}`)];
		}

		if (hvBundles === undefined) {
			return [undefined, Error(`error getting bundles: bundles should not be undefined`)];
		}

		let [bundles, err2] = await KVBundleService.initBundles(zarf, hvBundles.bundles);

		if (err2 !== undefined) {
			return [undefined, Error(`error getting bundles: ${err2}`)];
		}

		let sharedBundles: HvBundle[] = [];
		if (hvBundles.shared_bundles !== undefined) {
			Object.keys(hvBundles.shared_bundles).forEach((k) => {
				let sb = hvBundles.shared_bundles[k];
				let b: HvBundle = {
					created: sb.created,
					path: sb.path,
					id: sb.id,
					owner_entity_id: sb.owner_entity_id,
					users: []
				};
				sharedBundles.push(b);
			});
		}

		let [sbs, err3] = await KVBundleService.initBundles(zarf, sharedBundles);

		if (err3 !== undefined) {
			return [undefined, Error(`error getting bundles: ${err3}`)];
		}

		return [{ bundles: bundles, sharedBundles: sbs }, undefined];
	}

	static async createBundle(
		zarf: Zarf,
		name: string
	): Promise<[Bundle | undefined, Error | undefined]> {
		let [path, err] = await zarf.Api.CreateBundle();
		if (err !== undefined) {
			return [undefined, Error(`error creating bundle: ${err}`)];
		}

		if (path === undefined) {
			return [undefined, Error(`error creating bundle: bundles should not be undefined`)];
		}

		let b: Bundle = {
			Type: 'bundle',
			Path: path,
			Name: name,
			Owner: userService.getEntityID(),
			IsAdmin: false,
			ID: '',
			Users: []
		};

		let bundleService = new KVBundleService(zarf, b, () => {});
		await bundleService.init();

		return [b, undefined];
	}

	async updateSharedBundleUsers(
		zarf: Zarf,
		ownerEntityID: string,
		bundleID: string,
		users: BundleUser[]
	): Promise<[any[] | undefined, Error | undefined]> {
		let [pubkey, err] = await zarf.Api.updateSharedBundleUsers(ownerEntityID, bundleID, users);
		if (err !== undefined) {
			return [undefined, Error(`error updating bundle users: ${err}`)];
		}

		if (pubkey === undefined) {
			return [undefined, Error(`error updating shared bundle: bundles should return pubkeys`)];
		}

		return [pubkey, undefined];
	}
}

export class CategoryBundleService {}
