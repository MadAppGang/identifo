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

_Example:_

```yaml
general: 
  port: 8081
  host: http://localhost:8081 
  issuer: http://localhost:8081 
  algorithm: es256 
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

## Key storage

## Static files serving

## Login options

## External services and integrations



