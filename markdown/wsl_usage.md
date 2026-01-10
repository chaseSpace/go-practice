## WSL 使用

### 使用镜像网络

windows 搜索 `wsl settings`，修改网络——网络模式——Mirrored，重启 WSL 。

### WSL 指令

```bash
# 帮助信息
wsl --help
wsl -h

# 更新WSL
wsl --update

# 查看可用发行版
wsl --list --online
wsl -l -o

# 安装发行版
wsl --install -d Debian

# 卸载发行版
wsl --unregister <Distro>

# 设置默认WSL版本
wsl --set-default-version 2

# 设置特定发行版为WSL2
wsl --set-version ubuntu 2

# 设置默认发行版
wsl --set-default <Distro>
wsl -s <Distro>

# 查看状态
wsl --status
wsl --list --verbose
wsl -l -v

# 运行命令
wsl ls
wsl -d <Distro> <command>

# 关闭子系统
wsl --terminate <Distro>
wsl -t <Distro>
wsl --shutdown
```