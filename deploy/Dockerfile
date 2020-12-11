# base image
FROM golang:1.13.3-alpine3.10 as build

RUN apk add --no-cache --update git build-base openssh-client

WORKDIR /go/src/api

ARG WEB_PRIVATE_KEY
ARG GIT_DOMAIN

RUN mkdir ~/.ssh && \
    echo "$WEB_PRIVATE_KEY" | tr -d '\r' > ~/.ssh/id_rsa && \
    chmod 600 ~/.ssh/id_rsa && \
    ssh-keyscan -H $GIT_DOMAIN >> ~/.ssh/known_hosts

COPY . .

RUN git config --global http.sslVerify true &&\
    go get -v github.com/swaggo/swag/cmd/swag &&\
    $GOPATH/bin/swag init -g cmd/main.go &&\
    go get -v ./... &&\
    go build -o demo-echo .

FROM alpine

WORKDIR /app
COPY --from=build /go/src/api/demo-echo /app/demo-echo

CMD ["./demo-echo"]
