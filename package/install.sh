#!/bin/bash
set -o -e

yum install -y gcc gcc-c++ autoconf automake libtool vim  wget

cd /home/ec2-user/gopath/src/github.com/fakewechat/package

#golang install
curl -O -L https://storage.googleapis.com/golang/go1.6.2.linux-amd64.tar.gz -kvv
tar zxvf go1.6.2.linux-amd64.tar.gz
mv go /home/ec2-user/

#setup env
echo -e "export GOROOT=/home/ec2-user/go\n export PATH=$PATH:/home/ec2-user/go/bin\n export GOPATH=/home/ec2-user/gopath" >> /home/ec2-user/.bashrc
. /home/ec2-user/.bashrc

#install protbuf
cd /home/ec2-user/gopath/src/github.com/fakewechat/package
tar zxvf protobuf-master.tar.gz
cd protobuf-master
./autogen.sh
./configure 
make
cd python
sudo python setup.py install



# build package
cd /home/ec2-user/gopath/src/github.com/fakewechat/bin
./build.sh

#install redis for python
wget https://bootstrap.pypa.io/ez_setup.py -O - | sudo python
sudo easy_install redis

# redis
cd /home/ec2-user/gopath/src/github.com/fakewechat/package
tar zxvf redis-2.8.24.tar.gz
cd redis-2.8.24
make
sudo make install
cp src/redis-server /home/ec2-user/gopath/src/github.com/fakewechat/bin



rpm -Uvh http://dl.fedoraproject.org/pub/epel/7/x86_64/e/epel-release-7-6.noarch.rpm
rpm --import /etc/pki/rpm-gpg/RPM-GPG-KEY-EPEL-7

yum install -y ansible 





