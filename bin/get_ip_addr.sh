# 测试ifconfig是否存在
type ifconfig

# 其中的"$?"即为上一条执行的命令"type ifconfig"，命令执行成功会返回0，否则返回其他数字
# if语句用于判断上条命令是否执行成功，若执行成功则执行if中的命令，否则执行else中的命令
if [ $? -eq 0 ]; then
  IPAddresses=$(ifconfig -a|grep inet|grep -v 127.0.0.1|grep -v inet6|awk '{print $2}'|tr -d "addr:")
  echo "--------------------------------------------------------------------"
  echo "Local IP Address(es):"
  echo "$IPAddresses"
  echo "Plesas use your ip address to use the storage system."
  echo "--------------------------------------------------------------------"
else
  echo "FAILED: command 'ifconfig' not found, we can not find your IP Address.
  Please try to run command 'sudo apt install net-tools' to install ifconfig."
fi
