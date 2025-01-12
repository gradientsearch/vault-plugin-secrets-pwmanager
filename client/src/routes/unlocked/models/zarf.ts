import type { Api } from "$lib/api";
import type { KeyPair } from "$lib/asym_key_store";


/**
 * A Zarf is a collection of objects/services children components need
 */
export interface Zarf {
    Api: Api
    Keypair: KeyPair
}

