```shell
# 找到私钥
$ ls ~/.ssh/
config  id_ed25519  id_ed25519.pub  known_hosts  known_hosts.old

# 写入ssh config
# 443代替22端口（运营商屏蔽了某些节点的22端口）
echo 'Host github.com
        User git
        Hostname ssh.github.com
        PreferredAuthentications publickey
        IdentityFile ~/.ssh/id_ed25519
        Port 443
        ProxyCommand connect -S 127.0.0.1:7890 %h %p' > ~/.ssh/config
```