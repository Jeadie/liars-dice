package network

import (
	"encoding/json"
	"github.com/Jeadie/liars-dice/pkg/game"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
)

func handleIncomingEvents(c *websocket.Conn, output chan game.Event) {
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		var event game.Event
		err = json.Unmarshal(message, &event)
		if err != nil {
			log.Println("Error parsing message: ", err)
			break
		}
		output <- event
		log.Println("recv: ", event)
	}
}
func ConnectToServer(addr string) (*websocket.Conn, error) {
	u := url.URL{Scheme: "ws", Host: addr, Path: "/"}
	log.Printf("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	return c, err
}

func Handle(c *websocket.Conn, incoming chan game.Event, outgoing chan game.Action) {
	done := make(chan interface{})    // Channel to indicate that the receiverHandler is done
	interrupt := make(chan os.Signal) // OS interrupts
	go handleIncomingEvents(c, incoming)

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			CloseConn(c)
		case act := <-outgoing:
			err := sendAction(c, act)
			if err != nil {
				log.Println("write:", err)
				return
			}
		}
	}
}

func sendAction(conn *websocket.Conn, act game.Action) error {
	v, err := json.Marshal(act)
	if err != nil {
		return err
	}
	return conn.WriteMessage(websocket.TextMessage, v)
}

func CloseConn(c *websocket.Conn) {
	err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("write close:", err)
		return
	}
	// TODO: handle CloseMessage response from server.
}
