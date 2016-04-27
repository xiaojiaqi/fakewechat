yum install -y gcc gcc-c++ autoconf automake libtool vim  wget

cd ~/gopath/src/github.com/fakewechat/package

#golang install
curl -O -L https://storage.googleapis.com/golang/go1.6.2.linux-amd64.tar.gz -kvv
tar zxvf go1.6.2.linux-amd64.tar.gz
mv go ~/

#setup env
echo -e "export GOROOT=~/go\n export PATH=$PATH:~/go/bin\n export GOPATH=~/gopath" >> ~/.bashrc
. ~/.bashrc

#install protbuf
cd ~/gopath/src/github.com/fakewechat/package
tar zxvf protobuf-master.tar.gz
cd protobuf-master
./autogen.sh
./configure 
make
cd python
sudo python setup.py install



# build package
cd ~/gopath/src/github.com/fakewechat/bin
./build.sh

#install redis for python
wget https://bootstrap.pypa.io/ez_setup.py -O - | sudo python
sudo easy_install redis

# redis
cd ~/gopath/src/github.com/fakewechat/package
tar zxvf redis-2.8.24.tar.gz
cd redis-2.8.24
make
sudo make install
cp src/redis-server ~/gopath/src/github.com/fakewechat/bin



rpm -Uvh http://dl.fedoraproject.org/pub/epel/7/x86_64/e/epel-release-7-6.noarch.rpm
rpm --import /etc/pki/rpm-gpg/RPM-GPG-KEY-EPEL-7

yum install -y ansible 





