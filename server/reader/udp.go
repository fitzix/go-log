package reader

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/apex/log"
	"github.com/fatih/color"
)

// UDPReader udp struct
type UDPReader struct {
	conn *net.UDPConn // UDP连接
	Reader
}

// ReadLog 读取日志并放入channel
func (r *UDPReader) ReadLog() {
	buf := make([]byte, ServerConf.Reader.ReadByte)
	n, _, err := r.conn.ReadFromUDP(buf)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if n > 0 {
		rec := string(buf[:n])
		r.logs <- rec
	}
}

// UDPStart udp server start
func UDPStart() {
	var s UDPReader
	s.logs = make(chan string, ServerConf.Reader.ReadChan)

	udpAddr, err := net.ResolveUDPAddr(ServerConf.Reader.Network, ":"+strconv.Itoa(ServerConf.Reader.Port))
	if err != nil {
		log.Fatalf("解析监听地址失败----> %v", err)
	}
	s.conn, err = net.ListenUDP("udp4", udpAddr)
	if err != nil {
		log.Fatalf("监听端口失败----->%v", err)
	}
	if ServerConf.Reader.ReadBuffer == 0 {
		s.conn.SetReadBuffer(1048576)
	} else {
		s.conn.SetReadBuffer(ServerConf.Reader.ReadBuffer)
	}
	if !strings.HasSuffix(ServerConf.LogDir, "/") {
		ServerConf.LogDir += "/"
	}

	log.Infof(color.CyanString("开始监听%s", s.conn.LocalAddr()))

	defer s.conn.Close()

	go s.HandleLog()

	for {
		s.ReadLog()
	}
}
