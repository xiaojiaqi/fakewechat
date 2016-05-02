#!/bin/bash
set -o -e


cd $HOME/gopath/src/github.com/fakewechat/bin

./listcm.py

cp cmd/*.* .

chmod +x *.sh


sudo chmod 777 /etc/ansible/hosts

cat cmd/hosts >> /etc/ansible/hosts

./genData.sh

./ansible.sh

#ansible all -a   "sh $HOME/bin/kill.sh"

ansible all -a  "sh $HOME/bin/stopredis.sh"
ansible all -a   "sh $HOME/bin/startredis.sh"

ansible all -a  "sh $HOME/bin/gosaveDb.sh"

ansible 10.0.2.20 -a  "sh $HOME/bin/monitor.sh"

./resign.sh

ansible all -a "sh $HOME/bin/server.sh"


