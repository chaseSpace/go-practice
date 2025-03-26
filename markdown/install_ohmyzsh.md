## 快速安装ohmyzsh

[什么是ohmyzsh？](https://www.google.com.hk/search?q=什么是ohmyzsh)

```shell
yum install -y git zsh wget

# hk主机
wget -qO- https://gitee.com/mirrors/oh-my-zsh/raw/master/tools/install.sh | bash

wget https://gitee.com/mirrors/oh-my-zsh/raw/master/tools/install.sh -O install_onymzsh.sh
vi install_onymzsh.sh
# 修改下面两行
# REPO=${REPO:-ohmyzsh/ohmyzsh}
# REMOTE=${REMOTE:-https://github.com/${REPO}.git}
# 为
# REPO=${REPO:-mirrors/oh-my-zsh}
# REMOTE=${REMOTE:-https://gitee.com/${REPO}.git}
# 保存 并 执行
chmod +x install_onymzsh.sh && ./install_onymzsh.sh

# 修改主题
ls ~/.oh-my-zsh/themes
vi ~/.zshrc
# 找到 ZSH_THEME 行，修改为自己想用的主题名称即可

# 安装插件
git clone https://gitee.com/jsharkc/zsh-autosuggestions.git $ZSH_CUSTOM/plugins/zsh-autosuggestions
git clone https://gitee.com/jsharkc/zsh-syntax-highlighting.git $ZSH_CUSTOM/plugins/zsh-syntax-highlighting

# 配置插件
sed -i 's/plugins=(git)/plugins=(git zsh-autosuggestions zsh-syntax-highlighting)/' ~/.zshrc
# 设置别名
echo 'alias kk="kubectl"' >> ~/.zshrc
#echo 'alias m="minikube"' >> ~/.zshrc # 如果安装了minikube
echo 'DISABLE_AUTO_UPDATE=true' >> ~/.zshrc

# 生效
source ~/.zshrc
```