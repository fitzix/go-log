package sender

import (
	"log"
	"net"
	"strconv"
	"time"

	"github.com/fitzix/go-log/models"
	"github.com/fitzix/go-log/utils"
	"github.com/fitzix/go-log/utils/pool"
)

type Request struct {
	fn func() int
	c  chan int
}

type Worker struct {
	pool   pool.Pool
	state  bool
	weight int
	msg    chan string
}

type Master struct {
	Workers []*Worker
	Request chan string
}

func NewMaster(serverConf models.SenderConf) (master Master) {
	if serverConf.ChannelSize == 0 {
		master.Request = make(chan string, 10000)
	} else {
		master.Request = make(chan string, serverConf.ChannelSize)
	}
	for _, value := range serverConf.RemoteServer {
		conn, err := resolveLister(value)
		if err != nil {
			log.Println("连接", value.Ip, value.Port, "失败")
			continue
		}
		master.Workers = append(master.Workers, &Worker{conn, true, value.Weight, make(chan string, 10000)})
	}
	return
}

func resolveLister(remote models.RemoteServerConf) (pool pool.Pool, err error) {
	pool, err = createPool(remote)
	return
}

func createPool(remote models.RemoteServerConf) (p pool.Pool, err error) {
	factory := func() (interface{}, error) { return net.Dial(remote.Protocol, remote.Ip+":"+strconv.Itoa(remote.Port)) }
	closePool := func(v interface{}) error { return v.(net.Conn).Close() }
	if remote.Protocol != "tcp" || !remote.PoolEnabled {
		remote.ConnPool.IdleTimeout = 0
		remote.ConnPool.MaxCap = 1
		remote.ConnPool.InitialCap = 1
	}
	poolConfig := &pool.PoolConfig{
		InitialCap: remote.ConnPool.InitialCap,
		MaxCap:     remote.ConnPool.MaxCap,
		Factory:    factory,
		Close:      closePool,
		// 链接最大空闲时间，超过该时间的链接 将会关闭，可避免空闲时链接EOF，自动失效的问题
		IdleTimeout: time.Duration(remote.ConnPool.IdleTimeout) * time.Second,
	}
	p, err = pool.NewChannelPool(poolConfig)
	return
}

func (m *Master) HandleSenderLog() {
	balance := &utils.W1{}
	for _, value := range m.Workers {
		balance.Add(value, value.weight)
		go value.handleLog()
	}
	for {
		worker := balance.Next().(*Worker)
		worker.msg <- <-m.Request
	}
}

func (w *Worker) handleLog() {
	for {
		conn, err := w.pool.Get()
		if err != nil {
			log.Println("sender:88 获取连接池连接失败", err)
			continue
		}
		_, err = conn.(net.Conn).Write([]byte(<-w.msg))
		if err != nil {
			log.Println("sender:93 socker写入数据失败", err)
		}
	}
}

func (m *Master) Release() {
	for _, value := range m.Workers {
		value.pool.Release()
	}
}
