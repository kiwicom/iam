FROM golang:1.12.3 as builder
RUN mkdir /app
COPY . /app/
WORKDIR /app
ARG CI_COMMIT_SHORT_SHA
RUN make build

FROM alpine:3.9.3
RUN  apk add --no-cache --virtual=.run-deps ca-certificates &&\
  mkdir /app

WORKDIR /app
COPY --from=builder /app/main ./main
COPY --from=builder /app/.env.yaml .env.yaml
EXPOSE 8080

USER nobody
ENTRYPOINT ./main
