#!/bin/bash
set -o -e

sudo yum install -y screen gcc gcc-c++ autoconf automake libtool vim  wget psmisc yum-utils make

cd $HOME/gopath/src/github.com/fakewechat/package

sudo yumdownloader  psmisc

cp psmisc*.rpm $HOME/gopath/src/github.com/fakewechat/bin/

#golang install
curl -O -L https://storage.googleapis.com/golang/go1.6.2.linux-amd64.tar.gz -kvv
tar zxvf go1.6.2.linux-amd64.tar.gz
mv go $HOME/

#setup env
echo -e "export GOROOT=$HOME/go\n export PATH=$PATH:$HOME/go/bin\n export GOPATH=$HOME/gopath" >> $HOME/.bashrc
. $HOME/.bashrc



# build package
cd $HOME/gopath/src/github.com/fakewechat/bin
./build.sh



# redis
cd $HOME/gopath/src/github.com/fakewechat/package
tar zxvf redis-2.8.24.tar.gz
cd redis-2.8.24
make
sudo make install

cp src/redis-server $HOME/gopath/src/github.com/fakewechat/bin
cp src/redis-cli $HOME/gopath/src/github.com/fakewechat/bin



sudo rpm -Uvh http://dl.fedoraproject.org/pub/epel/7/x86_64/e/epel-release-7-6.noarch.rpm
sudo rpm --import /etc/pki/rpm-gpg/RPM-GPG-KEY-EPEL-7

sudo yum install -y ansible 





