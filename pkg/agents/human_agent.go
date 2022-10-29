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
	line, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return line[:len(line)-1]
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
		if err != nil || b <= 0 || b > 6 {
			fmt.Printf("Dice must be within [1, 6]. %s is invalid\n", betPts[0])
			h.Play(r)
		}
		bet = game.Bet{uint(a), uint(b)}
	}

	return game.Action{
		T:     at,
		Raise: bet,
	}
}

func (h *HumanAgent) Initialise(r game.Round, dice []uint) {
	h.dice = dice
	fmt.Printf("Your dice: %d\n", h.dice)

}
