# base image
FROM golang:1.21.4-alpine3.18 as build

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
    go install github.com/swaggo/swag/cmd/swag@v1.8.10 &&\
    $GOPATH/bin/swag init -g cmd/main.go &&\
    go get -v ./... &&\
    go build -o demo-echo ./cmd &&\
    go build -o entry ./migrations

FROM alpine

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.7.3/wait /wait
RUN chmod +x /wait

WORKDIR /app
COPY --from=build /go/src/api/demo-echo /app/demo-echo
COPY --from=build /go/src/api/entry /app/entry

RUN ["chmod", "+x", "./demo-echo"]
RUN ["chmod", "+x", "./entry"]

CMD /wait &&\
    ./entry --verbose &&\
    ./demo-echo
