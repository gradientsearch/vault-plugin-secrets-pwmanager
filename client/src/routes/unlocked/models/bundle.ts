/**
 * A Bundle is a description of a group of entries. Could be a group of KV Entries or a group of
 * Entry types such as password.
 */
interface Bundle {
	Type: string;
	Path: string;
	Name: string;
	Owner: string;
	IsAdmin: boolean;
	ID: string;
	Users: BundleUser[]
}

interface BundleUser {
	EntityName: string
	EntityID: string
	Capabilities: any
	IsAdmin: boolean
}