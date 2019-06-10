package websocket

type ClientParameters struct {
	URL     string
	LocalIP string
}

func NewDefaultParameters() *ClientParameters {
	return &ClientParameters{
		URL:     WS_URL,
		LocalIP: Local_IP,
	}
}
