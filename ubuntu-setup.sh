# DOCKER INSTALLATION
#curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
#sudo add-apt-repository \
#   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
#   $(lsb_release -cs) \
#   stable"
#sudo apt-get update
#sudo apt-get install docker-ce
wget https://get.docker.com/builds/Linux/x86_64/docker-17.04.0-ce.tgz
tar xzvf docker-17.04.0-ce.tgz
echo 'export PATH="$HOME/docker:$PATH"' >> .bashrc
source .bashrc

# GOLANG INSTALLATION
sudo add-apt-repository ppa:gophers/archive
sudo apt update
sudo apt upgrade
sudo apt-get install golang-1.9-go
echo 'export PATH="/usr/lib/go-1.9/bin:$PATH"' >> .bashrc
mkdir ~/go
echo 'export GOPATH="$HOME/go"' >> .bashrc
source .bashrc
go version

# PROTOCOL BUFFERS INSTALLATION
sudo apt-get install unzip
wget -c https://github.com/google/protobuf/releases/download/v3.1.0/protoc-3.1.0-linux-x86_64.zip
sudo unzip protoc-3.1.0-linux-x86_64.zip -d /usr/local

# RUNC INSTALLATION
# Use revision from https://github.com/containerd/containerd/blob/master/RUNC.md
go get github.com/opencontainers/runc
cd $GOPATH/src/github.com/opencontainers/runc
git checkout 74a17296470088de3805e138d3d87c62e613dfc4
make
sudo make install

# CONTAINERD INSTALLATION
# Use latest release from https://github.com/containerd/containerd/releases
wget https://github.com/containerd/containerd/releases/download/v1.0.0-beta.3/containerd-1.0.0-beta.3.linux-amd64.tar.gz
tar xf containerd-1.0.0-beta.3.linux-amd64.tar.gz
sudo cp bin/* /usr/local/bin

# 46BIT/DISCOVERY INSTALLATION
go get github.com/46bit/discovery
cd $GOPATH/src/github.com/46bit/discovery
make
