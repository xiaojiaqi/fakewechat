#!/bin/bash

set -x 
set -o 
set -e

for i in ../gateway/gateway.go ../router/router.go ../poster/poster.go   ../goredis/goredis.go ../localposter/localposter.go ../cacheserver/cacheserver.go  ../cacheserver/cacheserver.go  ../monitorsystem/monitorclient/monitorclient.go   ../monitorsystem/monitorserver/monitorserver.go ../client/client.go ../parselog/parselog.go;
do
go build -ldflags "-X github.com/fakewechat/lib/version.Buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X github.com/fakewechat/lib/version.Githash=`git rev-parse HEAD`  -X github.com/fakewechat/lib/version.ProgramVersion=0.1 " $i
done
