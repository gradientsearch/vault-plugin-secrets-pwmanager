import { expect, test } from "vitest";
import { buildUUK, twoSkd } from "./uuk";

test('buildUUK', async () => {
    const textEncoder = new TextEncoder();
    let password =textEncoder.encode("typingcats")
    let mount =textEncoder.encode("pwmanager")
    let password2 =textEncoder.encode("typingcats")

    let secretKey = crypto.getRandomValues(new Uint8Array(16))
    let entityID = textEncoder.encode(crypto.randomUUID())

   let uuk =  await buildUUK(password, mount, secretKey, entityID)
   let bits = await twoSkd(uuk, password, mount, secretKey, entityID)
   let bits2 = await twoSkd(uuk, password2, mount, secretKey, entityID)

    expect(bits).toEqual(bits2);
});

test('crypto', async () => {

    let dec = new TextDecoder()
    let enc = new TextEncoder()

    let symmetricKey = crypto.getRandomValues(new Uint8Array(32)).toString()
    const iv = crypto.getRandomValues(new Uint8Array(12));

});