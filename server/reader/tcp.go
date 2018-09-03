package reader

import (
	"net"
	"strconv"
	"strings"

	"github.com/apex/log"
	"github.com/fatih/color"
)

type TcpReader struct {
	listener *net.TCPListener //TCP连接
	Reader
}

// 读取日志并放入channel
func (r *TcpReader) ReadLog(c net.Conn) {
	defer c.Close()
	buf := make([]byte, ServerConf.Reader.ReadByte)
	for {
		n, err := c.Read(buf)
		if err != nil {
			return
		}
		r.logs <- string(buf[:n])
	}
}

func (r *TcpReader) Start() {
	r.logs = make(chan string, ServerConf.Reader.ReadChan)

	tcpAddr, err := net.ResolveTCPAddr(ServerConf.Reader.Network, ":"+strconv.Itoa(ServerConf.Reader.Port))
	if err != nil {
		log.Fatalf("解析监听地址失败----> %v", err)
	}
	r.listener, err = net.ListenTCP(ServerConf.Reader.Network, tcpAddr)
	if err != nil {
		log.Fatalf("监听端口失败----->%v", err)
	}

	if !strings.HasSuffix(ServerConf.LogDir, "/") {
		ServerConf.LogDir += "/"
	}

	log.Infof(color.CyanString("开始监听%s", r.listener.Addr()))

	defer r.listener.Close()

	go r.HandleLog()

	for {
		c, err := r.listener.Accept()
		if err != nil {
			log.WithError(err).Error(color.RedString("accept 接收失败"))
			continue
		}
		go r.ReadLog(c)
	}
}
