package network

import (
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/http"
)

func UpgradeToWebsockets(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	log.Debug().Msg("Received HTTP request")
	u := websocket.Upgrader{}
	return u.Upgrade(w, r, nil)
}

func CloseConn(c *websocket.Conn) {
	err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Warn().Err(err).Msg("Failed to close websocket connection normally.")
		return
	}
	log.Warn().Msg("We did not wait to hear back from the server.")
	// TODO: handle CloseMessage response from server.
}
