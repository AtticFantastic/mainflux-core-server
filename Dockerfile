###
# Mainflux Dockerfile
###
# Set the base image to Node, onbuild variant: https://registry.hub.docker.com/_/node/

FROM golang:alpine
MAINTAINER Mainflux

ENV MAINFLUX_NATS_PORT=4222

###
# Install
###

RUN apk update && apk add git && rm -rf /var/cache/apk/*

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/mainflux/mainflux-core-server

# Get and install the dependencies
RUN go get github.com/mainflux/mainflux-core-server

###
# Run main command from entrypoint and parameters in CMD[]
###
CMD [""]

# Run mainflux command by default when the container starts.
ENTRYPOINT /go/bin/mainflux-core-server

EXPOSE $MAINFLUX_NATS_PORT
