# Identifo.js
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[![https://nodei.co/npm/@identifo/identifo-auth-js.png?downloads=true&downloadRank=true&stars=true](https://nodei.co/npm/@identifo/identifo-auth-js.png?downloads=true&downloadRank=true&stars=true)](https://www.npmjs.com/package/@identifo/identifo-auth-js)

Browser library for authentication through [Identifo](https://github.com/madappGang/identifo).

## Install
```bash
# Run this command in your project root folder.
# yarn
$ yarn add @identifo/identifo-auth-js

# npm
$ npm install @identifo/identifo-auth-js
```

## Usage

### Initialize
```javascript
  import IdentifoAuth from '@identifo/identifo-auth-js';

  const identifo = new IdentifoAuth({
    url: 'http://localhost:8081', // URI of your Identifo server.
    appId: 'your app ID', // ID of application that you want to get access to.
  });
```
### Parametrs
All parameters can be considered optional unless otherwise stated.
Option | Type | Description
--- | --- | ---
`url` | string (required) | The url of your Identifo server.
`appId` | string (required) | ID of application that you want to get access to.
`issuer` | string | The issuer claim identifies the principal that issued the JWT.
`scopes` | string [] | The default scopes used for all authorization requests.
`redirectUri` | string | The URL where Identifo will call back to with the result of a successful authentication. It must be added to the "Redirect URLs" in your Identifo Application's settings.
`postLogoutRedirectUri` | string | The URL where Identifo will call back to after a successful log-out. It must be added to the "Redirect URLs" in your Identifo Application's settings.
`tokenManager`| TokenManager instance | The tokenManager provides access to client storage for specific purposes
`autoRenew`| boolean | Renew tokens before they expire. By default autoRenew is false.

### API
#### init()
Method which initialize your identifo instance and actualize your states. Call it before other methods. 
#### handleAuthentication():Promise<boolean>
Processing the received token. Call it on the redirect page to handle and verify token. 
#### renewSession():Promise<string>
Allows you to renew session manually. Returns new JWT.
```javascript
identifo.renewSession().then((token) => {})
```
#### signin() 
Redirect user to sign-in page. After successful sign-in you will be redirected to `redirectUri`
#### signup() 
Redirect user to sign-up page. After successful sign-up you will be redirected to `redirectUri`
#### logout() 
Logout and redirects back to `postLogoutRedirectUri` the user out of their current Identifo session and clears all tokens stored locally in the TokenManager. By default, the refresh token (if any) and access token are revoked so they can no longer be used. 
#### getToken()
 Returns token and token payload 
```javascript
{ token: 'JWT', payload: 'JWT payload' }
```
#### TokenManager
```javascript
import IdentifoAuth, { SessionStorageManager } from '@identifo/identifo-auth-js';

const identifo = new IdentifoAuth({
    appId: 'your app ID',
    url: 'http://localhost:8081',
    tokenManager: new SessionStorageManager() // you can pass your custom key for storage. Bt default it`s identifo_access_token
  })
```
You can import SessionStorageManager | LocalStorageManager. By default it`s LocalStorageManager.
#### isAuth:boolean
Returns actual auth statis of the user.
```javascript
const isAuth = identifo.isAuth;
```
### redirectUri

The url that is redirected to when using `identifo.signin` or `identifo signup` methods. This must be listed in your Identifo application's redirect URLs. If no redirectUri is provided, defaults to the current url (window.location.href). 

### postLogoutRedirectUri

The url that is redirected to when using `identifo.logout` method. This must be listed in your Identifo application's redirect URLs. If not specified, user will be redirected to login page.

## Author

[madappgang](https://madappgang.com)

## License

This project is licensed under the MIT license. See the [LICENSE](LICENSE) file for more info.
