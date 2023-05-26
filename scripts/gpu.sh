#!/bin/bash

echo "File name is $1"

if [ -n "$1" ]
then
    echo "The \$1 is $1"
else
    echo "\$1 未提供."
fi

sshpass -p 'h4L&$IQW' ssh stu1642@pilogin.hpc.sjtu.edu.cn  "cd data && sbatch $1"