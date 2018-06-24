#!/bin/bash

release=$1
if [ -z "$release" ]; then
    echo "usage:   .scripts/build_crossplatform.sh <release name>"
    exit 1
fi

export GOPATH=$HOME/.go

echo go version

rm -rf dist
mkdir dist

export GOOS=linux
export GOARCH=386
go get -v -d ./
go build ./
mv osem_notify dist/osem_notify_${release}_linux32
export GOARCH=amd64
go get -v -d ./
go build ./
mv osem_notify dist/osem_notify_${release}_linux64

export GOOS=windows
export GOARCH=386
go get -v -d ./
go build ./
mv osem_notify.exe dist/osem_notify_${release}_win32.exe
export GOARCH=amd64
go get -v -d ./
go build ./
mv osem_notify.exe dist/osem_notify_${release}_win64.exe

export GOOS=darwin
export GOARCH=386
go get -v -d ./
go build ./
mv osem_notify dist/osem_notify_${release}_mac32
export GOARCH=amd64
go get -v -d ./
go build ./
mv osem_notify dist/osem_notify_${release}_mac64

# export GOOS=linux
# export GOARCH=arm
# go get -v -d ./
# go build ./
# mv osem_notify dist/osem_notify_${release}_linux_arm

# export GOOS=android
# export GOARCH=arm
# go get -v -d ./
# go build ./
# mv osem_notify dist/osem_notify_${release}_android
