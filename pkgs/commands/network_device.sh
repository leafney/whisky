#!/bin/sh
# 获取同网络下设备列表

devices=$(cat /tmp/dhcp.leases)
echo "$devices"