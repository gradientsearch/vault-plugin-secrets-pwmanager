# 🔐 vault-plugin-secrets-pwmanager

This project creates a plugin that can be registered with `HashiCorp Vault`. The plugin is responsible for managing `pwmanager` user policies and user KV-V2 secret mounts, which store user passwords.

In addition to the plugin, a web client is provided to offer a user-friendly interface for interacting with the plugin’s password management capabilities. User passwords are encrypted client-side with a symmetric key before being sent to `HashiCorp Vault`. This ensures that if a `HashiCorp Vault` administrator views a user’s password, it will be encrypted. The user’s private key is stored in an encrypted bundle known as the User Unique Key (UUK). The UUK contains the necessary information to decrypt the user’s private key using a two-secret key derivation function. This private key can then decrypt the required keys for encrypting and decrypting data.


## ⚠️ Early Development - Not for Production Use
This repository is currently in early development. The code may be unstable, incomplete, or subject to significant changes. Do not use it in production environments at this time.

We recommend that you use this repository for testing and experimentation purposes only, and proceed with caution. Contributions and feedback are welcome as we continue to improve the project!


### Register User Flow

![register-flow drawio](https://github.com/user-attachments/assets/f590c38c-683e-483c-a813-fca52cce3b37)
