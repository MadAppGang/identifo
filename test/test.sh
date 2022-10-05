#!/bin/sh
# setup test enviroment and run tests

export IDENTIFO_TEST_INTEGRATION=1
export IDENTIFO_TEST_AWS_ENDPOINT="http://localhost:5001"
export AWS_ACCESS_KEY_ID='testing'
export AWS_SECRET_ACCESS_KEY='testing'
export AWS_SECURITY_TOKEN='testing'
export AWS_SESSION_TOKEN='testing'


export IDENTIFO_STORAGE_MONGO_TEST_INTEGRATION=1
export IDENTIFO_STORAGE_MONGO_CONN="mongodb://admin:password@localhost:27017/billing-local?authSource=admin"
export IDENTIFO_REDIS_HOST="127.0.0.1:6379"

docker-compose up -d

sleep 1
echo "dependencies started"

go test -race -timeout=60s -count=1 ../...
test_exit=$?

# docker-compose down -v
docker-compose rm -s -f -v

exit $test_exit
