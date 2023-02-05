#!/usr/bin/env bash

DOCKER_COMPOSE_TEST="./build/docker/server/test/docker-compose-test.yml"
EVAL_TEST_SERVER="eval-server-test"
EVAL_TEST_DATABASE="eval-database-test"
DB_TEST_USER="myuser"
DB_TEST_NAME="eval-test"

# start test server
docker compose -f $DOCKER_COMPOSE_TEST up -d --build

# wait for docker to start up
until docker exec -i $EVAL_TEST_DATABASE psql -U $DB_TEST_USER $DB_TEST_NAME<<EOF
EOF
do
    echo "Waiting for postgres to start..."
    sleep 1.5
done

# run tests
echo "RUNNING INTEGRATION TESTS"
go test $(go list ./... | grep tests/integration/) -v

# clean resources
echo "FINISHED RUNNING INTEGRATION TESTS. CLEANING RESOURCES"
docker compose -f $DOCKER_COMPOSE_TEST down -v
echo "RESOURCES CLEANED"
