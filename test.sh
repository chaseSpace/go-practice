echo "now=$(pwd) && cd /usr/local && \
wget https://studygolang.com/dl/golang/$aim && \
tar xzf $aim && \
rm -rf $aim && \
mkdir -p /usr/local/go/gopath && \
echo export PATH=$PATH:/usr/local/go/bin:/usr/local/gopath/bin >>/root/.bashrc && \
echo export GOPATH=/usr/local/gopath >> /root/.bashrc && \
echo export GOPROXY=https://goproxy.cn >> /root/.bashrc && \
echo export GO111module=on >> /root/.bashrc  && \
echo export GOROOT=/usr/local/go >> /root/.bashrc  && \
source  /root/.bashrc && \
cd $now && \
go version

" > install_go.sh