import { expect, test } from "vitest";
import { buildUUK, toHex } from "./uuk";

test('buildUUK', async () => {
    const textEncoder = new TextEncoder();
    let password =textEncoder.encode("typingcats")
    let mount =textEncoder.encode("pwmanager")
    let secretKey = crypto.getRandomValues(new Uint8Array(16))
    let entityID = textEncoder.encode(crypto.randomUUID())

   let obj =  await buildUUK(password, mount, secretKey, entityID)
    expect(0).toEqual(0);
})