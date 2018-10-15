package machine

import (
	"fmt"
	"math/rand"
	"strings"
)

type Symbol int8

const EmptySym Symbol = 0

type Spin struct {
	Type  string
	Total int
	Stops []int
}

type Machine struct {
	ReelStrips    [][]Symbol
	Paylines      [][]int
	Pays          [][]int
	SymbolNames   map[Symbol]string
	RetriggerFree bool
	ScatterSymbol Symbol
	WildSymbol    Symbol
	NReels        int
}

func (m *Machine) Run(n, bet int) []Spin {
	stops := make([]int, m.NReels)
	for i := 0; i < m.NReels; i++ {
		stops[i] = rand.Intn(len(m.ReelStrips[i]))
	}

	lines := m.Spin(stops)
	total := m.CalcRTP(lines, n, bet)

	spins := []Spin{
		{
			Type:  "main",
			Total: total,
			Stops: stops,
		},
	}

	return spins
}

func (m *Machine) PrintLines(lines []Symbol) string {
	b := &strings.Builder{}
	for i := 0; i < 3; i++ {
		offset := i * m.NReels
		for j := 0; j < m.NReels; j++ {
			fmt.Fprint(b, m.SymbolNames[lines[offset+j]])
			if j != m.NReels-1 {
				fmt.Fprint(b, ", ")
			}
		}
		fmt.Fprint(b, "\n")
	}

	return b.String()
}

func (m *Machine) Spin(stops []int) []Symbol {
	lines := make([]Symbol, m.NReels*3)

	for i := 0; i < m.NReels; i++ {
		lastStop := len(m.ReelStrips[i]) - 1
		var stop1, stop2, stop3 int

		stop2 = stops[i]
		if stop2 == 0 {
			stop1 = lastStop
			stop3 = 1
		} else if stop2 == lastStop {
			stop1 = lastStop - 1
			stop3 = 0
		} else {
			stop1 = stop2 - 1
			stop3 = stop2 + 1
		}

		lines[i] = m.ReelStrips[i][stop1]
		lines[m.NReels+i] = m.ReelStrips[i][stop2]
		lines[m.NReels*2+i] = m.ReelStrips[i][stop3]
	}

	return lines
}

func (m *Machine) CheckPayline(lines []Symbol, i int) (Symbol, int) {
	startSym := EmptySym
	symCount := 0

	for _, idx := range m.Paylines[i] {
		sym := lines[idx]
		// ignore scatter
		if sym == m.ScatterSymbol {
			break
		}

		if startSym == EmptySym || startSym == m.WildSymbol {
			startSym = sym
		}
		if sym == startSym || sym == m.WildSymbol {
			symCount++
		} else {
			break
		}
	}

	return startSym, symCount
}

func (m *Machine) CalcRTP(lines []Symbol, n, bet int) int {
	scatters := 0
	if m.ScatterSymbol != 0 {
		for _, sym := range lines {
			if sym == m.ScatterSymbol {
				scatters++
			}
		}
	}

	fmt.Printf("lines:\n%s\n", m.PrintLines(lines))
	total := m.Pays[m.ScatterSymbol][scatters]
	for i := 0; i < n; i++ {
		sym, count := m.CheckPayline(lines, i)
		coeff := m.Pays[sym][count]
		if coeff > 0 {
			fmt.Printf("payline %d, sym %s, count %d, coeff %d\n", i+1, m.SymbolNames[sym], count, coeff)
		}
		total += m.Pays[sym][count]
	}

	return total * bet
}
