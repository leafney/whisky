#!/bin/sh
# 调用 yacd 接口获取 clash 运行状态
# 如果更改了 yacd 的端口号，需要手动指定，否则默认为 9999

default_port=9999

port="${1:-$default_port}"

yacd_stats=$(curl -s http://127.0.0.1:${port}/configs)
echo "$yacd_stats"