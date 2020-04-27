# distributed-object-storage
源自《分布式对象存储——原理、架构及Go语言实现》

**下载链接**：[《分布式对象存储——原理、架构及Go语言实现》](https://www.lanzous.com/i7dvkzg)

## 前言

在执行前，需要根据自身网络环境，修改`/bin/startup.sh`中的网络配置参数，主要是IP、网关以及网卡名，如下：

```shell script
# -----------START: 网络配置，需根据自己的机器自行修改-----------
port=12345  # 服务端口号
netCard=wlp2s0  # 网卡名
netID=192.168.0  # 网络标识，指的是某一个物理网络的网络标识，所有连接在该网络中的主机共用相同的网络标识
hostID=101  # 主机标识，指的是为连接在某一个物理网络上的主机分配的用于区分其他主机的标识
gateway=24  # 网关
# -----------END: 网络配置，需根据自己的机器自行修改-----------
```

## 第一章

已实现书中的代码，能够正确运行，步骤如下：

### 1. 启动服务与停止服务

#### 1.1. 启动服务

执行`bin`目录下的`startup.sh`脚本

```shell script
cd bin && chmod +x startup.sh && ./startup.sh
```

执行启动脚本后，会在当前用户的目录下生成文件夹。路径为：`~/storage-system/storage/objects`

#### 1.2. 停止服务

执行`bin`目录下的`startup.sh`脚本，在服务进行中的任何时候都可以执行该脚本

```shell script
cd bin && chmod +x stop.sh && ./stop.sh
```

执行`stop.sh`后，会停止服务并清除服务生成的文件夹及文件，即删除`~/storage-system`

### 2. 执行PUT及GET操作

服务启动后，如果能够通过`ifconfig`或`ip addr`命令获取本机ip地址，则会给出相应的ip地址。我们可以直接使用给出的ip地址，执行PUT和GET操作：

#### 2.1 PUT操作

```shell script
# 命令格式
curl -v IP地址:12345/objects/对象文件名 -XPUT -d "This is the object content"

# 示例
curl -v 192.168.43.120:12345/objects/test -XPUT -d "This is the object content"
```

可能的执行结果如下
- 执行成功
  - 执行命令后即可从对象存储系统的路径`~/storage-system/storage/objects`下看到`test`文件
  - 在Server窗口会反馈类似`2020/03/24 17:35:18 PUT SUCCESS: object '/home/jihonghe/storage-system/storage/objects/test3'`的消息
  - 在Client窗口的反馈消息中会存在`HTTP/1.1 200 OK`
- 执行失败
  - `HTTP/1.1 405 Method Not Allowed`，指的是请求方法不支持，需要检查执行请求的命令行是否正确


#### 2.2 GET操作

```shell script
# 命令格式
curl -v IP地址:12345/objects/对象文件名

# 示例
curl -v 192.168.43.120:12345/objects/test
```

可能的执行结果如下

- 执行成功
  - 在Server窗口会反馈类似`2020/03/24 17:36:32 GET SUCCESS: object '/home/jihonghe/storage-system/storage/objects/test3'`的消息
  - 在Client窗口的反馈消息中会存在`HTTP/1.1 200 OK`及在最后一行会看到`This is content`数据
- 执行失败
  - 在Server窗口会反馈类似`2020/03/24 17:38:26 GET FAILED: open /home/jihonghe/storage-system/storage/objects/test4: no such file or directory`的消息
  - 在Client窗口则会反馈类似`HTTP/1.1 500 Internal Server Error`的数据

## 第二章

已实现第二章的全部功能，能正确运行，具体步骤如下：

### 1. 启动服务与停止服务

#### 1.1. 启动服务

执行`bin`目录下的`startup.sh`脚本

```shell script
cd bin && chmod +x startup.sh && ./startup.sh
```

启动脚本主要完成如下几件事情：

- 启动`RabbitMQ`，具体的配置可从书中查看
- 创建各个服务节点的虚拟ip，可用`ifconfig`或`ip addr`查看
- 创建服务相关的文件夹，其公共路径为：`~/storage-system/storage/`，`storage`目录树如下
    ```
    └── storage
        ├── dataNode1
        │   └── objects
        ├── dataNode2
        │   └── objects
        ├── dataNode3
        │   └── objects
        ├── dataNode4
        │   └── objects
        ├── dataNode5
        │   └── objects
        └── dataNode6
            └── objects
    ```
- 启动各个节点的服务，在终端后输出如下内容
    ```
    INFO: finish startup RabbitMQ & set RABBITMQ_SERVER environment
    INFO: new dataServer started. storageRoot=/home/jihonghe/storage-system/storage/dataNode1, listenAddr=192.168.0.102:12345
    INFO: new dataServer started. storageRoot=/home/jihonghe/storage-system/storage/dataNode2, listenAddr=192.168.0.103:12345
    INFO: new dataServer started. storageRoot=/home/jihonghe/storage-system/storage/dataNode3, listenAddr=192.168.0.104:12345
    INFO: new dataServer started. storageRoot=/home/jihonghe/storage-system/storage/dataNode4, listenAddr=192.168.0.105:12345
    INFO: new dataServer started. storageRoot=/home/jihonghe/storage-system/storage/dataNode5, listenAddr=192.168.0.106:12345
    INFO: new dataServer started. storageRoot=/home/jihonghe/storage-system/storage/dataNode6, listenAddr=192.168.0.107:12345
    INFO: new apiServer started. listenAddr=192.168.0.108:12345
    INFO: new apiServer started. listenAddr=192.168.0.109:12345
    ```

#### 1.2. 停止服务

执行`bin`目录下的`startup.sh`脚本，在服务进行中的任何时候都可以执行该脚本，该脚本会删除`~/storage-system`及创建的虚拟IP

```shell script
cd bin && chmod +x stop.sh && ./stop.sh
```

### 2. 执行PUT及GET操作

本章的实现了`PUT`与`GET`操作，即对象的存储与查询功能，其操作命令与第一章相同，但是数据存储是随机存储，不再赘述。
**注意**：在使用`curl`命令时，我们需要使用的是`apiServer`的请求地址，即在启动服务时生成的虚拟`IP:PORT`：`192.168.0.108:12345`或`192.168.0.109:12345`，
若是没有使用`apiServer`的虚拟IP，则会是定向存储而非随机存储

## 第三章

已实现第三章的所有功能，但需要说明的是：笔者所使用的ES版本是截止目前的最新版（7.6.2），书中使用的ES版本是6.x。另外，ES接口的实现也没有按照书中的方式来，笔者使用
的是ES的第三方包 [elastic](https://github.com/olivere/elastic) 实现的相关操作。并且，在测试期间遇到了关于命令`curl`的一点问题，会在后文提出。后续的内容，
主要阐述基于第二章的基础上新增的功能与注意事项：

**新增功能**
- 获取指定对象的所有版本的元数据
- 获取指定对象的具体版本的元数据
- 删除对象（逻辑删除，仅通过修改元数据信息实现逻辑删除）

### 1. 环境准备

#### 1.1 在`startup.sh`中添加环境变量

基于第二章的基础上，新增了Elasticsearch，因此我们需要安装ES。推荐 [华为的镜像站下载](https://mirrors.huaweicloud.com/elasticsearch/) ，速度极快，有所有的版本。
安装完ES后，我们需要创建索引映射（不熟悉ES的胖友请自行学习）。ES7.x版本的创建方式与6.x的方式有些许不同，7.x去不建议使用type。因此，若是使用7.x版本的ES则不能按照书中的来：
**使用kibana创建**
```
PUT /metadata
{
    "mappings": {
        "properties": {
            "name": {"type": "keyword"},
            "version": {"type": "integer"},
            "size": {"type": "integer"},
            "hash": {"type": "keyword"}
        }
    }
}
```

创建完索引与映射后，需要在`bin/startup.sh`中添加相应的环境变量，除此之外`startup.sh`无需做任何改变：
```shell script
export ES_SERVER=localhost:9200
```

#### 1.2 hash值shell脚本

该shell脚本主要是为了方便生成哈希值，并且生成哈希值所采用的加密方式采用的是`SHA256`，这是因为在使用`SHA256`测试时发现其生成的哈希值存在`/`，这会导致在解析`url`获取
对象名时出现问题，因此改用了`SHA1`。脚本调用方式如下：
```shell script
./hash-gen.sh "content"
# 示例
./hash-gen.sh "This is object content version-2"
```

### 2. 对象的PUT，GET，DELETE操作及定位服务

#### 2.1 PUT操作

PUT对象的命令相对于第二章的区别仅仅是在Digest头部添加了对象内容的哈希值，命令如下：
```shell script
curl -v 192.168.0.108:12345/objects/test -XPUT -d "This is object content version-1" -H "Digest: SHA-256=9AimTha2kCISf8bVfi1jPXo2BzY="
```

#### 2.2 Locate服务
Locate服务相较于上一章的变化主要是使用对象的哈希值作为对象名，其他的无变化：
```shell script
curl -v 192.168.0.108:12345/locate/SoxiAi+lEo63eTZ6rc62tzw8kSA=
```

#### 2.3 GET操作

**注意**：按照书中的方式`curl -v 192.168.0.108:12345/objects/test?version`请求数据，在Linux下是无效操作，需要对问号使用`\`进行转义

```shell script
# 查看对象名为test的所有版本的元数据信息
curl -v 192.168.0.108:12345/objects/test
# 查看对象名为test的指定版本的元数据信息，注意携带的参数的问号要转义
curl -v 192.168.0.108:12345/objects/test\?version
```

#### 2.4 DELETE操作

对象的删除是逻辑删除，规则：基于最新版本新增一个版本号，将size和hash值置空，表示该对象已被删除。命令：
```shell script
curl -v 192.168.0.108:12345/objects/test -XDELETE
```

删除后，我们使用GET操作相关的命令会看到新增一个版本，且其hash和size为空值。除此外，我们仍然可以获取旧版本的信息。

## 第四章

已实现所有功能，能够正常地校验、去重及清理临时文件。并且，无相关功能调用命令的变化。下面会给出一些shell脚本的变化：

### 1. `bin/startup.sh`

在该脚本中，新增了关于ES的启动与索引和映射的创建命令。这样一来，关于环境准备的所有事情都无需进行额外操作。
```shell script
# 启动ES
sudo systemctl start elasticsearch.service
# 创建metadata索引和映射
curl -XPUT localhost:9200/metadata -H 'Content-Type:application/json' -d '{
  "mappings": {
    "properties": {
      "name": {"type": "keyword"},
      "version": {"type": "integer"},
      "size": {"type": "integer"},
      "hash": {"type": "keyword"}
    }
  }
}'
echo ""
echo "INFO: finish start elasticsearch && create index 'metadata' and mappings"
```

### 2. `bin/stop.sh`

在该脚本中，新增了删除ES索引的命令
```shell script
# 删除ES索引
curl -X DELETE 'localhost:9200/metadata'
```

### 3. `bin/hash-gen.sh`

在第三章中阐述的原因是错误的，并不是因为存在斜线的原因造成的，故在此纠正，现改为SHA256加密，BASE64编码。
```shell script
# 生成哈希值
content=$1
# shellcheck disable=SC2006
echo "Start generate `${content}` hash code..."
echo -n "$content" | openssl dgst -sha256 -binary | base64
```

## 第五章

已实现本章所有功能。另外，通过新增shell脚本改进了测试的便利性，封装了RESTful api的命令请求，实现了通过执行脚本及提供相关的参数便可进行PUT、GET和DELETE操作：

`bin/client.sh`
```shell script
# PUT操作
./client.sh put objectName "content"
# GET操作
./client.sh get objectName
# DELETE操作
./client.sh delete objectName
```
