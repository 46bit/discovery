set -euv

test "$USER" == "ubuntu"
test -d "$HOME/go/src/github.com/46bit/discovery"
sudo chown -R "$USER":"$USER" "$HOME/go"

sudo apt-get update -yqq
sudo apt-get upgrade -yqq

# echo "Installing Docker"
# curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
# sudo add-apt-repository -y \
#    "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
#    $(lsb_release -cs) \
#    stable"
# sudo apt-get update -yqq
# sudo apt-get install -yqq docker-ce

echo "Installing Golang"
sudo add-apt-repository -y ppa:gophers/archive
sudo apt-get update -yqq
sudo apt-get install -yqq golang-1.9-go
test -d $HOME/go
echo 'export GOPATH="$HOME/go"' >> ~/.bashrc
echo 'export PATH="$HOME/go/bin:/usr/lib/go-1.9/bin:$PATH"' >> ~/.bashrc
export GOPATH="$HOME/go"
export PATH="$HOME/go/bin:/usr/lib/go-1.9/bin:$PATH"
go version
go get github.com/onsi/ginkgo/ginkgo

echo "Installing Protobufs"
sudo apt-get install -yqq unzip
wget -nvc https://github.com/google/protobuf/releases/download/v3.5.0/protoc-3.5.0-linux-x86_64.zip
sudo unzip protoc-3.5.0-linux-x86_64.zip -d /usr/local

echo "Installing runc"
sudo apt-get install -yqq libseccomp-dev
go get github.com/opencontainers/runc
cd $GOPATH/src/github.com/opencontainers/runc
git checkout 7f24b40cc5423969b4554ef04ba0b00e2b4ba010
make
sudo make install

echo "Installing containerd"
sudo apt-get install -yqq btrfs-tools
wget -nv https://github.com/containerd/containerd/releases/download/v1.0.0/containerd-1.0.0.linux-amd64.tar.gz
tar xf containerd-1.0.0.linux-amd64.tar.gz
sudo cp bin/* /usr/local/bin

echo "Starting containerd"
cd $GOPATH/src/github.com/46bit/discovery/rainbow
sudo cp containerd.service /etc/systemd/system/containerd.service
sudo systemctl start containerd
sudo systemctl status containerd
sudo systemctl enable containerd

echo "Installing rainbow"
cd $GOPATH/src/github.com/46bit/discovery/rainbow
make cmd NAME=rainbowd
chmod +x bin/rainbowd
sudo cp bin/rainbowd /usr/local/bin/rainbowd

echo "Starting rainbow on localhost:4601"
sudo cp rainbowd.service /etc/systemd/system/rainbowd.service
sudo systemctl start rainbowd
sudo systemctl status rainbowd
sudo systemctl enable rainbowd
