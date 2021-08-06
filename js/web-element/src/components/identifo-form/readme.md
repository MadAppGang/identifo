# <identifo-form></identifo-form>



<!-- Auto Generated Below -->


## Properties

| Property                | Attribute                  | Description | Type                                                                                                                                                                                          | Default     |
| ----------------------- | -------------------------- | ----------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------- |
| `appId`                 | `app-id`                   |             | `string`                                                                                                                                                                                      | `undefined` |
| `callbackUrl`           | `callback-url`             |             | `string`                                                                                                                                                                                      | `undefined` |
| `debug`                 | `debug`                    |             | `boolean`                                                                                                                                                                                     | `undefined` |
| `federatedRedirectUrl`  | `federated-redirect-url`   |             | `string`                                                                                                                                                                                      | `undefined` |
| `postLogoutRedirectUri` | `post-logout-redirect-uri` |             | `string`                                                                                                                                                                                      | `undefined` |
| `route`                 | `route`                    |             | `"callback" \| "error" \| "loading" \| "login" \| "logout" \| "otp/login" \| "password/forgot" \| "password/forgot/success" \| "password/reset" \| "register" \| "tfa/setup" \| "tfa/verify"` | `'login'`   |
| `scopes`                | `scopes`                   |             | `string`                                                                                                                                                                                      | `''`        |
| `theme`                 | `theme`                    |             | `"dark" \| "light"`                                                                                                                                                                           | `'light'`   |
| `token`                 | `token`                    |             | `string`                                                                                                                                                                                      | `undefined` |
| `url`                   | `url`                      |             | `string`                                                                                                                                                                                      | `undefined` |


## Events

| Event      | Description | Type                         |
| ---------- | ----------- | ---------------------------- |
| `complete` |             | `CustomEvent<LoginResponse>` |
| `error`    |             | `CustomEvent<ApiError>`      |


----------------------------------------------

*Built with [StencilJS](https://stenciljs.com/)*
