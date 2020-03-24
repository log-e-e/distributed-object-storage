# 10s后执行，目的是等待服务启动
sleep 10
# 获取第二个参数，如果没有第二个参数
systemFolderPath=$1
if [ "$1" = "&" ]; then
  systemFolderPath=storage-system
fi
# 判断路径是否有效，无效则终止脚本的执行
if [ ! -d "$systemFolderPath" ]; then
  exit 1
fi

# 端口号
port=12345

# 监控进程是否停止
while true
do
  # 根据port获取对应的进程的pid
  pid=$(netstat -nlp | grep :$port | awk '{print $7}' | awk -F"/" '{ print $1 }');
  if [ "$pid" = "" ]; then
    rm -rf "$systemFolderPath"
    # 删除nohub产生的文件
    rm ./nohup.out
    break
  fi
  sleep 60
done
