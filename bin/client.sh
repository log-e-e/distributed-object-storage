# 提供客户端相对便捷的操作命令
# PUT: ./client.sh put objectName "object content..."
# GET: ./client.sh get objectName
# DELETE: ./client.sh delete objectName

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
  objectContent=$3  # 对象数据
  hash=$(echo -n "${objectContent}" | openssl dgst -sha256 -binary | base64)  # 哈希值
  if [[ "${objectName}" = "" ]] || [[ "${objectContent}" = "" ]]
  then
    echo "Syntax Error: missing parameter objectName or objectContent"
    echo 'Example: put objectName1 "objectContent...."'
  else
    curl -v ${apiServer}/objects/${objectName} -XPUT -d "${objectContent}" -H "Digest: SHA-256=${hash}"
  fi
fi

# GET对象
if [ "${operationType}" = ${GET} ]
then
  objectName=$2  # 对象名
  if [ "${objectName}" = "" ]
  then
    echo "Syntax Error: missing objectName"
    echo "Example: get objectName1"
  else
    curl -v ${apiServer}/objects/${objectName}
  fi
fi

# DELETE对象
if [ "${operationType}" = ${DELETE} ]
then
  objectName=$2  # 对象名
  if [ "${objectName}" = "" ]
  then
    echo "Syntax Error: missing objectName"
    echo "Example: delete objectName1"
  else
    curl -v ${apiServer}/objects/${objectName} -XDELETE
  fi
fi

# 获取帮助示例
echo ""
echo "---------------------------------------operation manual---------------------------------------"
echo 'PUT: ./client.sh put objectName "objectContent..."'
echo "GET: ./client.sh get objectName"
echo "DELETE: ./client.sh delete objectName"
echo "---------------------------------------operation manual---------------------------------------"
echo ""
