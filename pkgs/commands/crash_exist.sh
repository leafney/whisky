#!/bin/sh
# check crash service exist

file_path="/etc/init.d/shellcrash"

if [ -f "$file_path" ]; then
  echo "true"
else
  echo "false"
fi