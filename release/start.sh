#!/bin/sh

work_path=`pwd`

a=`ps -ef|grep chakala_go|grep -v grep|wc -l`

if(($a==0))
then
nohup "$work_path""/chakala_go" &
else 
echo "Have Started!"
fi

