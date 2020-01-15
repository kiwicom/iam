FROM golang:1.13.6 as builder
RUN mkdir /app
WORKDIR /app

# This step is done separately than `COPY . /app/` in order to
# cache dependencies.
COPY go.mod go.sum Makefile /app/
RUN make install_deps

COPY . /app/
ARG CI_COMMIT_SHORT_SHA
RUN make build

FROM alpine:3.11.2
RUN  apk add --no-cache --virtual=.run-deps ca-certificates &&\
  mkdir /app

WORKDIR /app
COPY --from=builder /app/main ./main

# README.md is not used by the app, but is needed for COPY to not fail in case
# .env.yaml or .well-known don't exist.
COPY --from=builder /app/README.md /app/.env.yaml* ./
COPY --from=builder /app/README.md /app/.well-known/* ./.well-known/
RUN rm -f ./README.md ./.well-known/README.md

EXPOSE 8080 8090

USER nobody
ENTRYPOINT ./main
