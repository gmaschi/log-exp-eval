sqlc-gen:
	./scripts/codegen/sqlc/sqlc-generate.sh

mock-generate:
	go install github.com/golang/mock/mockgen@v1.6.0
	go generate ./...
	go mod download

########################################################################################################################

server:
	docker compose -f ./build/docker/server/docker-compose.yml up -d --build

server-down:
	docker compose -f ./build/docker/server/docker-compose.yml down -v

########################################################################################################################

unit-test:
	chmod +x ./scripts/tests/unit/unit-test.sh
	./scripts/tests/unit/unit-test.sh

integration-test:
	chmod +x ./scripts/tests/unit/unit-test.sh
	./scripts/tests/integration/datastore/postgresql/exp/integration-test.sh

test: unit-test integration-test

########################################################################################################################

swag-install:
	which swagger || (brew tap go-swagger/go-swagger && brew install go-swagger)

swag-generate:
	swagger generate spec -o ./internal/docs/generated/swagger.json --scan-models

swag-validate:
	swagger validate ./internal/docs/generated/swagger.json

swag-serve: swag-install swag-generate swag-validate
	swagger serve -F=swagger ./internal/docs/generated/swagger.json