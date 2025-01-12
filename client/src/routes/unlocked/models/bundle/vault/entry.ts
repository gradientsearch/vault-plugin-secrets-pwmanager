export interface EncryptedEntry {
	entry: string;
	iv: string;
}

export  interface HvEncryptedEntry {
  request_id: string;
  lease_id: string;
  renewable: boolean;
  lease_duration: number;
  data: Data2;
  wrap_info: null;
  warnings: null;
  auth: null;
  mount_type: string;
}

interface Data2 {
  data: EncryptedEntry;
  metadata: Metadata;
}

interface Metadata {
  created_time: string;
  custom_metadata: null;
  deletion_time: string;
  destroyed: boolean;
  version: number;
}
