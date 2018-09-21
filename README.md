# go-log
go udp/tcp log collection server

# windows版安装方法

> 解压文件go_log_0.0.8_windows_amd64.zip  
初始化程序,进行到解压文件目录,执行下面命令,进入到解压文件目录,执行下面命令,会在当前目录下生成一个config.toml文件  
`go-log.exe init`

# linux/mac版安装方法

> 解压文件  
mac  
`tar -xf go_log_0.0.8_darwin_amd64.tar.gz`  
linux  
`tar -xf go_log_0.0.8_linux_amd64.tar.gz`  
初始化程序,进入到解压文件目录，执行下面命令，会在当前目录下生成一个config.toml文件  
`./go-log init`

# 修改配置文件

打开配置文件，会看到以下内容

```toml
title = "udp server 配置文件"
# 日志存储地址
log_dir = "/logs"
# 日志类型 raw json
log_type = "raw"

[reader]
# 监听类型 http tcp4 tcp6 udp4 udp6
network = "udp4"
# 监听端口
port = 8888
#时间间隔重新生成文件 单位:分 60min
interval = 60
# (udp有效)读取缓冲区大小 byte
read_buffer = 1048576
# channel 容量(理论上channel容量越大  缓冲性能越好但会消耗更多的内存)
read_chan = 10000
# 一次读取长度(http 无效)
read_byte = 1024
```

主要配置log_dir与port,注意windows下请用\\\\,例如：log_dir = "D:\\\\logs"
其它参数请按需配置

# 运行程序

> windows下请执行go-log.exe文件  
linux/mac下直接执行脚本  
./go-log


