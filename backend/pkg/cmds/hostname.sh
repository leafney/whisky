#!/bin/sh

hostname=$(cat /proc/sys/kernel/hostname)
echo -n "$hostname"