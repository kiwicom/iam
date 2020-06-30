FROM golang:1.12.9-alpine3.9

RUN apk add --no-cache \
    git~=2.20.4 \
    autoconf~=2.69 \
    libtool~=2.4.6 \
    curl~=7.64.0 \
    gcc~=8.3.0 \
    musl-dev~=1.1.20 \
    zstd~=1.3.8 \
    lz4~=1.8.3

RUN curl -sSL https://github.com/ovh/venom/releases/download/v0.26.0/venom.linux-amd64 -o /usr/local/bin/venom
RUN chmod u+x /usr/local/bin/venom

COPY . /venom
WORKDIR /venom

CMD ["venom", "run", "--var-from-file", "variables.yml", "--parallel", "5", "--format=xml", "--output-dir", "."]
