#!/bin/sh
# https://www.ipip.net/myip.html

ip=$(curl -s http://myip.ipip.net)
echo -n "$ip"