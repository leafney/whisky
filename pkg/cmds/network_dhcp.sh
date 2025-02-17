#!/bin/sh
# 获取同网络下设备列表

dhcp=$(cat /tmp/dhcp.leases)
echo -n "$dhcp"