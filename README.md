# 🔐 vault-plugin-secrets-pwmanager

This project creates a plugin that can be registered with `HashiCorp Vault`. The plugin is responsible for managing `pwmanager` user policies and user KV-V2 secret mounts, which store user passwords.

In addition to the plugin, a web client is provided to offer a user-friendly interface for interacting with the plugin’s password management capabilities. User passwords are encrypted client-side with a symmetric key before being sent to `HashiCorp Vault`. This ensures that if a `HashiCorp Vault` administrator views a user’s password, it will be encrypted. The user’s private key is stored in an encrypted bundle known as the User Unique Key (UUK). The UUK contains the necessary information to decrypt the user’s private key using a two-secret key derivation function. This private key can then decrypt the required keys for encrypting and decrypting data.


## ⚠️ Early Development - Not for Production Use
This repository is currently in early development. The code may be unstable, incomplete, or subject to significant changes. Do not use it in production environments at this time.

We recommend that you use this repository for testing and experimentation purposes only, and proceed with caution. Contributions and feedback are welcome as we continue to improve the project!

## Milestone Demo

> [!NOTE] 
> I will continue to make short videos like this to show progress and review code, so future contributors can ramp up quickly. You can find these updates on my YouTube channel: [@gradientsearch](https://www.youtube.com/@gradientsearch).


[create password demo](https://github.com/user-attachments/assets/5176f4a6-4bfd-4d4e-9862-b7f8e61d97fa)

> In this demo, I showcase an open-source password manager I’m currently developing, which integrates with HashiCorp Vault via a Vault plugin. I walk through the process of registering a user, unlocking the password manager, and creating and viewing passwords. Additionally, I explain how password entries are stored in Vault using the KV-v2 secret engine, including the encryption and decryption processes, and how Vault policies control access to these entries.


### Register User Flow

![register-flow drawio](https://github.com/user-attachments/assets/f590c38c-683e-483c-a813-fca52cce3b37)





## License

The client code for this project is licensed under the MIT License. You are free to use, modify, and distribute this code, including for commercial purposes, provided that you include the original copyright notice and disclaimers.

The plugin code is licensed under the Mozilla Public License 2.0 (MPL-2.0). This is because I used the HashiCups demo as a starting point for developing the plugin. The MPL allows for the use, modification, and distribution of the code, but it requires that any modifications to the plugin code be released under the same MPL license.
