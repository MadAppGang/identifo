![components structure](https://raw.githubusercontent.com/MadAppGang/identifo/master/docs/identifo.jpg)


# identifo
Universal authentication framework for web, created with go

# components structure

![components structure](https://raw.githubusercontent.com/MadAppGang/identifo/master/docs/structure.png)


# token

As the result of auth process the user will get the [JWT token](https://tools.ietf.org/html/rfc7519).
To get more information about the JWT and get some amazing tools around that, please follow [official JWT token portal](https://jwt.io).

## JWT signature

Token signatire is used to ansure the data in the token is not changedby 3-rd party and could be trusted. 
There is [JWA specification](https://tools.ietf.org/html/rfc7518), that declare the JWT signature standarts. 

The identifo follows this specification rules and supports mandatory SH256 algoritm (HMAC using SHA-256).The Identifo implements two recommended algorithms as well:
- RS256 (RSASSA PKCS1 V1.5 using SHA256)
- ES256 (ECDSA using P-256 and SHA-256)

Future support of optional algorithms is open for discussion:
- HS384, HS512
- RS384, RS512
- ES384, ES512
- PS256, PS384, PS512

Some summary on algorithms:
- HMAC allows the messages to be verified throght the shared secret. Anyone with the key could also create the signature.
- SHA-XXX is the family of hashing algorythms. They gey the data of arbitary length and produce the output of fixed length. Trying to reduce collision and not reversable. So even one bit of change of original message should produce the completely different output with good hashing function. And it should be fast as possible.
- RSASSA is RSA algorithm, adapted fro signatures. As RSA it's assymetric keys algorithm, using keypair of public and secret keys. The main difference of RSASSA is that private key could be used create and verify signature, and public only to verify. Public key could not be used to create signatures. It's kind of one-way one-to-many pattern.
- ECDSA - the replacement for RSA. It's ulilise the different math. And because of thay could prvide the same level of security with less power usage. It's signiffically faster. 

The key is provided by identifo regarding [JWK (JSON web key) cpecification](https://tools.ietf.org/html/rfc7517).

## JWT encryption

We are not supporting JWT encyption in this verison, buy you could embed some encrypted data in the user payload.
