import { expect, test } from "vitest";
import { generateSymmetricKey, exportJwkKey } from "./helper";

test('test create/export key', async () => {
    let key = await generateSymmetricKey()
    let jwk = await exportJwkKey(key)
    let json = JSON.stringify(jwk)
    expect(json, 'should export symmetric key as jwk').toContain('"alg":"A256GCM"')
})