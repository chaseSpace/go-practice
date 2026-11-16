cat << 'EOF' > restore_sshd_log.sh
#!/bin/bash

if [ "$EUID" -ne 0 ]; then
  echo "❌ 错误：请使用 root 权限运行此脚本 (例如: sudo bash $0)"
  exit 1
fi

echo "=== 开始恢复 SSH 日志、登录记录与系统审计 ==="

# 自动识别服务名
if systemctl list-unit-files | grep -q "^sshd.service"; then
  SSH_SERVICE="sshd"
else
  SSH_SERVICE="ssh"
fi

# 1. 解除 wtmp / btmp / lastlog 的只读锁定（必须最先执行，否则后续服务写入会失败）
for bin_log in /var/log/wtmp /var/log/btmp /var/log/lastlog; do
  if [ -f "$bin_log" ]; then
    chattr -i "$bin_log" 2>/dev/null
  fi
done
echo "[1/5] 已解除 wtmp / btmp / lastlog 文件锁定"

# 2. 还原 sshd_config 配置
LATEST_BAK=$(ls -t /etc/ssh/sshd_config.bak.* 2>/dev/null | head -n 1)
if [ -n "$LATEST_BAK" ]; then
  cp "$LATEST_BAK" /etc/ssh/sshd_config
  echo "[2/5] 已从最新备份还原 sshd_config: $LATEST_BAK"
else
  # 若无备份文件，则直接将 LogLevel 改回默认 INFO
  sed -i 's/^LogLevel QUIET/LogLevel INFO/' /etc/ssh/sshd_config
  echo "[2/5] 未找到备份文件，已手动将 LogLevel 恢复为 INFO"
fi

# 校验 sshd 语法
if ! sshd -t; then
  echo "❌ 警告：sshd_config 语法校验存在异常，请手动检查文件！"
  exit 1
fi

# 3. 移除 rsyslog 拦截规则并重启
if [ -f /etc/rsyslog.d/00-drop-sshd.conf ]; then
  rm -f /etc/rsyslog.d/00-drop-sshd.conf
fi
systemctl restart rsyslog 2>/dev/null
echo "[3/5] 已移除 rsyslog 屏蔽规则并重启服务"

# 4. 移除 systemd-journald 的 nolog 覆盖配置并重启 SSH
SYSTEMD_DIR="/etc/systemd/system/${SSH_SERVICE}.service.d"
if [ -f "${SYSTEMD_DIR}/nolog.conf" ]; then
  rm -f "${SYSTEMD_DIR}/nolog.conf"
  rmdir "$SYSTEMD_DIR" 2>/dev/null
fi
systemctl daemon-reload
systemctl restart "$SSH_SERVICE"
echo "[4/5] 已恢复 systemd-journald 记录并重启 $SSH_SERVICE"

# 5. 重新启用并启动 auditd 内核审计服务
if command -v auditctl >/dev/null 2>&1 || systemctl list-unit-files | grep -q "^auditd.service"; then
  systemctl enable auditd 2>/dev/null
  service auditd start 2>/dev/null || systemctl start auditd 2>/dev/null
  echo "[5/5] 已重新启用并启动 auditd 审计服务"
else
  echo "[5/5] 系统未安装 auditd，跳过恢复"
fi

# 卸载绑定挂载
sudo umount /run/utmp

# 删除占位文件
sudo chattr -i /etc/empty_utmp 2>/dev/null
sudo rm -f /etc/empty_utmp


echo "=== ✅ 全部恢复完成！所有登录与审计日志已恢复正常记录 ==="
EOF
