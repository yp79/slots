package atkins

import (
	"github.com/yupi/slots/machine"
)

const (
	EmptySym machine.Symbol = machine.EmptySym
	Scale    machine.Symbol = iota
	Atkins
	Steak
	Ham
	BuffaloWings
	Sausage
	Eggs
	Butter
	Cheese
	Bacon
	Mayonnaise
)

var SymbolNames = map[machine.Symbol]string{
	EmptySym:     "no sym",
	Scale:        "scale",
	Atkins:       "atkins",
	Steak:        "steak",
	Ham:          "ham",
	BuffaloWings: "buffalo wings",
	Sausage:      "sausage",
	Eggs:         "eggs",
	Butter:       "butter",
	Cheese:       "cheese",
	Bacon:        "bacon",
	Mayonnaise:   "mayonnaise",
}

func New() *machine.Machine {
	reelStrips := [][]machine.Symbol{
		{
			Scale,
			Mayonnaise,
			Ham,
			Sausage,
			Bacon,
			Eggs,
			Cheese,
			Mayonnaise,
			Sausage,
			Butter,
			BuffaloWings,
			Bacon,
			Eggs,
			Mayonnaise,
			Steak,
			BuffaloWings,
			Butter,
			Cheese,
			Eggs,
			Atkins,
			Bacon,
			Mayonnaise,
			Ham,
			Cheese,
			Eggs,
			Scale,
			Butter,
			Bacon,
			Sausage,
			BuffaloWings,
			Steak,
			Butter,
		},
		{
			Mayonnaise,
			BuffaloWings,
			Steak,
			Sausage,
			Cheese,
			Mayonnaise,
			Ham,
			Butter,
			Bacon,
			Steak,
			Sausage,
			Mayonnaise,
			Ham,
			Atkins,
			Butter,
			Eggs,
			Cheese,
			Bacon,
			Sausage,
			BuffaloWings,
			Scale,
			Mayonnaise,
			Butter,
			Cheese,
			Bacon,
			Eggs,
			BuffaloWings,
			Mayonnaise,
			Steak,
			Ham,
			Cheese,
			Bacon,
		},
		{
			Ham,
			Butter,
			Eggs,
			Scale,
			Cheese,
			Mayonnaise,
			Butter,
			Ham,
			Sausage,
			Bacon,
			Steak,
			BuffaloWings,
			Butter,
			Mayonnaise,
			Cheese,
			Sausage,
			Eggs,
			Bacon,
			Mayonnaise,
			BuffaloWings,
			Ham,
			Sausage,
			Bacon,
			Cheese,
			Eggs,
			Atkins,
			BuffaloWings,
			Bacon,
			Butter,
			Cheese,
			Mayonnaise,
			Steak,
		},
		{
			Ham,
			Cheese,
			Atkins,
			Scale,
			Butter,
			Bacon,
			Cheese,
			Sausage,
			Steak,
			Eggs,
			Bacon,
			Mayonnaise,
			Sausage,
			Cheese,
			Butter,
			Ham,
			Mayonnaise,
			Bacon,
			BuffaloWings,
			Sausage,
			Cheese,
			Eggs,
			Butter,
			BuffaloWings,
			Bacon,
			Mayonnaise,
			Eggs,
			Ham,
			Sausage,
			Steak,
			Mayonnaise,
			Bacon,
		},
		{
			Bacon,
			Scale,
			Steak,
			Ham,
			Cheese,
			Sausage,
			Butter,
			Bacon,
			BuffaloWings,
			Cheese,
			Sausage,
			Ham,
			Butter,
			Steak,
			Mayonnaise,
			Eggs,
			Sausage,
			Ham,
			Atkins,
			Butter,
			BuffaloWings,
			Mayonnaise,
			Eggs,
			Ham,
			Bacon,
			Butter,
			Steak,
			Mayonnaise,
			Sausage,
			Eggs,
			Cheese,
			BuffaloWings,
		},
	}
	paylines := [][]int{
		{5, 6, 7, 8, 9},
		{0, 1, 2, 3, 4},
		{10, 11, 12, 13, 14},
		{0, 6, 12, 8, 4},
		{10, 6, 2, 8, 14},
		{5, 1, 2, 3, 9},
		{5, 11, 12, 13, 9},
		{0, 1, 7, 13, 14},
		{10, 11, 7, 3, 4},
		{5, 1, 7, 13, 9},
		{5, 11, 7, 3, 9},
		{0, 6, 7, 8, 4},
		{10, 6, 7, 8, 14},
		{0, 6, 2, 8, 4},
		{10, 6, 12, 8, 14},
		{5, 6, 2, 8, 9},
		{5, 6, 12, 8, 9},
		{0, 1, 12, 3, 4},
		{10, 11, 2, 13, 14},
		{0, 11, 12, 13, 4},
	}

	pays := [][]int{
		{0, 0, 0, 0, 0, 0},    // no symbol, for testing purposes
		{0, 0, 0, 5, 25, 100}, // Scale
		{0, 0, 5, 50, 500, 5000},
		{0, 0, 3, 40, 200, 1000},
		{0, 0, 2, 30, 150, 500},
		{0, 0, 2, 25, 100, 300},
		{0, 0, 0, 20, 75, 200},
		{0, 0, 0, 20, 75, 200},
		{0, 0, 0, 15, 50, 100},
		{0, 0, 0, 15, 50, 100},
		{0, 0, 0, 10, 25, 50},
		{0, 0, 0, 10, 25, 50},
	}

	return &machine.Machine{
		ReelStrips:    reelStrips,
		Paylines:      paylines,
		Pays:          pays,
		SymbolNames:   SymbolNames,
		RetriggerFree: true,
		ScatterSymbol: Scale,
		WildSymbol:    Atkins,
		NReels:        5,
	}
}