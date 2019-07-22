FROM golang:1.12.7 as builder
RUN mkdir /app
COPY . /app/
WORKDIR /app
ARG CI_COMMIT_SHORT_SHA
RUN make build

FROM alpine:3.10.1
RUN  apk add --no-cache --virtual=.run-deps ca-certificates &&\
  mkdir /app

WORKDIR /app
COPY --from=builder /app/main ./main
COPY --from=builder /app/.env.yaml .env.yaml
COPY --from=builder /app/.well-known .well-known/
COPY --from=builder /app/configs config/
EXPOSE 8080 8090

USER nobody
ENTRYPOINT ./main
