###
# Mainflux Dockerfile
###

FROM golang:alpine
MAINTAINER Mainflux

###
# Install
###

RUN apk update && apk add git && rm -rf /var/cache/apk/*

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/mainflux/mainflux-core-server

RUN mkdir -p /config/core
COPY config/config-docker.yml /config/core/config.yml

# Get and install the dependencies
RUN go get github.com/mainflux/mainflux-core-server

###
# Run main command from entrypoint and parameters in CMD[]
###
CMD ["/config/core/config.yml"]

# Run mainflux command by default when the container starts.
ENTRYPOINT /go/bin/mainflux-core-server

