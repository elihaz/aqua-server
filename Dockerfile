FROM golang:latest AS Builder

COPY . /go/src/app
WORKDIR /go/src/app

RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build *.go

FROM alpine:latest AS Runner
RUN apk add --update ca-certificates

COPY  --from=Builder /go/src/app/main  /usr/local/bin/main
COPY ./server.crt /usr/local/share/ca-certificates/server.crt
COPY ./server.key /usr/local/share/ca-certificates/server.key

ENV SERVER_CRT /usr/local/share/ca-certificates/server.crt
ENV SERVER_KEY /usr/local/share/ca-certificates/server.key

ENTRYPOINT ["main"]

