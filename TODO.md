# TODO refactoring plan

## Config refactoring

- [x] refactor server composer
- [x] implement config files instead fo set of main files
- [x] refactor config file from flat to tree structure
- [x] refactor session service
- [x] refactor token service model
- [x] refactor key storage
- [x] refactor token service creation with configurator
- [x] refactor server creation
- [x] Check crash scenarious
- [x] release V2 beta branch
- [x] implement dump data import
- [x] implement app setting to validate HMAC signature (Web apps disabled by default)
- [x] implement integration testing
- [x] html/routes.go - check for we need static files handler?
- [x] check config file change monitoring
- [x] check S3 config file support
- [x] check S3 config file reloading
- [x] check template data deployment
  - [x] from static data
  - [x] from S3
  - [x] from dynamodb
- [x] refactor and fix jwt/token_test.go
- [x] implement initializer with JWKS URL .well-known


## User model and auth refactoring
- [x] New user model
- [ ] New app model
- [ ] new memory storage for user model
- [ ] new boltdb storage for user model
- [ ] new dynamodb storaGE for user model
- [ ] new grpc storage fro user model
- [ ] new mongodb storage for user model
- [ ] new plugin storage for user model
- [ ] new rest storage fro user model