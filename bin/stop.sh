# 终止各个服务节点的进程
# IP:PORT信息存储
ipAddrFilePath=.ipAddrs
while read LINE
do
  pid=$(netstat -nlp | grep $LINE | awk '{print $7}' | awk -F"/" '{ print $1 }');
  if [ "${pid}" != "" ]; then
    netstat -pan | grep ${LINE}
    kill -9 $pid
    echo "${pid} killed"
  fi
done < ${ipAddrFilePath}

# 删除存储数据的文件
storageRoot=/home/${USER}/storage-system
if [ -e ${storageRoot} ]; then
  rm -rf ${storageRoot}
fi

# 删除所有的虚拟IP
netCard=wlp2s0  # 网卡名
serverNodeAmount=8  # IP总数
for i in `seq 1 ${serverNodeAmount}` ; do
  sudo ifconfig ${netCard}:${i} down
done
