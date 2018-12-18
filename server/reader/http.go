package reader

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/apex/log"
	"github.com/fatih/color"
)

type HttpReader struct {
	Reader
}

func (r *HttpReader) Start() {
	r.logs = make(chan string, ServerConf.Reader.ReadChan)

	if !strings.HasSuffix(ServerConf.LogDir, "/") {
		ServerConf.LogDir += "/"
	}

	go r.HandleLog()

	log.Infof(color.CyanString("开始监听%d", ServerConf.Reader.Port))

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if req.Method == "POST" {
			result, _ := ioutil.ReadAll(req.Body)
			r.logs <- string(result)
		} else {
			w.Write([]byte("请使用POST发送日志!"))
		}
	})

	if err := http.ListenAndServe(":"+strconv.Itoa(ServerConf.Reader.Port), nil); err != nil {
		log.Fatalf("监听端口失败----->%v", err)
	}
}
