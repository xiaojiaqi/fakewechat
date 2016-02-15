#!/usr/bin/env sh

while [ True ];
do


zabbix_sender -z 172.16.11.194 -k "client.KeepConnection" -o $RANDOM -vv -s "FakeweChat Server"
zabbix_sender -z 172.16.11.194 -k "longCon.KeepConnection" -o $RANDOM -vv -s "FakeweChat Server"

sleep 10
done
