#!/bin/sh
# 调用 yacd 接口切换 运行模式 功能
# 三种模式：Global | Rule | Direct
# curl -s -o /dev/null -w "%{http_code}" -X PATCH -H "Content-Type:application/json" -d '{"mode":"Direct"}'  http://127.0.0.1:9999/configs
# 需要传入两个参数：1为模式值 rule、direct、global; 2为使用的端口号，默认为 9999

# 设置默认端口号
default_port=9999

# 获取mode参数，并验证格式
mode_type="$1"
if [ "$mode_type" != "rule" ] && [ "$mode_type" != "direct" ] && [ "$mode_type" != "global" ]; then
  echo "首个参数格式错误，应该为 rule 、direct 或 global"
  exit 1
fi

# 获取端口号参数，如果为空则使用默认值
port="${2:-$default_port}"

# 使用 case 语句判断参数值
case "$mode_type" in
  rule)
    mode_result=$(curl -s -o /dev/null -w "%{http_code}" -X PATCH -H "Content-Type:application/json" -d '{"mode":"Rule"}'  http://127.0.0.1:${port}/configs)
    ;;
  direct)
    mode_result=$(curl -s -o /dev/null -w "%{http_code}" -X PATCH -H "Content-Type:application/json" -d '{"mode":"Direct"}'  http://127.0.0.1:${port}/configs)
    ;;
  global)
    mode_result=$(curl -s -o /dev/null -w "%{http_code}" -X PATCH -H "Content-Type:application/json" -d '{"mode":"Global"}'  http://127.0.0.1:${port}/configs)
    ;;
  *)
    echo "参数值错误，必须为 rule 、direct 或 global"
    exit 1
    ;;
esac

# 将结果输出到标准输出
echo -n "$mode_result"
