package reader

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/apex/log"
	"github.com/fatih/color"
	"github.com/fitzix/go-log/server/models"
)

// ServerConf ss
var ServerConf models.SerConf

// Reader reader
type Reader struct {
	logs chan string //日志消息
	file *os.File
}

type XReader interface {
	HandleLog()
	WriteContent(content string)
	Start()
}

// 收取日志
func (reader *Reader) HandleLog() {
	for {
		rec := <-reader.logs
		if !utf8.ValidString(rec) {
			log.Errorf("编码错误: %v", rec)
		} else {
			reader.WriteContent(rec)
		}
	}
}

// WriteContent 向文件内写数据
func (reader *Reader) WriteContent(content string) {
	if reader.file == nil {
		err := errors.New("")
		reader.file, err = os.OpenFile(ServerConf.LogDir+strconv.Itoa(int(time.Now().Unix()))+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.WithError(err).Error(color.RedString("创建日志文件失败"))
			return
		}
		go func() {
			select {
			case <-time.After(time.Duration(ServerConf.Reader.Interval) * time.Minute):
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

func Start() {
	var readerImpl = getReader()
	readerImpl.Start()
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
