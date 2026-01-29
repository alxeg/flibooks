
FROM golang:1.24-alpine AS build

ARG FB2C_VERSION=v1.78.4

ARG TARGETARCH
ARG TARGETOS


WORKDIR     /build

RUN         apk add --update --no-cache ca-certificates git curl

ENV         GOBIN=/build/bin

ADD         . /build

RUN         mkdir -p /build/bin && cd /build && \
            go build -mod=vendor ./cmd/... && \
            curl -L "https://github.com/rupor-github/fb2converter/releases/download/${FB2C_VERSION}/fb2c-${TARGETOS}-${TARGETARCH}.zip" -o fb2c.zip && \
            unzip -d bin/ fb2c.zip && \
            rm -rf fb2c.zip

FROM alpine:3.23

WORKDIR /flibooks

COPY        --from=build /build/flibooks /flibooks/
COPY        --from=build /build/config/flibooks.properties /flibooks/
COPY        --from=build /build/bin/fb2c /flibooks/

EXPOSE      8000

ENTRYPOINT [ "/flibooks/flibooks" ]
