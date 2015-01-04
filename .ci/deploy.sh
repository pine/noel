#!/bin/bash

sudo apt-get install zip
zip noel.zip noel.exe

GIT_BRANCH=`git name-rev --name-only HEAD`
GIT_REV=`git rev-parse --short HEAD`

mkdir dist
mv noel.exe dist/noel_$GIT_BRANCH.exe
mv noel.zip dist/noel_$GIT_BRANCH.zip

if [ "$GIT_BRANCH" = "master" ]; then
	go get github.com/tcnksm/ghr
	ghr --username $GITHUB_USERNAME --token $GITHUB_TOKEN --replace --prerelease --debug auto-release-$GIT_BRANCH-$GIT_REV dist/
fi
