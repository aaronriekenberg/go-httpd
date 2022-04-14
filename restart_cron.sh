#!/bin/sh

pgrep go-httpd > /dev/null 2>&1
if [ $? -eq 1 ]; then
  cd ~/go-httpd
  ./restart.sh > /dev/null 2>&1
fi
