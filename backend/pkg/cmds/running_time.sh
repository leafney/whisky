#!/bin/sh
# 获取系统运行时间

current_timestamp=$(date +%s)
boot_timestamp=$(cat /proc/stat | awk '/btime/ {print $2}')
uptime_seconds=$((current_timestamp - boot_timestamp))

days=$(($uptime_seconds / 86400))
hours=$(($uptime_seconds % 86400 / 3600))
minutes=$(($uptime_seconds % 3600 / 60))
seconds=$(($uptime_seconds % 60))

running_time="${days}d ${hours}h ${minutes}m ${seconds}s"
echo -n "$running_time"