# Identifo

![components structure](https://raw.githubusercontent.com/MadAppGang/identifo/master/docs/identifo.jpg)

## General information

Universal authentication framework for web, created with go.

## Components structure

![components structure](https://raw.githubusercontent.com/MadAppGang/identifo/master/docs/structure.png)

## Motivation

There is a lot of Authorization and Authentication services on the market now.

Some of them are commercial SaaS systems with one click integration like [Auth0](https://auth0.com), [AWS Cognito](https://aws.amazon.com/cognito/), [Firebase auth](https://firebase.google.com/docs/auth/).

Some of them are self-hosted and open-sourced, like [OauthD](https://github.com/oauth-io/oauthd), [Gluu Server](https://www.gluu.org/),
[Anvil](http://anvil.io/#features).

Almost all of them provide OpenID Connect and Oauth2 login flow. It makes an excellent sense for expected use cases. The service is external to your servers, and OpenID Connect is one of the most secure ways to provide auth service in this architectural solution.

Someday we have started to implement Strongbox project. It's highly secure team messenger.  We are unable to use external services in this case, because it decreases the security level by storing user data and credentials on 3rd party services somewhere in the cloud. Also, it makes impossible to host this solution on premises, creating external dependency.

We could use a self-hosted solution. However, we had to keep all our codebase with Go. All of the current candidates provided OpenID Connect, OAuth2 login flow. It's overcomplicated for single app solution. If you have one backend, that hosts your web app and provides REST API for native clients - you own users passwords. So there is no sense in using OIDC in that case.

So we decided to implement Identity Provider with these key features:

- easy to integrate into any infrastructure
- fast, small and binary distributed
- uses JWT and other OIDC principles, avoid all levels of redirections
- provide default user persistent storage service
- provide an easy way to integrate any other user persistent layer
- provide default login/register/reset password forms 
- provide the way to restyle login/register/reset form
- implement token introspection
- implement machine-machine tokens (analog to Oauth2 Client Credentials Grant)
- optional management console to set up, monitor and configure the service

Although Identifo is not intended to implement OpenID Connect flow, you could easily integrate it with [Ory Hydra](https://www.ory.sh). Just with a couple of minutes, you are able to get the complete OIDC support. The tutorial and instructions TBA.

## Token

As the result of auth process the user will get the [JWT token](https://tools.ietf.org/html/rfc7519).
To get more information about the JWT and get some amazing tools around that, please follow [official JWT token portal](https://jwt.io).

### JWT signature

Token signature is used to ensure the data in the token is not changed by 3-rd party and could be trusted. 
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
- SHA-XXX is the family of hashing algorythms. They get the data of arbitary length and produce the output of fixed length. Trying to reduce collision and not reversable. So even one bit of change of original message should produce the completely different output with good hashing function. And it should be fast as possible.
- RSASSA is RSA algorithm, adapted from signatures. As RSA it's assymetric keys algorithm, using keypair of public and secret keys. The main difference of RSASSA is that private key could be used to create and verify signature, and public key only to verify. Public key could not be used to create signatures. It's kind of one-way one-to-many pattern.
- ECDSA - the replacement for RSA. It ulilises different math. And because of that could provide the same level of security with less power usage. It's significantly faster. 

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

## Getting started
First generate token with `./jwt/generate_token.sh `
Start Identifo in docker with `docker-compose up`
Open http://localhost:8081/adminpanel/ and use default credentials `admin@admin.com` `password` to login

## Useful information

[Understanding sessions](https://blog.questionable.services/article/map-string-interface/)

## Test

We have covered everything we count not trivial. Tests are created for documentation and code validation purposes. If you don't understand something and could not find it in the docs, please try to find in in the tests as well. If you think there is missing test, PR is always welcome or just let us know.

[Test http requests](https://blog.questionable.services/article/testing-http-handlers-go/)

## References

[HMAC signature](https://blog.andrewhoang.me/how-api-request-signing-works-and-how-to-implement-it-in-nodejs-2/)
[AWS HMAC Signature](https://docs.aws.amazon.com/AWSECommerceService/latest/DG/HMACSignatures.html)