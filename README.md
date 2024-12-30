# 🔐 vault-password-manager-plugin

This project creates a plugin that can be registered with HashiCorp Vault. The plugin is responsible for managing password manager user policies and user vaults (i.e., a user vault is a KV-V2 secret store). In addition to the plugin, a web client is provided, offering a user-friendly way to interact with the plugin's password management capabilities. User passwords are encrypted client-side with the user's private key before the data is sent to Vault. This means Vault admins cannot view users' secrets. The user's private key is encrypted in an encryption bundle called the User Unique Key (UUK). The UUK contains the information required to decrypt the user's private key through a two-secret key derivation function.



## ⚠️ Early Development - Not for Production Use
This repository is currently in early development. The code may be unstable, incomplete, or subject to significant changes. Do not use it in production environments at this time.

We recommend that you use this repository for testing and experimentation purposes only, and proceed with caution. Contributions and feedback are welcome as we continue to improve the project!


### Register User Flow

![register-flow drawio](https://github.com/user-attachments/assets/f590c38c-683e-483c-a813-fca52cce3b37)
