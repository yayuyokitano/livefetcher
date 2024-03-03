SHELL=/bin/bash -o pipefail
export CGO_LDFLAGS := $(shell mecab-config --libs)
export CGO_CFLAGS := -I$(shell mecab-config --inc-dir)

default:
	@echo "Please specify an argument."

watch:
	gow -e=go,mod,html,css,toml run ./cmd/livefetcher start

migrate:
	go run . migrate

runtest:
	go test -v ./... -timeout 120s

run:
	set -e
	CONTAINERIZED=true ./livefetcher migrate
	# TESTING=true CONTAINERIZED=true go test -v ./... -timeout 120s
	TESTING=false CONTAINERIZED=true ./livefetcher start