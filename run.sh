#!/bin/bash

# 如果 pid.txt 存在，读取其中的 PID 并终止对应的进程
if [ -f pid.txt ]; then
    old_pid=$(cat pid.txt)
    if ps -p $old_pid > /dev/null 2>&1; then
        echo "Killing previous process with PID $old_pid"
        kill $old_pid
    fi
    rm pid.txt
fi

chmod +x ./openmcp-discord

# 使用 nohup 运行 main.py 并将输出重定向到 log 文件
nohup ./openmcp-discord  > /dev/null 2>&1 &

# 将新的 PID 保存到 pid.txt
echo $! > pid.txt
echo "已启动进程，PID为：$(cat pid.txt)"