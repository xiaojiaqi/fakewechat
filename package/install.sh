yum install -y gcc gcc-c++ autoconf automake libtool vim  wget


curl -O -L https://storage.googleapis.com/golang/go1.6.2.linux-amd64.tar.gz -kvv

wget https://bootstrap.pypa.io/ez_setup.py -O - | sudo python


tar zxvf go1.6.2.linux-amd64.tar.gz

mv go ~/

echo -e "export GOROOT=~/go\n export PATH=$PATH:~/go/bin\n export GOPATH=~/gopath" >> ~/.bashrc

. ~/.bashrc

tar zxvf protobuf-master.tar.gz

cd protobuf-master

./autogen.sh

./configure 
make



cd python

sudo python setup.py install

cd ~/gopath/src/github.com/fakewechat/bin


sudo easy_install redis







