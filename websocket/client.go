package websocket

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/banbanpeppa/huobi-future-go/utils"

	"github.com/gorilla/websocket"
	"github.com/spf13/cast"
)

type Message struct {
	Ts      int         `json:"ts"`
	Status  string      `json:"status"`
	ErrCode string      `json:"err-code"`
	ErrMsg  string      `json:"err-msg"`
	Ping    int         `json:"ping"`
	Data    chan string `json:"data"`
}

type Client struct {
	Name   string
	Params *ClientParameters
	Ws     *websocket.Conn
}

// ws-client的监控对象
type Moniter struct {
	clientNum  int      //连接到huobi的ws-client数目
	addChan    chan int // 当添加ws-client的时候，使用管道对象管理
	subChan    chan int // 当关闭了ws-client之后，subChan减1
	lastUseSec int
}

var (
	mon           *Moniter
	clientNameNum int
	Msg           *Message
)

func GetMsg() *Message {
	return Msg
}

func initMoniter() {
	mon = &Moniter{}
	mon.addChan = make(chan int, 1000)
	mon.subChan = make(chan int, 1000)
	go func() {
		for {
			select {
			case <-mon.addChan: //接收管道消息并处理
				mon.clientNum++
			case <-mon.subChan:
				mon.clientNum--
			}
		}
	}()
}

func AddClientNum() {
	mon.addChan <- 1 //发送管道消息
}

func SubClientNum() {
	mon.subChan <- 1
}

func client(params *ClientParameters, name string) *Client {
	return &Client{Name: name, Params: params}
}

func NowSec() int {
	return int(time.Now().UnixNano() / 1000000000)
}

func NewHuobiWSClient(params *ClientParameters) *Client {
	clientNameNum++
	c := client(params, cast.ToString(clientNameNum))
	return c
}

func (cli *Client) Sub(messages []string) {
	initMoniter()
	go cli.sub(messages)
	time.Sleep(1 * time.Millisecond)
	cli.reCreateClient()
	for {
		time.Sleep(10 * time.Second)
	}
}

func (cli *Client) reCreateClient() {
	go func() {
		time.Sleep(time.Second * 100)
		checkTicker := time.NewTicker(time.Second * 20)
		for {
			select {
			case <-checkTicker.C:
				if mon.clientNum <= 0 {
					NewHuobiWSClient(cli.Params)
				}
			}
		}
	}()
}

func (cli *Client) sub(messages []string) {
	AddClientNum()
	dialer := websocket.DefaultDialer
	dialer.NetDial = func(network, addr string) (net.Conn, error) {

		addrs := []string{string(cli.Params.LocalIP)}
		localAddr := &net.TCPAddr{IP: net.ParseIP(addrs[rand.Int()%len(addrs)]), Port: 0}
		d := net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			LocalAddr: localAddr,
			DualStack: true,
		}
		c, err := d.Dial(network, addr)
		return c, err
	}
	c, _, err := dialer.Dial(cli.Params.URL, nil)
	if err != nil {
		log.Println("Dial Erro:", err)
		SubClientNum()
		return
	}
	log.Println(c.LocalAddr().String())

	defer func() {
		c.Close()
		SubClientNum()
	}()

	for _, message := range messages {
		messgeByte := []byte(message)
		err = c.WriteMessage(websocket.TextMessage, messgeByte)
		if err != nil {
			log.Println("write err :", err)
		}
	}

	go func() {
		pangTicker := time.NewTicker(time.Second * 3)
		for {
			select {
			case <-pangTicker.C:
				message := []byte(fmt.Sprintf("{\"pong\":%d}", time.Now().Unix()))
				err = c.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Println("send msg err:", err)
					return
				}
			}
		}
	}()

	for {
		_, zipmsg, err := c.ReadMessage()
		if err != nil {
			log.Println("Read Error : ", err, cli.Name)
			c.Close()
			return
		}

		msg, err := utils.ParseGzip(zipmsg)
		if err != nil {
			log.Println("gzip Error : ", err)
		}
		go func() {
			Msg = &Message{}
			Msg.Data = make(chan string)
			testm := string(msg)
			// log.Println(testm)
			Msg.Data <- string(testm)
		}()

		// log.Println(string(msg[:]))

	}
}
