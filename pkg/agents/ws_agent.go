package agents

import (
	"encoding/json"
	"fmt"
	"github.com/Jeadie/liars-poker/pkg/game"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

type agentServer struct {
	wsChan chan WsAgent
}

// ConstructWsAgents by starting a server on `addr` and then upgrading the connection to a websocket.
func ConstructWsAgents(addr string, numAgents uint) chan WsAgent {
	result := make(chan WsAgent, numAgents)
	aSrvr := agentServer{wsChan: result}
	server := http.Server{Addr: addr, Handler: &aSrvr}

	go server.ListenAndServe()

	return result
}

func (aSrvr *agentServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{}
	conn, err := u.Upgrade(w, r, nil)

	if err != nil {
		log.Print("Error during upgrading:", err)
		return
	}

	aSrvr.wsChan <- WsAgent{
		conn:     conn,
		commands: make([]game.Action, 1),
	}
}

type WsAgent struct {
	conn     *websocket.Conn
	commands []game.Action // Returned from the websocket.
}

// ConstructWsAgent by starting a server on `addr` and then upgrading the connection to a websocket.
func ConstructWsAgent(addr string) *WsAgent {
	x := <-ConstructWsAgents(addr, 1)
	return &x
}

// listenForCommands from the websocket connection. Store locally for use when invoked.
func (agent *WsAgent) listenForCommands() {
	for {
		// Get ws message
		_, message, err := agent.conn.ReadMessage()
		if err != nil {
			log.Println("Error during message reading:", err)
			break
		}

		// Unmarshall bytes -> game.Action
		var act game.Action
		err = json.Unmarshal(message, &act)
		if err != nil {
			log.Println("Error parsing message: ", err)
			break
		}

		// Append to command queue. Will be used by agent.Play()
		agent.commands = append(agent.commands, act)
	}
}

// Play Action received via websocket.
func (agent *WsAgent) Play(r game.Round) game.Action {
	var act game.Action

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
		for len(agent.commands) == 0 {
			time.Sleep(time.Second)
		}

		// Call recursive to not reimplement local cache logic.
		agent.Play(r)
	}
	return act
}

// Handle events by forwarding to websocket
func (agent *WsAgent) Handle(e game.Event) {
	data, err := json.Marshal(e)
	if err != nil {
		fmt.Println("Could not marshal Event: ", e, err)
	}
	err = agent.conn.WriteMessage(0, data)
	if err != nil {
		fmt.Println("Could not send Event: ", e, err)
	}
}
