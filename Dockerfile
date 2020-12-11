# Start from golang base image
FROM golang:1.13.3-alpine3.10 as builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

# Set the current working directory inside the container
WORKDIR /app

RUN go get github.com/githubnemo/CompileDaemon
RUN go get github.com/swaggo/swag/cmd/swag

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.7.3/wait /wait
RUN chmod +x /wait

#Command to run the executable
CMD swag init -g cmd/main.go\
  && /wait \
  && go run migrations/entry.go --verbose \
  && CompileDaemon --build="go build cmd/main.go"  --command="./main" --color
