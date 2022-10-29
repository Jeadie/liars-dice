package agents

import (
	"bufio"
	"fmt"
	"github.com/Jeadie/liars-poker/pkg/game"
	"os"
	"strconv"
	"strings"
)

type HumanAgent struct {
	dice []uint
}

func ConstructHuman() *HumanAgent {
	return &HumanAgent{}
}

func GetInput() string {
	fmt.Printf("Your turn >> ")
	sc := bufio.NewScanner(os.Stdin)
	sc.Split(bufio.ScanLines)
	sc.Scan()
	return sc.Text()
}

func (h *HumanAgent) Play(r game.Round) game.Action {
	line := GetInput()
	input := strings.SplitN(line, " ", 2)
	at, err := game.ParseActionType(input[0])

	if err != nil {
		fmt.Printf("Invalid action %s\n", input[0])
		h.Play(r)
	}

	var bet game.Bet
	if at == game.Raise {
		betPts := strings.SplitN(input[1], " ", 2)
		a, err := strconv.Atoi(betPts[0])
		if err != nil || a < 0 {
			fmt.Printf("Bet quantity %s is invalid\n", betPts[0])
			h.Play(r)
		}

		b, err := strconv.Atoi(betPts[1])
		if err != nil {
			fmt.Printf("Bet quantity %s is invalid\n", betPts[1])
			h.Play(r)
		} else if b <= 0 || b > 6 {
			fmt.Printf("Dice must be within [1, 6], not %s", betPts[1])
			h.Play(r)
		}
		bet = game.Bet{uint(a), uint(b)}
	}

	return game.Action{
		T:     at,
		Raise: bet,
	}
}

func (h *HumanAgent) Handle(e game.Event) {
	switch e.EType {
	case game.Turn:
		fmt.Printf("Player %d %s\n", e.Turn.ActionAgent, e.Turn.Action.ToString())
	case game.RoundComplete:
		agentIdx, changeDice := e.RoundComplete.AffectedAgent, e.RoundComplete.ChangeInDice
		if changeDice > 0 {
			fmt.Printf("Player %d gains %d dice\n", agentIdx, changeDice)
		} else {
			fmt.Printf("Player %d loses %d dice(s)\n", agentIdx, -1*changeDice)
		}
	case game.AgentRemoved:
		fmt.Printf("Player %d has been eliminated\n", e.AgentRemoved.AffectedAgent)
	case game.InvalidAction:
		fmt.Printf("%s\n", e.InvalidAction.Err)
	case game.RoundStart:
		h.dice = e.RoundStart.DiceRolled
		fmt.Printf("Your dice: %d\n", h.dice)

	default:
		fmt.Printf("Unknown event %s\n", e.EType)
	}

}

func (h *HumanAgent) Initialise(r game.Round, dice []uint) {

}
