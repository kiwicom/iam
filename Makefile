GOLANGCI_LINT := $(shell command -v golangci-lint 2> /dev/null)

install_deps:
	go mod download
	GO111MODULE=off go get github.com/cespare/reflex
	GO111MODULE=off go get -u github.com/golang/protobuf/protoc-gen-go
	GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger

start:
	go run cmd/main.go

dev:
	reflex --start-service -r '\.go$$' make start

# Colorful output
color_off = \033[0m
color_cyan = \033[1;36m
color_green = \033[1;32m

define log_info
	@printf "$(color_cyan)$(1)$(color_off)\n"
endef
define log_success
	@printf "$(color_green)$(1)$(color_off)\n"
endef

lint:
	$(call log_info, Running golangci-lint)
	golangci-lint run
	$(call log_success,Linting with golangci-lint succeeded!)

go-mod-tidy:
	$(call log_info,Check that go.mod and go.sum don't contain any unnecessary dependency)
	go mod tidy -v
	git diff-index --quiet HEAD
	$(call log_success,Go mod check succeeded!)

test:
	$(call log_info,Run tests and check race conditions)
	# https://golang.org/doc/articles/race_detector.html
	go test -race -v ./... -cover
	$(call log_success,All tests succeeded)

test/ci: test go-mod-tidy

test/watch:
	reflex --start-service -r '\.go$$' make test

test/all: test/go go-mod-tidy lint e2e

build:
	CGO_ENABLED=0 go build cmd/main.go
ifndef CI_COMMIT_SHORT_SHA
	$(call log_info,SENTRY_RELEASE not set)
else
	@printf "${color_green}SENTRY_RELEASE: ${CI_COMMIT_SHORT_SHA}${color_off}\n"
	@echo "SENTRY_RELEASE: ${CI_COMMIT_SHORT_SHA}" >> .env.yaml
endif

proto:
	protoc --go_out=plugins=grpc:. ./api/grpc/v1/kiwi_iamapi.proto

swagger-validate:
	swagger validate ./api/swagger.yml

swagger-serve:
	swagger serve ./api/swagger.yml

-include ./tests/e2e/e2e.mk
-include ./www/docs.mk