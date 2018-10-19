package machine

import (
	"math/rand"
)

type Symbol int8

const EmptySym Symbol = 0

type Spin struct {
	Type  string
	Total int
	Stops []int
}

type MachineSettings struct {
	ReelStrips    [][]Symbol
	Paylines      [][]int
	Pays          [][]int
	SymbolNames   map[Symbol]string
	FreeSpins     int
	FreeSpinsMult int
	RetriggerFree bool
	ScatterSymbol Symbol
	WildSymbol    Symbol
	ReelsCount    int
}

type Machine struct {
	ms       *MachineSettings
	paylines int
	bet      int
	spinFunc func(*Machine) ([]int, []Symbol)
}

func New(ms *MachineSettings, paylines, bet int, sf func(*Machine) ([]int, []Symbol)) *Machine {
	if sf == nil {
		sf = DefaultSpinFunc
	}

	return &Machine{
		ms:       ms,
		paylines: paylines,
		bet:      bet,
		spinFunc: sf,
	}
}

func DefaultSpinFunc(m *Machine) ([]int, []Symbol) {
	stops := make([]int, m.ms.ReelsCount)
	for i := 0; i < m.ms.ReelsCount; i++ {
		stops[i] = rand.Intn(len(m.ms.ReelStrips[i]))
	}
	return stops, m.Spin(stops)
}

func (m *Machine) Run() []Spin {
	freeSpins, spin := m.run("main", 1)
	spins := []Spin{spin}

	for freeSpins != 0 {
		freeSpins--
		fs, spin := m.run("free", m.ms.FreeSpinsMult)
		spins = append(spins, spin)
		if m.ms.RetriggerFree {
			freeSpins += fs
		}
	}
	return spins
}

func (m *Machine) run(typ string, ml int) (int, Spin) {
	freeSpins := 0
	stops, lines := m.spinFunc(m)
	coeff := m.ScatterCoeff(lines)
	if coeff > 0 {
		freeSpins += m.ms.FreeSpins
	}
	coeff += m.PaylinesCoeff(lines)

	return freeSpins, Spin{typ, coeff * m.bet * ml, stops}
}

func (m *Machine) Spin(stops []int) []Symbol {
	lines := make([]Symbol, m.ms.ReelsCount*3)

	for i := 0; i < m.ms.ReelsCount; i++ {
		lastStop := len(m.ms.ReelStrips[i]) - 1
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

		lines[i] = m.ms.ReelStrips[i][stop1]
		lines[m.ms.ReelsCount+i] = m.ms.ReelStrips[i][stop2]
		lines[m.ms.ReelsCount*2+i] = m.ms.ReelStrips[i][stop3]
	}

	return lines
}

func (m *Machine) CheckPayline(lines []Symbol, i int) (Symbol, int) {
	startSym := EmptySym
	symCount := 0

	for _, idx := range m.ms.Paylines[i] {
		sym := lines[idx]
		// ignore scatter
		if sym == m.ms.ScatterSymbol {
			break
		}

		if startSym == EmptySym || startSym == m.ms.WildSymbol {
			startSym = sym
		}
		if sym == startSym || sym == m.ms.WildSymbol {
			symCount++
		} else {
			break
		}
	}

	return startSym, symCount
}

func (m *Machine) ScatterCoeff(lines []Symbol) int {
	scatters := 0
	if m.ms.ScatterSymbol != 0 {
		for _, sym := range lines {
			if sym == m.ms.ScatterSymbol {
				scatters++
			}
		}
	}
	return m.ms.Pays[m.ms.ScatterSymbol][scatters]
}

func (m *Machine) PaylinesCoeff(lines []Symbol) int {
	coeff := 0
	for i := 0; i < m.paylines; i++ {
		sym, count := m.CheckPayline(lines, i)
		coeff += m.ms.Pays[sym][count]
	}

	return coeff
}
