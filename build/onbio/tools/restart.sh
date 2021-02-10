#!/bin/sh
cd /data/code/simple_web;
master_pid=`ps -eaf | grep onbio_web | grep -v grep | awk '{print $2}'`

if [ $master_pid ]; then
    kill -9 $master_pid
fi

./bin/onbio_web ./conf/conf.json & 
sleep 1
master_pid_new=`ps -eaf | grep onbio | grep -v grep | awk '{print $2}'`
echo $master_pid_new > ./bin/onbio.pid
