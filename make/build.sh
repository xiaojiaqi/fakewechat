#!/bin/bash

function dobuild()  {
if [ $1 == 'clean' ];
then
    echo "go clean"
    go clean
else
    echo "go build"
    gofmt -w *.go
    go build -ldflags "-X github.com/fakewechat/lib/version.Buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X github.com/fakewechat/lib/version.Githash=`git rev-parse HEAD`  -X github.com/fakewechat/lib/version.ProgramVersion=0.1 "
fi
}

function foreachd(){
    for file in $1/*
    do
        if [ -d $file ];
        then
            echo $file
            cd  $file
            dobuild $2
            cd -
            foreachd $file $2
        fi
    done
}

if [ $1 == 'clean' ];
  then
    a='clean'
else
    a='build'
fi
foreachd `pwd` $a
