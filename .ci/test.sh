#!/bin/bash

echo "> go get -d ./..."
go get -d ./...

echo "> go test -v"
go test -v
