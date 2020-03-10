# base image
FROM golang:1.13.3-alpine3.10

# install system utils
RUN apk add --no-cache --update git build-base openssh-client

# change workdir
WORKDIR /go/src/api

# read build arguments
ARG WEB_PRIVATE_KEY
ARG GIT_DOMAIN

# configuring ssh client
RUN mkdir ~/.ssh && \
    echo "$WEB_PRIVATE_KEY" | tr -d '\r' > ~/.ssh/id_rsa && \
    chmod 600 ~/.ssh/id_rsa && \
    ssh-keyscan -H $GIT_DOMAIN >> ~/.ssh/known_hosts

# copy project
COPY . .

# install project dependecies
RUN go get -v ./... &&\
    # run tests
    go test -v ./... &&\
    # build application
    go build -o demo-echo .

CMD ["./demo-echo"]
