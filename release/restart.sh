#!/bin/sh

work_path=`pwd`

while((1))
do

a=`ps -ef|grep chakala_go|grep -v grep|wc -l`

if(($a==0))
then
nohup "$work_path""/chakala_go" &
else 
sleep 1
fi

done

