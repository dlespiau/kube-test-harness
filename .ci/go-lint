#!/bin/bash

set -e

if [ ! $(command -v gometalinter) ]
then
    go get github.com/alecthomas/gometalinter
    gometalinter --install --vendor
fi


gometalinter --tests --vendor --disable-all --deadline=600s \
    --enable=misspell \
    --enable=vet \
    --enable=ineffassign \
    --enable=gofmt \
    --enable=gocyclo --cyclo-over=20 \
    --enable=golint \
    ./...
