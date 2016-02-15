

nohup ./router  >/dev/null &
nohup ./monitorserver  >/dev/null &
nohup  ./cacheserver -hostid=ca1 -servertype=cache -listenaddress="127.0.0.1" -listenport=9601 -rgid=1 -process=ca -routeserverurl="http://127.0.0.1:8080/server/"     -routeregisturl="http://127.0.0.1:8080/regist/"  >/dev/null &
nohup  ./cacheserver -hostid=ca2 -servertype=cache -listenaddress="127.0.0.1" -listenport=9602 -rgid=2 -process=ca -routeserverurl="http://127.0.0.1:8080/server/"     -routeregisturl="http://127.0.0.1:8080/regist/"  >/dev/null &
nohup  ./cacheserver -hostid=ca3 -servertype=cache -listenaddress="127.0.0.1" -listenport=9603 -rgid=3 -process=ca -routeserverurl="http://127.0.0.1:8080/server/"     -routeregisturl="http://127.0.0.1:8080/regist/"  >/dev/null &
nohup  ./cacheserver -hostid=ca4 -servertype=cache -listenaddress="127.0.0.1" -listenport=9604 -rgid=4 -process=ca -routeserverurl="http://127.0.0.1:8080/server/"     -routeregisturl="http://127.0.0.1:8080/regist/"  >/dev/null &




nohup  ./gateway -hostid=gw5 -servertype=gw -listenaddress="127.0.0.1" -listenport=9501 -rgid=1 -process=ca -routeserverurl="http://127.0.0.1:8080/server/"     -routeregisturl="http://127.0.0.1:8080/regist/"  >/dev/null &
nohup  ./gateway -hostid=gw6 -servertype=gw -listenaddress="127.0.0.1" -listenport=9502 -rgid=2 -process=ca -routeserverurl="http://127.0.0.1:8080/server/"     -routeregisturl="http://127.0.0.1:8080/regist/"  >/dev/null &
nohup  ./gateway -hostid=gw7 -servertype=gw -listenaddress="127.0.0.1" -listenport=9503 -rgid=3 -process=ca -routeserverurl="http://127.0.0.1:8080/server/"     -routeregisturl="http://127.0.0.1:8080/regist/"  >/dev/null &
nohup  ./gateway -hostid=gw8 -servertype=gw -listenaddress="127.0.0.1" -listenport=9504 -rgid=4 -process=ca -routeserverurl="http://127.0.0.1:8080/server/"     -routeregisturl="http://127.0.0.1:8080/regist/"  >/dev/null &




nohup  ./localposter -hostid=ca9 -servertype=localposter -listenaddress="127.0.0.1" -listenport=9701 -rgid=1 -process=ca -routeserverurl="http://127.0.0.1:8080/server/"     -routeregisturl="http://127.0.0.1:8080/regist/"  >/dev/null &
nohup  ./localposter -hostid=ca10 -servertype=localposter -listenaddress="127.0.0.1" -listenport=9702 -rgid=2 -process=ca -routeserverurl="http://127.0.0.1:8080/server/"     -routeregisturl="http://127.0.0.1:8080/regist/"  >/dev/null &
nohup  ./localposter -hostid=ca11 -servertype=localposter -listenaddress="127.0.0.1" -listenport=9703 -rgid=3 -process=ca -routeserverurl="http://127.0.0.1:8080/server/"     -routeregisturl="http://127.0.0.1:8080/regist/"  >/dev/null &
nohup  ./localposter -hostid=ca12 -servertype=localposter -listenaddress="127.0.0.1" -listenport=9704 -rgid=4 -process=ca -routeserverurl="http://127.0.0.1:8080/server/"     -routeregisturl="http://127.0.0.1:8080/regist/"  >/dev/null &




nohup  ./poster -hostid=poster13 -servertype=poster -listenaddress="127.0.0.1" -listenport=9801 -rgid=1 -process=ca -routeserverurl="http://127.0.0.1:8080/server/"     -routeregisturl="http://127.0.0.1:8080/regist/"  >/dev/null &
nohup  ./poster -hostid=poster14 -servertype=poster -listenaddress="127.0.0.1" -listenport=9802 -rgid=2 -process=ca -routeserverurl="http://127.0.0.1:8080/server/"     -routeregisturl="http://127.0.0.1:8080/regist/"  >/dev/null &
nohup  ./poster -hostid=poster15 -servertype=poster -listenaddress="127.0.0.1" -listenport=9803 -rgid=3 -process=ca -routeserverurl="http://127.0.0.1:8080/server/"     -routeregisturl="http://127.0.0.1:8080/regist/"  >/dev/null &
nohup  ./poster -hostid=poster16 -servertype=poster -listenaddress="127.0.0.1" -listenport=9804 -rgid=4 -process=ca -routeserverurl="http://127.0.0.1:8080/server/"     -routeregisturl="http://127.0.0.1:8080/regist/"  >/dev/null &



redis-cli -p  1501 flushall
redis-cli -p  1502 flushall
redis-cli -p  1503 flushall
redis-cli -p  1504 flushall



./savetoredis.py -m 127.0.0.1 -p 1501 < 1.txt
./savetoredis.py -m 127.0.0.1 -p 1502 < 2.txt
./savetoredis.py -m 127.0.0.1 -p 1503 < 3.txt
./savetoredis.py -m 127.0.0.1 -p 1504 < 4.txt


curl -kvv "http://127.0.0.1:8080/regist/redis/1/?id=res1&host=127.0.0.1&port=1501&cellid=1"
curl -kvv "http://127.0.0.1:8080/regist/redis/2/?id=res2&host=127.0.0.1&port=1502&cellid=1"
curl -kvv "http://127.0.0.1:8080/regist/redis/3/?id=res3&host=127.0.0.1&port=1503&cellid=1"
curl -kvv "http://127.0.0.1:8080/regist/redis/4/?id=res4&host=127.0.0.1&port=1504&cellid=1"




