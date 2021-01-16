FROM golang:1.14 as buildContainer

ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOOS=linux
ENV GOPATH=/

COPY . /src/mqtt-logger
WORKDIR /src/mqtt-logger

RUN go get ./... &&\
    go build -ldflags -s -a -installsuffix cgo -o mqtt-logger ./cmd/mqtt-logger/


FROM alpine

COPY --from=buildContainer /src/mqtt-logger/mqtt-logger /mqtt-logger

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

USER 10000:10000

ENTRYPOINT ["/mqtt-logger"]