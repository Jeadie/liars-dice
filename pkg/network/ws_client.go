package network

import (
	"encoding/json"
	"github.com/Jeadie/liars-dice/pkg/game"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/url"
	"os"
	"os/signal"
)

func handleIncomingEvents(c *websocket.Conn, output chan game.Event) {
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Warn().Err(err).Msg("Failed to read message from ws socket.")
			return
		}
		var event game.Event
		err = json.Unmarshal(message, &event)
		if err != nil {
			log.Warn().Err(err).Bytes("rawPkt", message).Msg("Invalid payload received from websocket client.")
			break
		}
		log.Debug().Interface("recv_msg", event)
		output <- event
	}
}
func ConnectToServer(addr string) (*websocket.Conn, error) {
	var c *websocket.Conn
	u, err := url.Parse(addr)
	if err != nil {
		return c, err
	}
	c, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	return c, err
}

func Handle(c *websocket.Conn, incoming chan game.Event, outgoing chan game.Action) {
	done := make(chan interface{})    // Channel to indicate that the receiverHandler is done
	interrupt := make(chan os.Signal) // OS interrupts
	signal.Notify(interrupt, os.Interrupt)

	go handleIncomingEvents(c, incoming)

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			CloseConn(c)
			signal.Reset(os.Interrupt)
		case act := <-outgoing:
			err := sendAction(c, act)
			if err != nil {
				log.Warn().Err(err).Msg("Failed to write message to ws socket.")
				return
			}
		}
	}
}

func sendAction(conn *websocket.Conn, act game.Action) error {
	log.Debug().Interface("event", act).Msg("sending to server")
	v, err := json.Marshal(act)
	if err != nil {
		return err
	}
	return conn.WriteMessage(websocket.TextMessage, v)
}
