#!/bin/sh

KILL_CMD=pkill

$KILL_CMD go-httpd

sleep 2

export PATH=${HOME}/bin:$PATH

nohup ./go-httpd ./configfiles/apu2-config.json > output 2>&1 &
