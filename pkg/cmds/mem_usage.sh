#!/bin/sh
# 获取内存使用率

mem_usage=$(awk '/MemTotal/ {total=$2}
                   /MemFree/  {free=$2}
                   /Buffers/  {buffers=$2}
                   /Cached/   {cached=$2}
                   END {
                     used=total-free-buffers-cached
                     percent=used/total*100
                     printf "%.2f%%", percent
                   }' /proc/meminfo)
echo -n "$mem_usage"