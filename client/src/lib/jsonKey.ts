export type JSONValue = string | number | boolean | null | JSONObject | JSONArray;
export type JSONObject = { [key: string]: JSONValue };
export type JSONArray = JSONValue[];

/**
 *
 * @param key json key
 * @returns key converted to snake case
 */
function convertKeyToJson(key: string): string {
	let newKey = '';
	for (let i = 0; i < key.length; i++) {
		if (key[i] === key[i].toUpperCase() && !(key[i] >= '0' && key[i] <= '9')) {
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
 *
 * @param key json key
 * @returns key converted to snake case
 */
function convertKeyFromJson(key: string): string {
	let newKey = '';
	for (let i = 0; i < key.length; i++) {
		if (i == 0) {
			newKey += key[i].toUpperCase();
			continue;
		}

		if (key[i - 1] === '_') {
			newKey += key[i].toUpperCase();
			continue;
		}

		if (key[i] === '_') {
			continue;
		}

		newKey += key[i];
	}
	return newKey;
}

/**
 * Copied and updated from: @utkarshbhatt12
 * https://bigcodenerd.org/blog/converting-json-keys-upper-lower-case-typescript/
 *
 */
function _convertCase(obj: JSONObject, toJson: boolean): JSONObject {
	return Object.fromEntries(
		Object.entries(obj).map(([key, value]) => [
			toJson ? convertKeyToJson(key) : convertKeyFromJson(key),
			typeof value === 'object' && value !== null
				? Array.isArray(value)
					? value.map((item) =>
							typeof item === 'object' && item !== null
								? _convertCase(item as JSONObject, toJson)
								: item
						)
					: _convertCase(value as JSONObject, toJson)
				: value
		])
	);
}

export function convertCase(obj: any, toJson: boolean): string {
	let converted = _convertCase(JSON.parse(JSON.stringify(obj)), toJson);
	return JSON.stringify(converted);
}
export function revertCase<T>(obj: JSONObject, toJson: boolean): T {
	let reverted = _convertCase(obj, toJson);
	return JSON.parse(JSON.stringify(reverted)) as T;
}
