#!/bin/sh
# setup test enviroment and run tests

export IDENTIFO_TEST_INGRATION=1
export IDENTIFO_TEST_AWS_ENDPOINT="http://localhost:5001"
export AWS_ACCESS_KEY_ID='testing'
export AWS_SECRET_ACCESS_KEY='testing'
export AWS_SECURITY_TOKEN='testing'
export AWS_SESSION_TOKEN='testing'

docker-compose up -d

sleep 1
echo "dependencies started"

go test -race -timeout=60s -count=1 ../...
test_exit=$?

docker-compose down -v
docker-compose rm -s -f -v

exit $test_exit
