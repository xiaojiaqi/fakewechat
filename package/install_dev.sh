#!/bin/bash
set -o -e

sudo yum install -y gcc gcc-c++ autoconf automake libtool vim  wget psmisc

cd $HOME/gopath/src/github.com/fakewechat/package

#install redis for python
wget https://bootstrap.pypa.io/ez_setup.py -O - | sudo python
sudo easy_install redis


#install protbuf
cd $HOME/gopath/src/github.com/fakewechat/package
tar zxvf protobuf-master.tar.gz
cd protobuf-master
./autogen.sh
./configure 
make
cd python
sudo python setup.py install






