#!/bin/sh

#default_port=9999

# 获取mode参数，并验证格式
mode_type="$1"

#port="${2:-$default_port}"
#
#mode_result="${mode_type}-${port}"

echo -n "$mode_type"