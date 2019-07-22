GOLANGCI_LINT := $(shell command -v golangci-lint 2> /dev/null)

install_deps:
	go mod download
	GO111MODULE=off go get github.com/cespare/reflex
	GO111MODULE=off go get -u github.com/golang/protobuf/protoc-gen-go

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

test:
	$(call log_info,Run tests and check race conditions)
	# https://golang.org/doc/articles/race_detector.html
	go test -race -v ./... -cover
	$(call log_success,All tests succeeded)

go-mod-tidy:
	$(call log_info,Check that go.mod and go.sum don't contain any unnecessary dependency)
	go mod tidy -v
	git diff-index --quiet HEAD
	$(call log_success,Go mod check succeeded!)

test/ci: test go-mod-tidy

test/watch:
	reflex --start-service -r '\.go$$' make test

build:
	CGO_ENABLED=0 go build cmd/main.go
ifndef CI_COMMIT_SHORT_SHA
	$(call log_info,SENTRY_RELEASE not set)
else
	@printf "${color_green}SENTRY_RELEASE: ${CI_COMMIT_SHORT_SHA}${color_off}\n"
	@echo "SENTRY_RELEASE: ${CI_COMMIT_SHORT_SHA}" >> .env.yaml
endif

lint: ## Runs golangci-lint. It outputs to the code-climate json file in if CI is defined.
	$(call log_info, Running golangci-lint)
ifndef GOLANGCI_LINT
	@echo "Can\'t find executable of the golangci-lint. Make sure it is installed. See github.com\/golangci\/golangci-lint#install"
	@exit 1
endif
	golangci-lint run $(if $(CI),--out-format code-climate > gl-code-quality-report.json)
	$(call log_success,Linting with golangci-lint succeeded!)

coala:
	docker run -v "$(shell pwd)":/app -v /tmp/coala-cache:/cache --workdir=/app coala/base:0.11 coala -a -n -j 4

-include ./test/e2e/e2e.mk

proto:
	protoc --go_out=plugins=grpc:. ./api/grpc/v1/iam.proto