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
