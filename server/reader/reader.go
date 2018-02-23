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
	"github.com/fitzix/go-udp/server/utils"
)

var (
	ServerConf models.UdpConf
)

type Reader struct {
	conn  *net.UDPConn        //UDP连接
	logs  chan string         //日志消息
	files map[string]*os.File //用于保存当前已打开的日志文件 file descriptor
}

func (reader *Reader) WriteLog(log, filePath string) error {
	if !strings.HasSuffix(filePath, "/") {
		filePath += "/"
	}
	filename := filePath + utils.TimeIntervalFileName(ServerConf.Reader.Interval) + ".log"
	reader.WriteContent(filename, log)
	return nil
}

// 向文件内写数据
func (reader *Reader) WriteContent(filename, content string) {
	if _, ok := reader.files[filename]; !ok {
		err := errors.New("")
		reader.files[filename], err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			fmt.Println(err)
		}
		//关闭文件，删除file descriptor
		go func() {
			select {
			case <-time.After(time.Duration(ServerConf.Reader.Interval) * 60 * time.Second):
				reader.files[filename].Close()
				delete(reader.files, filename)
			}
		}()
	}
	file := reader.files[filename]
	// file.WriteString(content + "\n")
	file.WriteString(content)
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
func (reader *Reader) HandleLog(filePath string) {
	for {
		rec := <-reader.logs
		reader.WriteLog(rec, filePath)
	}
}

func Start() {
	var s Reader
	s.logs = make(chan string, ServerConf.Reader.ReadChan)

	s.files = make(map[string]*os.File, 5)
	udpAddr, err := net.ResolveUDPAddr("udp4", ":"+strconv.Itoa(ServerConf.Port))
	s.conn, err = net.ListenUDP("udp4", udpAddr)
	if err != nil {
		log.Fatalf("监听端口失败----->", err)
	}
	if ServerConf.Reader.ReadBuffer == 0 {
		s.conn.SetReadBuffer(1048576)
	} else {
		s.conn.SetReadBuffer(ServerConf.Reader.ReadBuffer)
	}
	log.Infof(color.CyanString("开始监听%s", s.conn.LocalAddr()))

	defer s.conn.Close()

	go s.HandleLog(ServerConf.LogDir)

	for {
		s.ReadLog()
	}
}