#!/bin/bash

sudo apt-get install zip
zip noel.zip noel.exe

mkdir dist
mv noel.zip dist/noel_$GIT_BRANCH.zip

if [ "$GIT_BRANCH" = "master" ]; then
	GIT_BRANCH=`git name-rev --name-only HEAD`
	go get github.com/tcnksm/ghr
	ghr --username $GITHUB_USERNAME --token $GITHUB_TOKEN --replace --prerelease --debug auto-release-$GIT_BRANCH dist/
fi
