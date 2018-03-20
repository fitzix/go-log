package reader

import (
	"errors"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/fatih/color"
	"github.com/fitzix/go-udp/server/models"
)

type TcpReader struct {
	listener *net.TCPListener //UDP连接
	logs     chan string      //日志消息
	//files map[string]*os.File //用于保存当前已打开的日志文件 file descriptor
	file *os.File
}

func (reader *TcpReader) WriteLog(log string) error {
	reader.WriteContent(log)
	return nil
}

// 向文件内写数据
func (reader *TcpReader) WriteContent(content string) {
	if reader.file == nil {
		err := errors.New("")
		reader.file, err = os.OpenFile(models.ServerConf.LogDir+strconv.Itoa(int(time.Now().Unix()))+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.WithError(err).Error(color.RedString("创建日志文件失败"))
			return
		}
		go func() {
			select {
			case <-time.After(time.Duration(models.ServerConf.Reader.Interval) * time.Minute):
				reader.file.Close()
				reader.file = nil
			}
		}()
	}

	if !strings.HasSuffix(content, "\n") {
		reader.file.WriteString(content + "\n")
		return
	}

	reader.file.WriteString(content)
}

// 读取日志并放入channel
func (reader *TcpReader) ReadLog(c net.Conn) {
	defer c.Close()
	buf := make([]byte, models.ServerConf.Reader.ReadByte)
	for {
		n, err := c.Read(buf[0:])
		if err != nil {
			return
		}
		reader.logs <- string(buf[:n])
	}
}

// 处理日志 写入文件
func (reader *TcpReader) HandleLog() {
	for {
		rec := <-reader.logs
		reader.WriteLog(rec)
	}
}

func TcpStart() {
	var s TcpReader
	s.logs = make(chan string, models.ServerConf.Reader.ReadChan)

	//s.files = make(map[string]*os.File, 1)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":"+strconv.Itoa(models.ServerConf.Reader.Port))
	if err != nil {
		log.Fatalf("解析监听地址失败----> %v", err)
	}
	s.listener, err = net.ListenTCP(models.ServerConf.Reader.Network, tcpAddr)
	if err != nil {
		log.Fatalf("监听端口失败----->%v", err)
	}

	if !strings.HasSuffix(models.ServerConf.LogDir, "/") {
		models.ServerConf.LogDir += "/"
	}

	log.Infof(color.CyanString("开始监听%s", s.listener.Addr()))

	defer s.listener.Close()

	go s.HandleLog()

	for {
		c, err := s.listener.Accept()
		if err != nil {
			log.WithError(err).Error(color.RedString("accept 接收失败"))
			continue
		}
		go s.ReadLog(c)
	}

}
