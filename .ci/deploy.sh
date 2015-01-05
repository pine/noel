#!/bin/bash

mkdir -p log

sudo apt-get install zip > log/apt-get_zip.log
zip noel.zip noel.exe

GIT_BRANCH=`git name-rev --name-only HEAD`
GIT_REV=`git rev-parse --short HEAD`

echo $GIT_BRANCH
echo $GIT_REV

mkdir dist
mv noel.exe dist/noel.exe
mv noel.zip dist/noel.zip

if [ "$GIT_BRANCH" = "master" ]; then
	go get github.com/tcnksm/ghr
	ghr --username $GITHUB_USERNAME --token $GITHUB_TOKEN --replace --prerelease --debug auto-release-$GIT_BRANCH-$GIT_REV dist/
fi
