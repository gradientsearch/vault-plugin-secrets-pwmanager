interface HvBundles {
	request_id: string;
	lease_id: string;
	lease_duration: number;
	renewable: boolean;
	data: Data;
	warnings: null;
	mount_type: string;
}

interface Data {
	bundles: HvBundle[];
	shared_bundles: any;
	// Only present on create bundle payload
	path: string;
}

interface HvBundle {
	created: number;
	path: string;
	id: string;
	owner_entity_id: string;
	users: HvUser[];
}

interface HvUser {
	capabilities: string;
	entity_id: string;
	entity_name: string;
	is_admin: boolean;
	shared_timestamp: number;
}
