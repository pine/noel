#!/bin/bash

sudo apt-get install gcc-mingw-w64-x86-64

pushd ~/

curl -s -o go.tar.gz https://storage.googleapis.com/golang/go1.4.src.tar.gz
tar xfz go.tar.gz

export GOROOT="$PWD/go"
export PATH="$GOROOT/bin:$PATH"

cd $GOROOT/src
env CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC_FOR_TARGET="x86_64-w64-mingw32-gcc -fno-stack-protector -D_FORTIFY_SOURCE=0 -lssp" ./make.bash --no-clean

popd
