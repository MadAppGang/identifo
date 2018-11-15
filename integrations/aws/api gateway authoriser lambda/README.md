# Identifo APU Authorizer

`identifo_api_auth` is an Amazon API Gateway Lambda authorizer. The task is simple:

- extract the bearer token from the header `Authentication: Bearer XXX...`
- validate token signature and claims (date, issuer)
- check token in the blacklist
- return  the claims as a JSON

For full documentation, please refer AWS documentation:
[AWS Lambda authorizer docs](https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-use-lambda-authorizer.html)

## Dependency manager

All dependencies are not included in git index. Before building the project, please update all dependencies. We are using standard dependency manager `dep`. For more details please refer to [official github repository](https://github.com/golang/dep).

To update all dependencies, just type:
`dep unsure -update`

or just type:
`make update`

## Local debug

To debug and run the project locally, you need Docker and AWS SAM. Please refer [official documentation](https://docs.aws.amazon.com/lambda/latest/dg/serverless_app.html) about it.
To run test version, type:

`make debug`

Other option is to use tests, basically this option is more preferable. Try to use unit tests for debug and AWS SAM for fine tunning and final functional and manual tests only.
To run all the tests, type:

`make test`
