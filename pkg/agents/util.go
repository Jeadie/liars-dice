package agents

func MakeAgents(n uint, humanIdx int, wsAgents chan *WsAgent, numWsAgents uint) []Agent {
	agents := make([]Agent, n)
	for i := 0; i < int(n); i++ {
		// If humanIdx == -1, no human agent is used.
		if i == humanIdx {
			agents[i] = ConstructHuman()

		} else if numWsAgents > 0 {
			numWsAgents--
			x := <-wsAgents
			agents[i] = x
		} else {
			agents[i] = ConstructProbAgent()
		}
	}
	return agents
}
