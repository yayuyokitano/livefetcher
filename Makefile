SHELL=/bin/bash -o pipefail
export CGO_LDFLAGS := $(shell mecab-config --libs)
export CGO_CFLAGS := -I$(shell mecab-config --inc-dir)

default:
	@echo "Please specify an argument."

watch:
	gow -e=go,mod,html,css,toml run ./cmd/livefetcher start

runlocal:
	go run ./cmd/livefetcher start

migrate:
	go run ./cmd/livefetcher migrate

runtest:
	go test -v ./... -timeout 120s

testconnector:
	CONNECTOR_ID=$(c) go test -v ./internal/core/connectors

runconnector:
	CONNECTOR_ID=$(c) go run ./cmd/livefetcher test

run:
	set -e
	CONTAINERIZED=true ./livefetcher migrate
	# TESTING=true CONTAINERIZED=true go test -v ./... -timeout 120s
	TESTING=false CONTAINERIZED=true ./livefetcher start