SHELL=/bin/bash -o pipefail

default:
	@echo "Please specify an argument."

watch:
	gow -e=go,mod,html,css,toml run ./cmd/livefetcher start

migrate:
	go run . migrate

run:
	set -e
	CONTAINERIZED=true ./livefetcher migrate
	# TESTING=true CONTAINERIZED=true go test -v ./... -timeout 120s
	TESTING=false CONTAINERIZED=true ./livefetcher start