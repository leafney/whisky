#!/bin/sh
# check crash service exist

file_path="/etc/init.d/shellcrash"

if [ -f "$file_path" ]; then
  echo -n "true"
else
  echo -n "false"
fi