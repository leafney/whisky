#!/bin/sh
# 获取 cpu 温度

temp_cpu=$(cat /sys/class/thermal/thermal_zone0/temp | awk '{printf "%.2f°C", $1/1000}')
echo -n "$temp_cpu"