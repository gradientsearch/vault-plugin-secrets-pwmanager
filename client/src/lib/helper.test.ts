import { expect, test } from "vitest";
import { createKey, exportKey } from "./helper";

test('test create/export key', async () => {
    let key = await createKey()
    let jwk = await exportKey(key)
    let json = JSON.stringify(jwk)
    expect(json, 'should export symmetric key as jwk').toContain('"alg":"A256GCM"')
})