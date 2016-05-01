#!/bin/bash
set -o -e


cd /home/ec2-user/gopath/src/github.com/fakewechat/bin

./listcm.py

cp cmd/*.* .

chmod +x *.sh


sudo chmod 777 /etc/ansible/hosts

cat cmd/hosts >> /etc/ansible/hosts

./gen.sh

./ansible.sh
#ansible all -a   "sh /home/ec2-user/bin/kill.sh"

ansible all -a  "sh /home/ec2-user/bin/stop_redis.sh"
ansible all -a   "sh /home/ec2-user/bin/start_redis.sh"

ansible all -a  "sh /home/ec2-user/bin/savetoredis.sh"

ansible 10.0.2.20 -a  "sh /home/ec2-user/bin/monitor.sh"

./resign.sh

ansible all -a "sh /home/ec2-user/bin/server.sh"


