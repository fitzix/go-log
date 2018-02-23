package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
	"github.com/fitzix/go-udp/server/models"
	"github.com/BurntSushi/toml"
	"log"
	"github.com/fitzix/go-udp/server/utils"
)

type Server struct {
	conn  *net.UDPConn        //UDP连接
	logs  chan string         //日志消息
	files map[string]*os.File //用于保存当前已打开的日志文件 file descriptor
}

var (
	ServerConf models.UdpConf
	configPath = flag.String("c", "config.toml", "指定配置文件位置")
)

func main() {
	flag.Parse()
	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		file, err := os.OpenFile(*configPath, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Fatalln("创建配置文件失败", err)
		}
		defer file.Close()
		if _, err = file.WriteString(models.DefaultServerConf); err != nil {
			log.Fatalln("配置文件写入失败", err)
		}
		fmt.Println("*************************************")
		fmt.Println("已生成默认配置文件                 **")
		fmt.Println("请先修改默认配置文件               **")
		fmt.Println("*************************************")
		os.Exit(1)
	}
	if _, err := toml.DecodeFile(*configPath, &ServerConf); err != nil {
		log.Fatal("解析配置文件失败---->", err)
	}
	ServerStart(ServerConf.Port, ServerConf.LogDir)
}

func (server *Server) WriteLog(log, filePath string) error {
	if !strings.HasSuffix(filePath, "/") {
		filePath += "/"
	}
	filename := filePath + utils.TimeIntervalFileName(ServerConf.Reader.Interval) + ".log"
	server.WriteContent(filename, log)
	return nil
}

func (server *Server) writeToLog(filename, content string) {
	if len(server.files) == 0 {
		err := errors.New("")
		server.files[filename], err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Println("打开日志文件失败", err)
		}
		// 关闭文件
		go func() {
			select {
			case <-time.After(time.Duration(utils.GetTimeInterval(ServerConf.Reader.Interval)) * 60 * time.Second):
				server.files[filename].Close()
				delete(server.files, filename)
			}
		}()
	}
	file := server.files[filename]
	file.WriteString(content)
}

func (server *Server) WriteContent(filename, content string) {
	if _, ok := server.files[filename]; !ok {
		err := errors.New("")
		server.files[filename], err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			fmt.Println(err)
		}
		go func() {
			select {
			case <-time.After(time.Duration(ServerConf.Reader.Interval) * 60 * time.Second):
				server.files[filename].Close()
				delete(server.files, filename)
			}
		}() //10分钟后关闭文件，删除file descriptor
	}
	file := server.files[filename]
	// file.WriteString(content + "\n")
	file.WriteString(content)
}

func (server *Server) ReadLog() {
	buf := make([]byte, 1024)
	n, _, err := server.conn.ReadFromUDP(buf)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if n > 0 {
		rec := string(buf[:n])
		server.logs <- rec
	}
}

func (server *Server) HandleLog(filePath string) {
	for {
		rec := <-server.logs
		server.WriteLog(rec, filePath)
	}
}

func ServerStart(port int, filePath string) {
	var s Server
	s.logs = make(chan string, ServerConf.Reader.ReadChan)

	s.files = make(map[string]*os.File, 1)
	udpAddr, err := net.ResolveUDPAddr("udp4", ":"+strconv.Itoa(port))
	s.conn, err = net.ListenUDP("udp4", udpAddr)
	if err != nil {
		log.Fatalln("监听端口失败----->", err)
	}
	if ServerConf.Reader.ReadBuffer == 0 {
		s.conn.SetReadBuffer(1048576)
	} else {
		s.conn.SetReadBuffer(ServerConf.Reader.ReadBuffer)
	}
	fmt.Println("*************************************")
	fmt.Printf("开始监听%s               **\n", s.conn.LocalAddr())
	defer s.conn.Close()

	go s.HandleLog(filePath)

	for {
		s.ReadLog()
	}
}
