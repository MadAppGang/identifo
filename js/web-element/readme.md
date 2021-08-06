# Identifo Web Element
Simple drop in auth web component for Identifo.

## Usage
You need to set `app-id` and `url` props.

```<identifo-form app-id="ххх" url="https://identifo.com" theme="light"></identifo-form>```

## Properties

| Property | Attribute | Description | Type                                                                                                                                    | Default     |
| -------- | --------- | ----------- | --------------------------------------------------------------------------------------------------------------------------------------- | ----------- |
| `appId`  | `app-id`  |             | `string`                                                                                                                                | `undefined` |
| `route`  | `route`   |             | `"callback" \| "error" \| "login" \| "otp/login" \| "password/forgot" \| "password/reset" \| "register" \| "tfa/setup" \| "tfa/verify"` | `'login'`   |
| `theme`  | `theme`   |             | `"dark" \| "light"`                                                                                                                     | `undefined` |
| `token`  | `token`   |             | `string`                                                                                                                                | `''`        |
| `url`    | `url`     |             | `string`                                                                                                                                | `undefined` |


## Events

| Event      | Description | Type                    |
| ---------- | ----------- | ----------------------- |
| `complete` |             | `CustomEvent<string>`   |
| `error`    |             | `CustomEvent<ApiError>` |


----------------------------------------------