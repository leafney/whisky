#!/bin/sh
# 获取磁盘使用率

disk_usage=$(df / | tail -n 1 | awk '{print $5}')
echo -n "$disk_usage"