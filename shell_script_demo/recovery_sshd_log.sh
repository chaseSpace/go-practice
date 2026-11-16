sudo rm -f /etc/rsyslog.d/00-drop-sshd.conf /etc/systemd/system/ssh*.service.d/nolog.conf
sudo cp /etc/ssh/sshd_config.bak.* /etc/ssh/sshd_config
sudo systemctl daemon-reload && sudo systemctl restart rsyslog sshd
