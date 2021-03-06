package models

type ServerConf struct {
	Title   string     `toml:"title"`
	LogDir  string     `toml:"log_dir"`
	LogType string     `toml:"log_type"`
	Reader  ReaderConf `toml:"reader"`
	Sender  SenderConf `toml:"sender"`
}

type ReaderConf struct {
	Network    string `toml:"network"`
	Port       int    `toml:"port"`
	Interval   int    `toml:"interval"`
	ReadBuffer int    `toml:"read_buffer"`
	ReadByte   int    `toml:"read_byte"`
	ReadChan   int    `toml:"read_chan"`
}

type SenderConf struct {
	Enabled      bool               `toml:"enabled"`
	ChannelSize  int                `toml:"channel_size"`
	RemoteServer []RemoteServerConf `toml:"remote_server"`
}

type RemoteServerConf struct {
	Protocol    string   `toml:"protocol"`
	Ip          string   `toml:"ip"`
	Port        int      `toml:"port"`
	Weight      int      `toml:"weight"`
	PoolEnabled bool     `toml:"pool_enabled"`
	ConnPool    ConnPool `toml:"conn_pool"`
}

type ConnPool struct {
	InitialCap  int `toml:"initial_cap"`
	MaxCap      int `toml:"max_cap"`
	IdleTimeout int `toml:"idle_timeout"`
}

// 配置文件默认值
var DefaultServerConf = `
title = "udp server 配置文件"
# 日志存储地址
log_dir = "/tmp"
# 日志类型 raw(默认) json(每次一条 校验并去掉换行符)
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
`
