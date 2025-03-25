# go_signer
Library for support multiples method to sign

## Configuration
This library support 2 types of signing methods: 
- **local**: it's a private key file
- **remote**: it's a call to a remote signer service that implements [remote signing APIs](https://github.com/ethereum/remote-signing-api?tab=readme-ov-file) as [web_3signer](https://docs.web3signer.consensys.io/) 

### Configuration local method
#### Generic configuration
The object `SignerConfig` can be read from this:
```
Method = "local"
Path = "/path/to/file"
Password = "password"
```
#### Specific configuration struct
`signercommon.KeystoreFileConfig`

### Configuration remote method
#### Generic configuration
The object `SignerConfig` need next params:
- `SignerConfig.Method` : `remote`
- `SignerConfig.Config["URL"]`: URL to web3_signer service
  `SignerConfig.Config["Address"]`: Public address to use if there are more than 1 in web3_signer service. If there are only 1 it can be empty and the first one will be used.

- Example of config file:
```
Method = "remote"
URL = "http://localhost:9000"
Address = "0xe34243804e1f7257acb09c97d0d6f023663200c39ee85a1e6927b0b391710bbb"
```

## Support

Feel free to [open an issue](https://github.com/agglayer/go_signer/issues/new) if you have any feature request or bug report.<br />


## License

Copyright (c) 2024 PT Services DMCC

Licensed under either of

* Apache License, Version 2.0, ([LICENSE-APACHE](LICENSE-APACHE) or http://www.apache.org/licenses/LICENSE-2.0)
* MIT license ([LICENSE-MIT](LICENSE-MIT) or http://opensource.org/licenses/MIT)