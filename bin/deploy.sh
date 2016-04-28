#!/bin/bash
set -o -e


cd /home/ec2-user/gopath/src/github.com/fakewechat/bin

./listcm.py

cp cmd/* .

chmod +x *.sh


chmd 777 /etc/ansible/hosts

cat cmd/hosts >> /etc/ansible/hosts


cat "" >> /etc/ansible/hosts


./ansible.sh


./gen.sh

ansible all -a -m raw "/home/ec2-user/bin/stop_redis.sh"
ansible all -a -m raw "/home/ec2-user/bin/start_redis.sh"

./savetoredis.sh
 
ansible 10.0.2.20 -a -m raw "/home/ec2-user/bin/monitor.sh"

./resign.sh

ansible all -a -m raw "/home/ec2-user/bin/server.sh"


