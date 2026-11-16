cat << 'EOF' > disable_sshd_log.sh
#!/bin/bash

# 确保以 root 权限运行
if [ "$EUID" -ne 0 ]; then
  echo "❌ 错误：请使用 root 权限运行此脚本 (例如: sudo bash $0)"
  exit 1
fi

echo "=== 开始配置彻底禁用 SSHD 日志 ==="

# 自动识别系统 SSH 服务名与日志路径
if systemctl list-unit-files | grep -q "^sshd.service"; then
  SSH_SERVICE="sshd"
else
  SSH_SERVICE="ssh"
fi

if [ -f /var/log/secure ]; then
  LOG_FILE="/var/log/secure"
elif [ -f /var/log/auth.log ]; then
  LOG_FILE="/var/log/auth.log"
else
  LOG_FILE=""
fi

echo "[*] 检测到 SSH 服务名: $SSH_SERVICE"
[ -n "$LOG_FILE" ] && echo "[*] 检测到认证日志文件: $LOG_FILE"

# 1. 备份并修改 sshd_config
BACKUP_FILE="/etc/ssh/sshd_config.bak.$(date +%F_%H%M%S)"
cp /etc/ssh/sshd_config "$BACKUP_FILE"
echo "[1/4] 已备份配置至: $BACKUP_FILE"

if grep -qE "^#*LogLevel" /etc/ssh/sshd_config; then
  sed -i 's/^#*LogLevel.*/LogLevel QUIET/' /etc/ssh/sshd_config
else
  echo "LogLevel QUIET" >> /etc/ssh/sshd_config
fi

# 校验 sshd 配置语法，防止断连
if ! sshd -t; then
  echo "❌ 配置文件校验失败，正在自动还原备份..."
  cp "$BACKUP_FILE" /etc/ssh/sshd_config
  exit 1
fi

# 2. 配置 rsyslog 丢弃 sshd 日志
if [ -d /etc/rsyslog.d ]; then
  cat << 'RSYSLOG_EOF' > /etc/rsyslog.d/00-drop-sshd.conf
:programname, isequal, "sshd" stop
if ($syslogfacility-text == 'auth' or $syslogfacility-text == 'authpriv') and ($msg contains 'sshd') then stop
RSYSLOG_EOF
  systemctl restart rsyslog 2>/dev/null
  echo "[2/4] 已配置 rsyslog 拦截规则并重启服务"
else
  echo "[2/4] 未检测到 /etc/rsyslog.d，跳过 rsyslog 配置"
fi

# 3. 禁用 systemd-journald 对 sshd 的日志收集
SYSTEMD_DIR="/etc/systemd/system/${SSH_SERVICE}.service.d"
mkdir -p "$SYSTEMD_DIR"
cat << 'SYSTEMD_EOF' > "${SYSTEMD_DIR}/nolog.conf"
[Service]
LogLevelMax=0
StandardOutput=null
StandardError=null
SYSTEMD_EOF

systemctl daemon-reload
systemctl restart "$SSH_SERVICE"
echo "[3/4] 已关闭 systemd-journald 收集并重启 $SSH_SERVICE"

# 4. 安全清空历史日志文件（保留原有权限与 SELinux 上下文）
if [ -n "$LOG_FILE" ]; then
  truncate -s 0 "$LOG_FILE"
  echo "[4/4] 已安全清空历史日志: $LOG_FILE"
else
  echo "[4/4] 未找到传统认证日志文件，跳过清理"
fi

echo "=== ✅ 全部完成！请新开一个终端测试 SSH 连接是否正常 ==="
EOF
