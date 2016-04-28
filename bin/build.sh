
set -x 
set -o 
set -e

for i in ../gateway/gateway.go  ../monitorsystem/monitorclient/monitorclient.go  ../router/router.go ../poster/poster.go ../localposter/localposter.go ../cacheserver/cacheserver.go  ../monitorsystem/monitorserver/monitorserver.go ../client/client.go ../parselog/parselog.go;
do
go build -ldflags "-X github.com/fakewechat/lib/version.Buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X github.com/fakewechat/lib/version.Githash=`git rev-parse HEAD` -X github.com/fakewechat/lib/version.ProgramVersion=0.1"  $i
done
