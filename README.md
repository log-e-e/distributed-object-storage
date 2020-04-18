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
