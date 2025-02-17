#!/bin/sh
# 调用 yacd 接口切换 Allow LAN 功能
# 需要传入两个参数：1为启用状态 true 或者 false; 2为使用的端口号，默认为 9999
# curl -s -o /dev/null -w "%{http_code}" -X PATCH -H "Content-Type:application/json" -d '{"allow-lan":false}'  http://127.0.0.1:9999/configs

# 设置默认端口号
default_port=9999

# 获取布尔类型参数，并验证格式
is_allow="$1"
if [ "$is_allow" != "true" ] && [ "$is_allow" != "false" ]; then
  echo "首个参数格式错误，应该为 true 或 false"
  exit 1
fi

# 获取端口号参数，如果为空则使用默认值
port="${2:-$default_port}"

# curl 请求
if [ "$is_allow" = "true" ]; then
  lan_result=$(curl -s -o /dev/null -w "%{http_code}" -X PATCH -H "Content-Type:application/json" -d '{"allow-lan":true}'  http://127.0.0.1:${port}/configs)
else
  lan_result=$(curl -s -o /dev/null -w "%{http_code}" -X PATCH -H "Content-Type:application/json" -d '{"allow-lan":false}'  http://127.0.0.1:${port}/configs)
fi

# 将结果输出到标准输出
echo -n "$lan_result"