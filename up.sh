#!/bin/bash
# 	GitHb: https://github.com/chasonnchen/webot
# 	Author: Chasonn Chen <185320860@qq.com>
cmd=$1
if [ -n ${cmd} ]; then
    echo "Start to "${cmd}" the project."
else
     echo "Cmd error, eg: sh wkteam.sh restart"
     exit
fi



if [ "$cmd" == "build" ];then
    go build
    rm -rf runtime/webot
    mv webot runtime/
fi

if [ "$cmd" == "kill" ];then
    ps -ef | grep webot | grep -v grep | awk '{print $2}' | xargs kill -9
fi

if [ "$cmd" == "restart" ];then
    ps -ef | grep webot | grep -v grep | awk '{print $2}' | xargs kill -9
    cd runtime
    nohup ./webot 2>&1 &
fi
