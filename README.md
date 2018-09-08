# Identifo

![components structure](https://raw.githubusercontent.com/MadAppGang/identifo/master/docs/identifo.jpg)

## General information

Universal authentication framework for web, created with go.
It follows the [OpenID connect 1.0](https://openid.net/specs/openid-connect-core-1_0.html)  and [OAuth2](https://tools.ietf.org/html/rfc6749) specifications.

OpenID connect 1.0 is new protocol.

From [openid.net](https://openid.net/connect/), “OpenID Connect 1.0 is a simple identity layer on top of the OAuth 2.0 protocol. It allows Clients to verify the identity of the End-User based on the authentication performed by an Authorization Server, as well as to obtain basic profile information about the End-User in an interoperable and REST-like manner.” This “REST-like manner” makes OIDC more like an API (in line with OAuth2) than the previous generations of OpenID. OIDC extends the OAuth2 Authorization Code Grant (three-legged OAuth).

We are highly recommend to dive in and understand the purpose and mechanisms of this specification. OIDC is the most secure, flexible and widely adopted way for authorization.

[Modern authentication standards explained series](https://medium.com/@robert.broeckelmann/saml-v2-0-vs-jwt-series-550551f4eb0d)
[Understanding OpenID Connect Part 1](https://medium.com/@robert.broeckelmann/saml2-vs-jwt-understanding-openid-connect-part-1-fffe0d50f953)
[Understanding OpenID Connect Part 2](https://medium.com/@robert.broeckelmann/saml2-vs-jwt-understanding-openid-connect-part-2-f361ca867baa)
[Understanding OpenID Connect Part 3](https://medium.com/@robert.broeckelmann/saml2-vs-jwt-understanding-openid-connect-part-3-b81c5aa9bc20)

## Components structure

![components structure](https://raw.githubusercontent.com/MadAppGang/identifo/master/docs/structure.png)

## Token

As the result of auth process the user will get the [JWT token](https://tools.ietf.org/html/rfc7519).
To get more information about the JWT and get some amazing tools around that, please follow [official JWT token portal](https://jwt.io).

### JWT signature

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
- ECDSA - the replacement for RSA. It's ulilise the different math. And because of thay could provide the same level of security with less power usage. It's signiffically faster. 

The key is provided by identifo regarding [JWK (JSON web key) specification](https://tools.ietf.org/html/rfc7517).

### JWT encryption

We are not supporting JWT encryption in this version, buy you could embed some encrypted data in the user payload.

## Components

Everything build on negroni and gorilla mux

Negroni uses this middlewares:

[Authorization](https://github.com/casbin/negroni-authz)
[Data Binding](https://github.com/mholt/binding)
[CORS Headers](https://github.com/rs/cors)
[GZIP response](https://github.com/phyber/negroni-gzip)
[JWT middleware](https://github.com/auth0/go-jwt-middleware)
[Logrus support](https://github.com/meatballhat/negroni-logrus)
[OAuth2 support](https://github.com/goincremental/negroni-oauth2)

## Useful information

[Understanding sessions](https://blog.questionable.services/article/map-string-interface/)