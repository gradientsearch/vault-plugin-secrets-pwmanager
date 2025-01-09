import type { Api } from "$lib/api";
import type { KeyPair } from "$lib/asym_key_store";

// collection of objects/services children components
// may need
export interface Zarf {

    Api?: Api
    Keypair?: KeyPair
}

