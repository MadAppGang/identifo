# Settings data structure

## Description

This section will describe all possible settings in `settings.yaml` config file.

This settings file is subject for changes and extendability.



## General 

| Field | Description |
| :--- | :--- |
| port | external port, exposed globally by load balancers and/or reverse proxy  |
| host | Identifo server URL. env variable `HOST_NAME` overrides this value from the config file. The host should have full URL, including scheme, hostname, path and port. |
| issuer | JWT token issuer, used as `iss` field value in JWT token. [Please refer to RFC7519 Section 4.1.1.](https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.1) |
| algorithm | Key signature algorithms for JWT tokens. [Please refer RFC7518 for details.](https://datatracker.ietf.org/doc/html/rfc7518) Supported options are: `es256`, `es256` or `auto`. Auto option will use keys algorithm as an option. |
| supported\_scopes | An array containing a list of the [OAuth 2.0](https://openid.net/specs/openid-connect-discovery-1_0.html#RFC6749) \[RFC6749\] scope values that this server supports. The server MUST support the `openid` scope value. Servers MAY choose not to advertise some supported scope values even when this parameter is used, although those defined in [\[OpenID.Core\]](https://openid.net/specs/openid-connect-discovery-1_0.html#OpenID.Core) SHOULD be listed, if supported. |

_Example:_

```yaml
general: 
  port: 8081
  host: http://localhost:8081 
  issuer: http://localhost:8081 
  algorithm: es256 
  supported_scopes:
    - openid
    - offline
    - admin_panel
    - ios
    - android
    - adv_manager
```

## Admin panel 

| Field | Description |
| :--- | :--- |
| loginEnvName | environment variable for admin account email address/login |
| passwordEnvName | environment variable for admin account password |

Example:

```yaml
adminAccount:
  loginEnvName: IDENTIFO_ADMIN_LOGIN
  passwordEnvName: IDENTIFO_ADMIN_PASSWORD
```

## Data storages

Storage settings hold together all storage settings. All settings for a particular database engine \(i.e, file paths for BoltDB, endpoints and regions for DynamoDB etc.\) are assumed to be the same across all stores. If they are not the same, the latest option in this file will be applied. For example, if there are two MongoDB-backed storage, `appStorage` and `tokenStorage`, and endpoint for `appStorage` is localhost:27017, while tokenStorage's endpoint is `localhost:27018`, the server will connect both stores to `localhost:27018`.

| Field | Description |
| :--- | :--- |
| appStorage | Application storage settings |
| userStorage | User accounts storage settings |
| tokenStorage | Tokens storage for all issues access and refresh tokens |
| tokenBlacklist | Storage for token blacklist |
| varificationCodeStorage | Storage to keep verification codes |
| inviteStorage | Storage for invitations for registration |

Now we support a list of storage types out of the box. It is easy to add a new one, so please free to implement it and send PR. And we have a plugin system, that will allow you to extend  the storage with custom logic on your favourite language with supported by [the Hashicorp plugin system](https://pkg.go.dev/github.com/hashicorp/go-plugin): Nodejs, python, RoR and any other language, which support gRPC.

Example:

```yaml
storage:
  appStorage: &storage_settings
    type: boltdb
    boltdb:
      path: ./db.db
  userStorage: *storage_settings
  tokenStorage: *storage_settings
  tokenBlacklist: *storage_settings
  verificationCodeStorage: *storage_settings
  inviteStorage: *storage_settings
```

Now we support the following types:

| Field 'type' value | Description |
| :--- | :--- |
| mongodb | MongoDB 4+ databases, you can use AtlasDB with a large free storage allowance to start. |
| dynamodb | AWS DynamoDB storage |
| boltDB | BoltDB local storage for simple solutions and single instance solutions |
| mem | In-memory storage for testing and development |

### MongoDB

| Field | Description |
| :--- | :--- |
| type | mongodb |
| mongo | object field to keep all settings for mongodb |
| mongo.database | the database name to keep all the data |
| mongo.connection | the connection string for cluster or single instance |

Example:

```yaml
storage:
  appStorage: &storage_settings
    type: mongodb
    mongo:
      database: identifo-test
      connection: mongodb://localhost:27017
```

### BoltDB

| Field | Description |
| :--- | :--- |
| type | boltdb |
| boltdb | object field to keep all settings for boltdb |
| boltdb.path | Full file path and name for boltdb file, could be absolute or relevant on local or network attached drive. |

Example:

```yaml
storage:
  appStorage: &storage_settings
    type: boltdb
    boltdb:
      path: ./db.db
```

### Memory

| Field | Description |
| :--- | :--- |
| type | fake |

Example:

```yaml
storage:
  appStorage:
    type: fake
```

### DynamoDb

| Field | Description |
| :--- | :--- |
| type | dynamodb |
| dynamo | Field to store all the relevant settings for dynamoDB |
| dynamo.endpoint | Full endpoint for DynamoDB |
| dynamo.region | Region for DynamoDB endpoint |

Example:

```yaml
storage:
  appStorage: &dynamo_settings
    type: dynamodb
    dynamo:
      endpoint: http://localhost:8000
      region: us-east-2
```

## Session storage 

Session storage keeps sessions for admin panel. 

| Field | Description |
| :--- | :--- |
| sessionStorage | root field to hold the session values |
| type | options are `memory`, `redis` and `dynamodb` |
| sessionDuration | session duration in seconds |
| redis | is object key to store values for session settings in redis |
| redis.address | redis address  |
| redis.password | redis password, optional |
| redis.db | redis database name |
| dynamo | is an object key to store session settings for dynamodb storage type |
| dynamo.region | AWS region for dynamodb session storage  |
| dynamo.endpoint | AWS Dynamodb endpoint for session storage |

Example:

```yaml
# Storage for admin sessions.
sessionStorage:
  type: memory # Supported values are "memory", "redis", and "dynamodb".
  # Admin session duration in seconds.
  # This value specifies the maximum time of inactivity in the admin panel before asking to relogin.
  sessionDuration: 300

  # example for redis session storage
  # redis:
  #   address: http://localhost:2073
  #   password: redis_password
  #   db: admin_sessions

  # example for dynamo session storage
  # dynamo:
  #   region: us-east1
  #   endpoint: dynamo_endpoint
```

## Key storage

Storage for keys used for signing and verifying JWT tokens. Technically, a private key is enough to generate the public key. But we are using both for convenience.

Currently we support keys from local file system and S3. Other options could be added in the future. like: base64 encoded env variable or etcd or Hashicorp vault or AWS KMS.

| Field | Description |
| :--- | :--- |
| keyStorage | root key for key settings |
| type | `local` or `s3` |
| file | key for local file type settings |
| file.private\_key\_path | absolute or relevant path for a private key |
| file.public\_key\_path | absolute or relevant path for a public key |
| s3 | key for s3 file type settings |
| s3.region | AWS S3 region |
| s3.bucket | AWS S3 bucket  |
| s3.public\_key\_key | Public key S3 key |
| s3.private\_key\_key | Private key S3 key |

Example:

```yaml
keyStorage: # Storage for keys used for signing and verifying JWTs.
  type: local # Key storage type. Supported values are "local" and "s3".
  #file/local key storage settings
  file:
    private_key_path: ./jwt/test_artifacts/private.pem
    public_key_path: ./jwt/test_artifacts/public.pem
```

## Static files serving

Static file storage is responsible to store the admin panel, login pages and email templates.

Now we support only one storage for all of those above.

It tries to get the custom static data from the settings here, and if it failed, it gets the data included locally in docker file.

| Field | Description |
| :--- | :--- |
| static | root key for static data storage settings |
| type | storage type, supported values are `local`, `s3`, `dynamodb` |
| serveAdminPanel | boolean value, if `true` - serves admin panel |
| local | key for local static files settings |
| local.folder | folder for static data  |
| s3 | key for aws s3 static files settings |
| s3.region | aws s3 region for s3 static files settings |
| s3.bucket | aws s3 bucket for static file storage |
| s3.folder | aws s3 key/folder/prefix for static file storage |
| dynamo | key for aws dynamodb static files settings |
| dynamo.region | aws dynamodb region |
| dynamo.endpoint | aws dynamodb enpoint |

Example:

```yaml
static:
  type: local
  serveAdminPanel: true
  local:
    folder: ./static
  # s3 storage example
  # s3:
  #   region: east-us1
  #   bucket: identifo-data
  #   folder: static
  # dynamodb storage example
  # dynamo:
  #   region: east-us1
  #   endpoint: dbb-endpoint
```

## Login options

Settings for login options supported by identifo.

| Field | Description |
| :--- | :--- |
| login | root key for login settings |
| loginWith | key for login types  |
| loginWith.phone | boolean value, indicating login with phone is supported |
| loginWith.email | boolean value, login with email is supported |
| loginWith.username | boolean value, login with username is supported |
| loginWith.federated | boolean value, federated login is supported |
| tfyType | Two-factor authentication, currently we support `app`, `sms` and `email`. |

Example:

```yaml
login: # Supported login ways.
  loginWith:
    phone: true
    email: true
    username: true
    federated: true
  # Type of two-factor authentication, if application enables it.
  # Supported values are: "app" (like Google Authenticator), "sms", "email".
  tfaType: app
```

## External services and integrations

Now we supporting two types of external services, for sending SMS and Emails.

`services` is a root key for external services settings.

All services are supporting `mock` type. This type prints everything to a console, instead of sending it somewhere. Use it for development purposes.

### Email external service

| Field | Description |
| :--- | :--- |
| email | root key for email service settings |
| email.type | email service type, supported types are mailgun, ses and mock. Mock types does not requeres any  |
| email.mailgun | mailgun email service settings key |
| email.mailgun.domain | mailgun domain name |
| email.mailgun.privateKey | mailgun private key  |
| email.mailgun.publicKey | mailgun public key |
| email.mailgun.sender | sender email address |
| email.ses | AWS SES email settings key |
| email.ses.sender | email sender for SES service |
| email.ses.region | email AWS SES region |

Example:

```yaml
services:
  email: # Email service settings.
    type: mock # Supported values are "mailgun", "aws ses", and "mock".
    # mailgun:
    #   domain: identifo.com 
    #   privateKey: ABXCDS 
    #   publicKey: AAABBBDDD 
    #   sender: admin@admin.com 
    # ses:
    #   sender: admin@admin.com 
    #   region: es-east1 
```

### SMS external service

| Field | Description |
| :--- | :--- |
| sms | root key for sms settings |
| sms.type | SMS services type, now we support `mock`, `twilio`, `nexmo`, `routmobile` |
| sms.twilio | key to store settings for Twilio SMS service |
| sms.twilio.accountSid | Twilio account SID |
| sms.twilio.authToken | Twilio authentication token |
| sms.twilio.serviceSid | Twilio service SID |
| sms.nexmo | key to store nexmo service settings |
| sms.nexmo.apiKey | Nexmo API Key |
| sms.nexmo.apiSecret | Nexmo API Secret |
| sms.routemobile | RouteMobile service settings |
| sms.routemobile.username | Routemobile username |
| sms.routemobile.password | Routemobile service password |
| sms.routemobile.source | Routemobile service source |
| sms.routemobile.region | Routemobile region settings |

```yaml
services:
  sms: # SMS service settings.
    type: mock # Supported values are: "twilio", "nexmo", "routemobile", "mock".
    # twilio:
    #   accountSid: SID1234 
    #   authToken: TOKENABCDS 
    #   serviceSid: SIDFFFF 
    # nexmo:
    #   apiKey: KEY1234 
    #   apiSecret: SECRET4433 
    # routemobile:
    #   username: identifo 
    #   password: secret 
    #   source: whatever 
    #   region: australia 
```

