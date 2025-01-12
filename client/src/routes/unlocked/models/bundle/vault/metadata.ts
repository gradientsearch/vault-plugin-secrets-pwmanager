import type { Metadata as EntryMetadata } from '../../entry';

export interface VaultMetadata {
	entries: EntryMetadata[];
}

/** TODO Prepend Hv to all Hashicorp Vault DAOs */
export interface HvMetadata {
  request_id: string;
  lease_id: string;
  lease_duration: number;
  renewable: boolean;
  data: Data2;
  warnings: null;
  mount_type: string;
}

interface Data2 {
  data: Data;
  metadata: VaultMetadata;
}

interface Metadata {
  created_time: string;
  custom_metadata: null;
  deletion_time: string;
  destroyed: boolean;
  version: number;
}

interface Data {
  entry: string;
  iv: string;
}