package models

type UdpConf struct {
	Title  string     `toml:"title"`
	Port   int        `toml:"port"`
	LogDir string     `toml:"log_dir"`
	Reader ReaderConf `toml:"reader"`
	Sender SenderConf `toml:"sender"`
}

type ReaderConf struct {
	Interval    int  `toml:"interval"`
	ReadBuffer  int  `toml:"read_buffer"`
	ReadeByte   int  `json:"reade_byte"`
	ReadChan    int  `toml:"read_chan"`
	AutoNewline bool `toml:"auto_newline"`
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
# 监听端口
port = 8888
# 日志存储地址
log_dir = "/tmp"

[reader]
#时间间隔重新生成文件 单位:分 60min
interval = 60
# 读取缓冲区大小 byte
read_buffer = 1048576
# channel 容量(理论上channel容量越大  缓冲性能越好但会消耗更多的内存)
read_chan = 10000
# 每行结尾是否追加换行符
auto_newline = false
# 一次读取长度
read_byte = 1024

#发送服务器配置
[sender]
# 是否启用sender 默认不启用
enabled = false

# channel缓存 
channel_size = 50000

[[sender.remote_server]]
protocol = "tcp"
ip = "127.0.0.1"
port = 8080
# 是否启用连接池
pool_enabled = true
# ;连接池配置
[sender.remote_server.conn_pool]
# 初始化连接数
initial_cap = 30
# 最大连接数
max_cap = 50
# 连接失效时间
idle_timeout = 10


[[sender.remote_server]]
protocol = "udp"
ip = "127.0.0.1"
port = 8080
`

//var ServerConf = UdpConf{
//	Title:  "配置",
//	Port:   8888,
//	LogDir: "",
//	Reader: ReaderConf{
//		Interval:   60,
//		ReadBuffer: 1024 * 1024,
//	},
//	Sender: SenderConf{
//		RemoteServer: []RemoteServerConf{
//			{
//				Port:     9000,
//				Protocol: "udp",
//				Ip:       "127.0.0.1",
//			},
//			{
//				Port:     9001,
//				Protocol: "tcp",
//				Ip:       "127.0.0.1",
//			},
//		},
//	},
//}
