service_name="monitor_file_lastchange_time"
bin_path="/usr/local/bin/monitor_file_lastchange_time"
description="监控指定目录下最后一个文件更新时间是否超过指定时间并发送微信告警"

## 注册为服务
## 准备配置文件
cat <<EOF > /etc/systemd/system/${service_name}.service
[Unit]
Description=${description}

[Service]
Restart=always
RestartSec=5
ExecStart=${bin_path}

[Install]
WantedBy=multi-user.target
EOF

## 启动并设置开机自动启动
systemctl daemon-reload
systemctl enable ${service_name}.service
systemctl stop ${service_name}.service
systemctl start ${service_name}.service
## systemctl status ${service_name}.service
journalctl -f -u ${service_name}.service