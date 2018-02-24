package reader

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/fatih/color"
	"github.com/fitzix/go-udp/server/models"
	"github.com/fitzix/go-udp/server/sender"
)

var (
	ServerConf models.UdpConf
	master     sender.Master
)

type Reader struct {
	conn *net.UDPConn //UDP连接
	logs chan string  //日志消息
	//files map[string]*os.File //用于保存当前已打开的日志文件 file descriptor
	file *os.File
}

func (reader *Reader) WriteLog(log string) error {
	reader.WriteContent(log)
	return nil
}

// 向文件内写数据
func (reader *Reader) WriteContent(content string) {
	if reader.file == nil{
		err := errors.New("")
		reader.file,err = os.OpenFile(ServerConf.LogDir + strconv.Itoa(int(time.Now().Unix())) + ".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil{
			log.WithError(err).Error(color.RedString("创建日志文件失败"))
		}
		go func() {
			select {
			case <-time.After(time.Duration(ServerConf.Reader.Interval) * time.Minute):
				reader.file.Close()
				reader.file = nil
			}
		}()
	}

	if ServerConf.Reader.AutoNewline {
		reader.file.WriteString(content + "\n")
		return
	}
	reader.file.WriteString(content)
}

// 读取日志并放入channel
func (reader *Reader) ReadLog() {
	buf := make([]byte, 1024)
	n, _, err := reader.conn.ReadFromUDP(buf)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if n > 0 {
		rec := string(buf[:n])
		reader.logs <- rec
	}
}

// 收取日志
func (reader *Reader) HandleLog() {
	if ServerConf.Sender.Enabled {
		for {
			rec := <-reader.logs
			master.Request <- rec
			reader.WriteLog(rec)
		}
	} else {
		for {
			rec := <-reader.logs
			reader.WriteLog(rec)
		}
	}
}

func Start() {
	var s Reader
	s.logs = make(chan string, ServerConf.Reader.ReadChan)

	//s.files = make(map[string]*os.File, 1)
	udpAddr, err := net.ResolveUDPAddr("udp4", ":"+strconv.Itoa(ServerConf.Port))
	if err != nil{
		log.Fatalf("解析监听地址失败----> %v",err)
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

	if ServerConf.Sender.Enabled {
		master = sender.NewMaster(ServerConf.Sender)
		go master.HandleSenderLog()
	}

	go s.HandleLog()

	for {
		s.ReadLog()
	}
}
