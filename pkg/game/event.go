package game

type EventType string

const (
	Turn          EventType = "turn"
	RoundComplete EventType = "round-complete"
	GameComplete  EventType = "game-complete"
	AgentRemoved  EventType = "agent-removed"
	AgentTurn     EventType = "agent-turn"
	RoundStart    EventType = "round-start"
	InvalidAction EventType = "invalid-action"
	GameStart     EventType = "game-start"
)

type GameCompleteEvent struct {
	Winner Agent
}

type GameStartEvent struct {
	NumDicePerAgent []uint
	AgentIdx        int
}

type InvalidActionEvent struct {
	InvalidAction Action
	Err           string
}

type AgentTurnEvent struct {
}

type RoundStartEvent struct {
	DiceRolled []uint
}

type AgentRemovedEvent struct {
	AffectedAgent Agent
}

type RoundCompleteEvent struct {
	AffectedAgent Agent
	ChangeInDice  int
}

type TurnEvent struct {
	Action      Action
	ActionAgent Agent
}

type Event struct {
	EType EventType `json:"type"`

	// Present based on EventType
	Turn          *TurnEvent          `json:"turn,omitempty"`
	RoundComplete *RoundCompleteEvent `json:"round-complete,omitempty"`
	GameComplete  *GameCompleteEvent  `json:"game-complete,omitempty"`
	AgentRemoved  *AgentRemovedEvent  `json:"agent-removed,omitempty"`
	InvalidAction *InvalidActionEvent `json:"invalid-action,omitempty"`
	RoundStart    *RoundStartEvent    `json:"round-start,omitempty"`
	GameStart     *GameStartEvent     `json:"game-start,omitempty"`
	AgentTurn     *AgentTurnEvent     `json:"agent-turn,omitempty"`
}
