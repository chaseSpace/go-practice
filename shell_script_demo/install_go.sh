# 执行此命令声明要下载的go版本
aim=go1.19.7.linux-amd64.tar.gz

#升级安装请先删除除了/usr/local/go/gopath以外的文件
#rm  /usr/local/go -r

# 执行下面的命令创建 install_go.sh
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

# 执行 . install_go.sh