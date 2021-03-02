FROM golangci/golangci-lint:v1.36

WORKDIR /app

ARG COMPOSE_USER_ID
ARG COMPOSE_GROUP_ID

RUN addgroup --gid $COMPOSE_GROUP_ID xdocker \
    && useradd -ms /bin/bash --gid $COMPOSE_GROUP_ID -u $COMPOSE_USER_ID xdocker \
    && adduser xdocker sudo \
    && echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers \
    && usermod -a -G www-data xdocker \
    && usermod -a -G root xdocker

USER xdocker