# go_signer
Library for support multiples method to sign

## Configuration
This library supports 3 types of signing methods: 
- **local**: it's a private key file
- **GCP**: google cloud KMS
- **AWS**: AWS KMS
- **remote**: it's a call to a remote signer service that implements [remote signing APIs](https://github.com/ethereum/remote-signing-api?tab=readme-ov-file) as [web_3signer](https://docs.web3signer.consensys.io/) **only support sign transactions**

There are a `None` method just for develop propouses

### Configuration local method
The object `SignerConfig` needs next fields:
- `SignerConfig.Method` : `local`  (you can use const `MethodLocal`)
- `SignerConfig.Config["Path"]`: full path to keystore
- `SignerConfig.Config["Password"]`: password for the keystore
```
Method = "local"
Path = "/path/to/file"
Password = "password"
```

### Configuration GCP method
The object `SignerConfig` needs next fields:
- `SignerConfig.Method` : `GCP`  (you can use const `MethodGCPKMS`)
- `SignerConfig.Config["KeyName"]`: Full path to key resource with version
Example of config file:
```
Method = "GCP" 
KeyName ="projects/your-prj-name/locations/your_location/keyRings/name_of_your_keyring/cryptoKeys/key-name/cryptoKeyVersions/version"
```
You can copy `KeyName` from console > Security > Key Managent 
Or executing `gcloud kms inventory list-keys`

### Configuration AWS method
The object `SignerConfig` needs next fields:
- `SignerConfig.Method` : `AWS`  (you can use const `MethodAWSKMS`)
- `SignerConfig.Config["KeyName"]`: Full path to key resource with version

The field `KeyName` is `KeyId` from AWS CLI: 
```
aws kms list-keys
{
    "Keys": [
        {
            "KeyId": "a47c263b-6575-4835-8721-af0bbb97XXXX",
            "KeyArn": "arn:aws:kms:us-east-1:467680962559:key/a47c263b-6575-4835-8721-af0bbb97XXX"
        }
    ]
}
```

### Configuration remote method
#### Generic configuration
The object `SignerConfig` needs next params:
- `SignerConfig.Method` : `remote` (you can use const `MethodRemoteSigner`)
- `SignerConfig.Config["URL"]`: URL to web3_signer service
- `SignerConfig.Config["Address"]`: Public address to use if there are more than 1 in web3_signer service. If there are only 1 it can be empty and the first one will be used.

- Example of config file:
```
Method = "remote"
URL = "http://localhost:9000"
Address = "0xe34243804e1f7257acb09c97d0d6f023663200c39ee85a1e6927b0b391710bbb"
```


### Configuration mock method
This method is for unittest and debug, it's not suitable for production. 
You can use a specific private key (without encryption) or generate it
#### Setting a random key
This allow to generate a random key to used: this is useful if you are setting a unittest that require a key but you don't care which one it is.
The object `MockSign` needs next params:
- `SignerConfig.Method` : `mock` (you can use const `MethodMock`)
Example: 
```
{
     Method: "mock",
}
```

#### Setting a specific private key
It is also specific for unittest, reading a key from a keystore is really fast, and maybe 
the key that are you usign is not a secret, just generated for this test. In that way you 
can setup directly in your test in no time. 

The object `MockSign` needs next params:
- `SignerConfig.Method` : `mock` (you can use const `MethodMock`)
- `SignerConfig.Config["PrivateKey"]`: private key in hex format (you can use const `FieldMockForcedPublicKey`)
Example: 
```
{
     Method: "mock",
     PrivateKey: "0xa574853f4757bfdcbb59b03635324463750b27e16df897f3d00dc6bef2997ae0",
}
```

## Support

Feel free to [open an issue](https://github.com/agglayer/go_signer/issues/new) if you have any feature request or bug report.<br />


## License

Copyright (c) 2024 PT Services DMCC

Licensed under either of

* Apache License, Version 2.0, ([LICENSE-APACHE](LICENSE-APACHE) or http://www.apache.org/licenses/LICENSE-2.0)
* MIT license ([LICENSE-MIT](LICENSE-MIT) or http://opensource.org/licenses/MIT)