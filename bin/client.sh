# 提供客户端相对便捷的操作命令
# PUT:
# 两种方式：一种是直接提供对象数据，另一种是提供对象数据文件路径
#     1. `./client.sh put objectName -content "object content"`
#     2. `./client.sh put objectName -path object-file-path`
# GET:
# 两种方式：一种是仅提供对象名，默认返回最新版本；另一种是给出特定版本
#     1. `./client.sh get objectName filepath` , 该方式有问题，暂时不用
#     2. `./client.sh get objectName version filepath`
# DELETE:
#     1. `./client.sh delete objectName`

apiServerNodeAmount=2  # 接口服务层节点数
dataServerNodeAmount=6  # 数据服务层节点数
# 随机选取apiServer节点
seed=$(date +%s)
randomIndex=$(($seed%${apiServerNodeAmount}+$(($dataServerNodeAmount+1))))
# 获取apiServer的IP:PORT，将其放到数组中
ipAddrFilePath=".ipAddrs"
index=1
apiServer=""
while read LINE
do
  if [ $index -eq $randomIndex ]
  then
    apiServer=$LINE
    break
  fi
  index=$(($index+1))
done < ${ipAddrFilePath}

# 操作类型
PUT="put"  # 存储对象
GET="get"  # 获取对象数据
DELETE="delete"  # 删除对象

# 命令参数
operationType=$1

# PUT对象
if [ "${operationType}" = ${PUT} ]
then
  objectName=$2  # 对象名
  ioType=$3  # 对象输入类型（对象值或对象文件）
  ioValue=$4  # 对象数据/对象文件
  if [ "${objectName}" = "" ] || [ "${ioType}" = "" ] || [ "${ioValue}" = "" ]
  then
    echo "Syntax Error: missing parameter, please checkout your command parameter"
  else
    go run ../client_server/client.go -apiServer="${apiServer}" -operationType="${operationType}" -objectName="${objectName}" -ioType="${ioType}" -ioValue="${ioValue}"
  fi
fi

# GET对象
if [ "${operationType}" = ${GET} ]
then
  objectName=$2  # 对象名
  version="$3"  # 对象版本号
  file=$4  # 保存数据的文件路径
  if [ "${objectName}" = "" ] && [ "${version}" = "" ] && [ "${file}" = "" ]
  then
    echo "Syntax Error: missing parameter"
  else
    if [ "${version}" = "" ]
    then
      go run ../client_server/client.go -apiServer="${apiServer}" -operationType="${operationType}" -objectName="${objectName}" -file="${file}"
    else
      go run ../client_server/client.go -apiServer="${apiServer}" -operationType="${operationType}" -objectName="${objectName}" -version="${version}" -file="${file}"
    fi
  fi
fi

# DELETE对象
if [ "${operationType}" = ${DELETE} ]
then
  objectName=$2  # 对象名
  if [ "${objectName}" = "" ]
  then
    echo "Syntax Error: missing parameter"
  else
    go run ../client_server/client.go -apiServer="${apiServer}" -operationType="${operationType}" -objectName="${objectName}"
  fi
fi

# 获取帮助示例
echo ""
echo "---------------------------------------operation manual---------------------------------------"
echo 'PUT:'
echo '    ./client.sh put objectName -content "objectContent..."'
echo '    ./client.sh put objectName -path object-file-path'
echo "GET: ./client.sh get objectName"
echo "DELETE: ./client.sh delete objectName"
echo "---------------------------------------operation manual---------------------------------------"
echo ""
