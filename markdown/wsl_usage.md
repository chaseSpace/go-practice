## WSL 使用

### 介绍

使用WSL在Windows进行开发工作可以提供极大的便利，相比VMware或其他虚拟机方式，WSL可以使用PowerShell携带的指令进行快速下载启动，还可指定喜好的发行版。

最大的优势是在Windows上为WSL设置网络的镜像模式，以便WSL可以直接使用Windows的Clash代理（设置http_proxy），注意只有普通用户可以直接clash代理，root用户无法使用，全程不需要设置https_proxy。

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

### 推荐的发行版

AlmaLinux，作为CentOS的替代品。

[wsl_almalinux.md](wsl_almalinux.md) 是笔者的使用笔记。