#!/usr/bin/env bash

go test $(go list ./... | grep -v tests/integration/ | grep -v /mocks | grep -v /mockgen)
