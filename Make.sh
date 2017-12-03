#!/bin/bash

PWD=$(pwd)
PROJECT="/go${PWD/"$GOPATH"/}"

case $1 in
    "x86")
        docker run --rm -v $GOPATH:/go -w $PROJECT \
          -e "CGO_ENABLED=1" -e "GOARCH=386" \
          go-x86:latest go build .
        ;;
    "armv6")
        docker run --rm -v $GOPATH:/go -w $PROJECT \
          -e "CGO_ENABLED=1" -e "GOARCH=arm" \
          armv6:latest go build .
        ;;
    *)
        echo "armv6 or x86"
esac
