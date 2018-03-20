package controller

import (
	"strconv"

	"github.com/zipper-project/zipper/common/log"
	"golang.org/x/net/websocket"
)

type notify struct {
	wsConn     *websocket.Conn
	heightChan chan uint32
}

func newNotify(wsURL, origin string) *notify {
	ws, err := websocket.Dial(wsURL, "", origin)
	if err != nil {
		panic(err)
	}
	return &notify{
		heightChan: make(chan uint32, 10),
		wsConn:     ws,
	}
}

func (n *notify) heightChannel() <-chan uint32 {
	go func() {
		for {
			var msg = make([]byte, 512)
			m, err := n.wsConn.Read(msg)
			if err != nil {
				log.Errorln(err)
			}
			height, _ := strconv.Atoi(string(msg[:m]))
			n.heightChan <- uint32(height)
		}
	}()
	return n.heightChan
}
