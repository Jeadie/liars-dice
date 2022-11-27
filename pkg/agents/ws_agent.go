package agents

import (
	"encoding/json"
	"github.com/Jeadie/liars-dice/pkg/game"
	"github.com/Jeadie/liars-dice/pkg/network"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

type agentServer struct {
	wsChan chan *WsAgent
}

// ConstructWsAgents by starting a server on `addr` and then upgrading the connection to a websocket.
func ConstructWsAgents(addr string, numAgents uint) chan *WsAgent {
	result := make(chan *WsAgent, numAgents)
	aSrvr := agentServer{wsChan: result}
	server := http.Server{Addr: addr, Handler: &aSrvr}
	go server.ListenAndServe()

	return result
}

func (aSrvr *agentServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("Received HTTP request")
	u := websocket.Upgrader{}
	conn, err := u.Upgrade(w, r, nil)
	log.Debug().Msg("Upgrading connection to ws success")

	if err != nil {
		log.Warn().Err(err).Msg("error occurred when upgrading connection from HTTP to websocket ")
		return
	}

	agent := &WsAgent{
		conn:     conn,
		commands: make([]game.Action, 0),
	}
	go agent.listenForCommands()
	aSrvr.wsChan <- agent
}

type WsAgent struct {
	conn       *websocket.Conn
	commands   []game.Action // Returned from the websocket.
	Terminated bool
}

func CreateWsAgent(conn *websocket.Conn) *WsAgent {
	ws := &WsAgent{
		conn:       conn,
		commands:   make([]game.Action, 0),
		Terminated: false,
	}
	go ws.listenForCommands()
	return ws
}

// listenForCommands from the websocket connection. Store locally for use when invoked.
func (agent *WsAgent) listenForCommands() {
	for !agent.Terminated {
		// Get ws message
		mType, message, err := agent.conn.ReadMessage()
		if err != nil {
			log.Warn().Err(err).Msg("Error reading game.Action from agent")
			break
		}
		log.Debug().Int("messageType", mType).Bytes("msg", message).Send()
		if mType == websocket.CloseMessage {
			agent.Terminated = true
			return
		}

		// Unmarshall bytes -> game.Action
		var act game.Action
		err = json.Unmarshal(message, &act)
		if err != nil {
			log.Warn().Err(err).Bytes("rawPkt", message).Msg("failed to unmarshal message to game.action")
			break
		}
		log.Debug().Interface("received action", act).Msg("GOt one")
		// Append to command queue. Will be used by agent.Play()
		agent.commands = append(agent.commands, act)
	}
}

// Play Action received via websocket.
func (agent *WsAgent) Play(r game.Round) game.Action {
	var act game.Action
	if agent.Terminated {
		return act
	}

	// Can play immediately from local cache.
	if len(agent.commands) > 0 {
		act = agent.commands[0]

		// Pop action from queue.
		if len(agent.commands) == 1 {
			agent.commands = make([]game.Action, 0)
		} else {
			agent.commands = agent.commands[1:]
		}

	} else {
		// Wait for response from ws (polling in goroutine)
		log.Debug().Int("commands", len(agent.commands)).Msg("waiting for action from agent")
		time.Sleep(time.Second)

		// Call recursive to not reimplement local cache logic.
		act = agent.Play(r)
	}
	log.Debug().Interface("action", act).Msg("Playing agent's action")
	return act
}

// Handle events by forwarding to websocket
func (agent *WsAgent) Handle(e game.Event) {
	data, err := json.Marshal(e)
	if err != nil {
		log.Warn().Interface("event", e).Err(err).Msg("Could not marshal Event")
	}
	err = agent.conn.WriteMessage(1, data)
	if err != nil {
		log.Warn().Interface("event", e).Err(err).Msg("Could not send Event")
	}

	if e.EType == game.GameComplete {
		network.CloseConn(agent.conn)
	}
}
