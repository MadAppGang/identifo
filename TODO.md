# TODO refactoring plan

## Config refactoring

- [x] refactor server composer
- [x] implement config files instead fo set of main files
- [x] refactor config file from flat to tree structure
- [x] refactor session service
- [x] refactor token service model
- [ ] html/routes.go - check to we need static files handler?
- [ ] refactor token service creation with configurator
- [ ] refactor key storage
- [ ] refactor main file
- [ ] refactor server creation
- [ ] implement dump data import
- [ ] implement integration testing
- [ ] check config file change monitoring
- [ ] check S3 config file support
- [ ] check S3 config file reloading
- [ ] check template data deployment
  - [ ] from static data
  - [ ] from S3
  - [ ] from dynamodb
- [ ] refactor and fix jwt/token.test.go
- [ ] implement initializer with JWKS URL .well-known
