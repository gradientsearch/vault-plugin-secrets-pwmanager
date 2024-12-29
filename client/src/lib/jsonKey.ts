export type JSONValue = string | number | boolean | null | JSONObject | JSONArray;
export type JSONObject = { [key: string]: JSONValue };
export type JSONArray = JSONValue[];

/**
 *
 * @param key json key
 * @returns key converted to snake case
 */
function convertKey(key: string): string {
	let newKey = '';
	for (let i = 0; i < key.length; i++) {
		if (key[i] === key[i].toUpperCase()) {
			if (i == 0) {
				newKey += key[i].toLowerCase();
			} else {
				newKey += '_' + key[i].toLowerCase();
			}
		} else {
			newKey += key[i];
		}
	}
	return newKey;
}

/**
 * Copied and updated from: @utkarshbhatt12
 * https://bigcodenerd.org/blog/converting-json-keys-upper-lower-case-typescript/
 *
 */
export function convertCase(obj: JSONObject): JSONObject {
	return Object.fromEntries(
		Object.entries(obj).map(([key, value]) => [
			convertKey(key),
			typeof value === 'object' && value !== null
				? Array.isArray(value)
					? value.map((item) =>
							typeof item === 'object' && item !== null ? convertCase(item as JSONObject) : item
						)
					: convertCase(value as JSONObject)
				: value
		])
	);
}
