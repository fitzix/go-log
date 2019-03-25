package reader

import (
	"log"
	"net"
	"strconv"
	"strings"
)

type TcpReader struct {
	listener *net.TCPListener
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

func (r *TcpReader) Start() (err error) {
	r.logs = make(chan string, ServerConf.Reader.ReadChan)

	tcpAddr, err := net.ResolveTCPAddr(ServerConf.Reader.Network, ":"+strconv.Itoa(ServerConf.Reader.Port))
	if err != nil {
		log.Printf("解析监听地址失败----> %v", err)
		return err
	}
	r.listener, err = net.ListenTCP(ServerConf.Reader.Network, tcpAddr)
	if err != nil {
		log.Printf("监听端口失败----->%v", err)
		return err
	}

	if !strings.HasSuffix(ServerConf.LogDir, "/") {
		ServerConf.LogDir += "/"
	}
	log.Printf("开始监听%s", r.listener.Addr())

	defer func() {
		err = r.listener.Close()
	}()

	go r.HandleLog()

	for {
		c, err := r.listener.Accept()
		if err != nil {
			log.Println("accept 接收失败", err)
			continue
		}
		go r.ReadLog(c)
	}
}
