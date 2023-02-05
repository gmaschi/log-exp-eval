FROM postgres:15.1-alpine3.17
WORKDIR /docker-entrypoint-initdb.d
COPY ./internal/services/datastore/postgresql/exp/schema/ .