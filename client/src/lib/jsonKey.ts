type JSONValue = string | number | boolean | null | JSONObject | JSONArray
type JSONObject = { [key: string]: JSONValue }
type JSONArray = JSONValue[]

function convertKey(key: string): string {
    let newKey = ""
    for (let i = 0; i < key.length; i++) {
        if (key[i] ===
            key[i].toUpperCase()) {
            if (i == 0) {
                newKey += key[i].toLowerCase();
            } else {
                newKey += "_" + key[i].toLowerCase();
            }
        } else {
            newKey += key[i]
        }
    }
    return newKey
}

function convertCase(obj: JSONObject): JSONObject {
    return Object.fromEntries(
        Object.entries(obj).map(([key, value]) => [
            typeof value === 'object' && value !== null
                ? Array.isArray(value)
                    ? value.map((item) =>
                        typeof item === 'object' && item !== null
                            ? convertCase(item as JSONObject)
                            : item
                    )
                    : convertCase(value as JSONObject)
                : convertKey(value as string),
        ])
    )
}