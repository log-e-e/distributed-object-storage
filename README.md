# distributed-object-storage
源自《分布式对象存储——原理、架构及Go语言实现》

**下载链接**：[《分布式对象存储——原理、架构及Go语言实现》](https://www.lanzous.com/i7dvkzg)

## 第一章

已实现书中的代码，能够正确运行，步骤如下：

### 1. 启动服务

执行`bin`目录下的`chapter-01-STORAGE-SYSTEM-startup.sh`脚本

```shell script
cd bin && chmod +x startup.sh && ./startup.sh
```

执行启动脚本后，会在当前用户的目录下生成文件夹。路径为：`~/storage-system/storage/objects`

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
