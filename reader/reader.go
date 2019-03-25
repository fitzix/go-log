package reader

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/fitzix/go-log/models"
)

// ServerConf ss
var ServerConf models.ServerConf

// Reader reader
type Reader struct {
	// 日志消息
	logs chan string
	file *os.File
}

type XReader interface {
	HandleLog()
	WriteContent(content string)
	Start() error
}

// 收取日志
func (reader *Reader) HandleLog() {
	for {
		rec := <-reader.logs
		if !utf8.ValidString(rec) {
			log.Printf("编码错误: %v", rec)
		} else {
			reader.WriteContent(rec)
		}
	}
}

// WriteContent 向文件内写数据
func (reader *Reader) WriteContent(content string) {
	if reader.file == nil {
		// err := errors.New("")
		var err error
		reader.file, err = os.OpenFile(ServerConf.LogDir+strconv.Itoa(int(time.Now().Unix()))+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Println("创建日志文件失败", err)
			return
		}
		go func() {
			<-time.After(time.Duration(ServerConf.Reader.Interval) * time.Minute)
			err := reader.file.Close()
			if err != nil {
				log.Println("close channel file error", err)
			}
			reader.file = nil
		}()
	}

	if ServerConf.LogType == "json" && strings.Contains(content, "\n") {
		content = strings.Replace(content, "\n", "", -1)
	}

	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	_, err := reader.file.WriteString(content)
	log.Println("write string error", err)
}

func Start() error {
	var readerImpl = getReader()
	return readerImpl.Start()
}

func getReader() XReader {
	if strings.HasPrefix(ServerConf.Reader.Network, "tcp") {
		return &TcpReader{}
	}
	if strings.HasPrefix(ServerConf.Reader.Network, "udp") {
		return &UDPReader{}
	}
	return &HttpReader{}
}
