package main

import (
	"fmt"

	"github.com/banbanpeppa/huobi-future-go/websocket"
)

func main() {
	// WebSocket 行情,交易 API
	fmt.Println()
	fmt.Println("********************websocket  Run******************************")

	// WebSocket 行情,交易 API
	// websocket.WSRun()  //无需本地IP地址，直接运行
	p := websocket.NewDefaultParameters()
	huobiClient := websocket.NewHuobiWSClient(p) //配置文件须填写本地IP地址，WS运行太久，外部原因可能断开，支持自动重连

	BTC_CW := "{\"Sub\":\"market.BTC_CW.trade.detail\"}"
	BTC_NW := "{\"Sub\":\"market.BTC_NW.trade.detail\"}"
	BTC_CQ := "{\"Sub\":\"market.BTC_CQ.trade.detail\"}"
	messages := []string{BTC_CW, BTC_NW, BTC_CQ}
	huobiClient.Sub(messages)

	m := websocket.GetMsg()

	for data := range m.Data {
		fmt.Println(data)
	}
}
