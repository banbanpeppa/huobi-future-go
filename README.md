# huobi-future-go
Huobi-Futures-Go

## Installation
```
go get github.com/banbanpeppa/huobi-future-go
```

## Usage

### Basic requests
```
var TICKER_ALL = []string{"BTC", "ETH", "BCH", "EOS", "LTC", "ETC", "BSV", "XRP"}

params := websocket.NewDefaultParameters()
huobiClient := websocket.NewHuobiWSClient(params)

//定义想要订阅的ws请求体
requests := []websocket.Request{}
for _, ticker := range TICKER_ALL {
    req_cw := websocket.Request{Id: "id7", Sub: "market." + ticker + "_CW.trade.detail"}
    req_nw := websocket.Request{Id: "id7", Sub: "market." + ticker + "_NW.trade.detail"}
    req_cq := websocket.Request{Id: "id7", Sub: "market." + ticker + "_CQ.trade.detail"}
    requests = append(requests, req_cw, req_nw, req_cq)
}

huobiClient.Subscribe(requests)
for obj := range huobiClient.Listen() {
    switch obj.(type) {
    case []byte:
        tradeDetail := &websocket.TradeDetail{}
        err := json.Unmarshal(obj.([]byte), &tradeDetail)
        if err == nil {
            if len(tradeDetail.Tick.Data) > 0 {
                price := tradeDetail.Tick.Data[0].Price
                fmt.Println(tradeDetail.Ch+": ", price)
            }
        }
    }
    default:
        fmt.Println("other type of obj")
}
```
### Huobi Websocket API

[合约Websocket 文档](https://github.com/huobiapi/API_Docs/wiki/WS_api_reference_Derivatives)