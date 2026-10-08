#!/bin/sh

# 遍历所有网络设备
for device in $(ls /sys/class/net); do
  # 检查设备是否为虚拟设备
  if [[ "$device" == "lo" ]]; then
    continue
  fi

  # 获取设备的 IP 地址
  ip_address=$(ip addr show dev "$device" | grep 'inet ' | awk '{print $2}' | cut -d/ -f1)

  # 输出结果
  if [[ -n "$ip_address" ]]; then
    # 这里 echo 不加 -n，使其输出内容换行显示
    echo "${device}#${ip_address}"
  fi
done