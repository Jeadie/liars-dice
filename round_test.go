package liars_dice

import "testing"

func TestBet_Raises(t *testing.T) {
	a := Bet{4, 2}
	b := Bet{3, 2}
	if !a.Raises(b) {
		t.Errorf("%d %ds raises %d %ds", a[0], a[1], b[0], b[1])
	}

	a = Bet{2, 4}
	b = Bet{2, 3}
	if !a.Raises(b) {
		t.Errorf("%d %ds raises %d %ds", a[0], a[1], b[0], b[1])
	}

	a = Bet{2, 2}
	b = Bet{2, 2}
	if a.Raises(b) {
		t.Errorf("%d %ds does not raise %d %ds", a[0], a[1], b[0], b[1])
	}
}

func TestRound_Calls(t *testing.T) {
	r := &Round{
		dice:    [][]uint{{2, 2}, {1, 2}},
		usedOne: false,
		bet:     Bet{3, 2},
		i:       0,
	}
	result, err := r.Calls(0)
	if err != nil {
		t.Errorf("Agent %d can call on %v", 0, r)
	}
	if !result {
		t.Errorf("Call on %v should be successful", r)
	}

	r = &Round{
		dice:    [][]uint{{2, 2}, {1, 2}},
		usedOne: false,
		bet:     Bet{4, 2},
		i:       0,
	}
	result, err = r.Calls(0)
	if err != nil {
		t.Errorf("Agent %d can call on %v", 0, r)
	}
	if !result {
		t.Errorf("Call on %v should be successful", r)
	}
}

func TestRound_Calls_withOneBetted(t *testing.T) {
	r := &Round{
		dice:    [][]uint{{2, 2}, {1, 2}},
		usedOne: false,
		bet:     Bet{0, 0},
		i:       0,
	}
	err := r.Raise(0, Bet{1, 1})
	if err != nil {
		t.Errorf("Agent %d can Raise", 0)
	}

	err = r.Raise(1, Bet{4, 2})
	if err != nil {
		t.Errorf("Agent %d can Raise", 1)
	}
	result, err := r.Calls(0)

	if err != nil {
		t.Errorf("Agent %d can call on %v. Error: %s", 0, r, err)
	}
	if result {
		t.Errorf("call on %v should be unsuccessful. One already used in round", r)
	}
}

func TestRound_Exact(t *testing.T) {
	r := &Round{
		dice:    [][]uint{{2, 2}, {1, 2}},
		usedOne: false,
		bet:     Bet{3, 2},
		i:       0,
	}
	result, _ := r.Exact(0)
	if result {
		t.Errorf("exact on %v should be successful", r)
	}

	r = &Round{
		dice:    [][]uint{{2, 2}, {1, 2}},
		usedOne: false,
		bet:     Bet{4, 2},
		i:       0,
	}
	result, _ = r.Calls(0)
	if !result {
		t.Errorf("exact on %v should be successful", r)
	}

	r = &Round{
		dice:    [][]uint{{2, 2}, {1, 2}},
		usedOne: true,
		bet:     Bet{4, 2},
		i:       0,
	}
	result, _ = r.Calls(0)
	if result {
		t.Errorf("exact on %v should be unsuccessful. One already used in round", r)
	}
}
