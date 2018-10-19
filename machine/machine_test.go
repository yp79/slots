package machine_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/yupi/slots/machine"
	"github.com/yupi/slots/machine/atkins"
)

var testMachine = atkins.DefaultSettings()

func printLines(lines []machine.Symbol, ms *machine.MachineSettings) string {
	b := &strings.Builder{}
	for i := 0; i < 3; i++ {
		offset := i * ms.ReelsCount
		for j := 0; j < ms.ReelsCount; j++ {
			fmt.Fprint(b, ms.SymbolNames[lines[offset+j]])
			if j != ms.ReelsCount-1 {
				fmt.Fprint(b, ", ")
			}
		}
		fmt.Fprint(b, "\n")
	}

	return b.String()
}

func TestSpin(t *testing.T) {
	ms := testMachine
	m := machine.New(ms, 20, 5, nil)

	expected := []machine.Symbol{
		atkins.Cheese,
		atkins.Mayonnaise,
		atkins.Cheese,
		atkins.Mayonnaise,
		atkins.BuffaloWings,
		atkins.Eggs,
		atkins.Ham,
		atkins.Mayonnaise,
		atkins.Bacon,
		atkins.Bacon,
		atkins.Atkins,
		atkins.Atkins,
		atkins.Butter,
		atkins.Ham,
		atkins.Scale,
	}

	stops := []int{18, 12, 5, 31, 0}
	lines := m.Spin(stops)

	if !reflect.DeepEqual(lines, expected) {
		t.Fail()
		t.Log("got:\n")
		t.Log(printLines(lines, ms))
		t.Log("expected:\n")
		t.Log(printLines(expected, ms))
	}
}

func TestCheckPayline(t *testing.T) {
	ms := testMachine
	m := machine.New(ms, 20, 5, nil)

	tests := []struct {
		lines    []machine.Symbol
		paylines []int
		symbol   []machine.Symbol
		count    []int
	}{
		{
			lines:    m.Spin([]int{19, 31, 18, 18, 26}),
			paylines: []int{5, 11, 16},
			symbol: []machine.Symbol{
				atkins.Bacon,
				atkins.Mayonnaise,
				atkins.Bacon},
			count: []int{3, 3, 3},
		},
		{
			lines:    m.Spin([]int{10, 8, 1, 21, 5}),
			paylines: []int{8},
			symbol:   []machine.Symbol{atkins.Butter},
			count:    []int{5},
		},
		{
			lines:    m.Spin([]int{14, 13, 25, 0, 27}),
			paylines: []int{12, 1, 5, 16, 17, 13, 15},
			symbol: []machine.Symbol{
				atkins.Mayonnaise,
				atkins.Steak,
				atkins.BuffaloWings,
				atkins.Steak,
				atkins.Steak,
				atkins.BuffaloWings,
				atkins.BuffaloWings,
			},
			count: []int{3, 3, 2, 2, 2, 3, 3},
		},
	}

	for _, test := range tests {
		for i, p := range test.paylines {
			sym, count := m.CheckPayline(test.lines, p-1)
			if sym != test.symbol[i] || count != test.count[i] {
				t.Errorf("payline %d, got %d %s, expected %d %s", p,
					count, ms.SymbolNames[sym],
					test.count[i], ms.SymbolNames[test.symbol[i]])
				t.Log("lines:\n" + printLines(test.lines, ms))
			}
		}
	}
}

func TestRTP(t *testing.T) {
	ms := testMachine
	m := machine.New(ms, 20, 5, nil)

	tests := []struct {
		lines []machine.Symbol
		total int
	}{
		{
			lines: m.Spin([]int{14, 13, 25, 0, 27}),
			total: 540,
		},
		{
			lines: []machine.Symbol{
				atkins.Scale,
				atkins.EmptySym,
				atkins.EmptySym,
				atkins.EmptySym,
				atkins.Scale,
				atkins.EmptySym,
				atkins.EmptySym,
				atkins.Scale,
				atkins.EmptySym,
				atkins.EmptySym,
				atkins.EmptySym,
				atkins.EmptySym,
				atkins.EmptySym,
				atkins.Scale,
				atkins.EmptySym,
			},
			total: 25 * 5,
		},
		{
			lines: m.Spin([]int{18, 12, 5, 31, 0}),
			total: 305,
		},
	}

	for _, test := range tests {
		total := m.ScatterCoeff(test.lines) * 5
		total += m.PaylinesCoeff(test.lines) * 5
		if total != test.total {
			t.Errorf("got %d, expected %d", total, test.total)
		}
	}
}

func TestFreeSpin(t *testing.T) {
	ms := testMachine
	ms.FreeSpins = 2

	spinFunc := func() func(m *machine.Machine) ([]int, []machine.Symbol) {
		m := machine.New(ms, 20, 5, nil)
		spinLines := [][]machine.Symbol{
			{
				atkins.Scale,
				atkins.EmptySym,
				atkins.EmptySym,
				atkins.EmptySym,
				atkins.Scale,
				atkins.EmptySym,
				atkins.EmptySym,
				atkins.Scale,
				atkins.EmptySym,
				atkins.EmptySym,
				atkins.EmptySym,
				atkins.EmptySym,
				atkins.EmptySym,
				atkins.Scale,
				atkins.EmptySym,
			},
			m.Spin([]int{18, 12, 5, 31, 0}),
			m.Spin([]int{14, 13, 25, 0, 27}),
		}

		return func(m *machine.Machine) ([]int, []machine.Symbol) {
			var line []machine.Symbol
			line, spinLines = spinLines[0], spinLines[1:]
			return []int{0, 0, 0, 0, 0}, line
		}
	}

	m := machine.New(ms, 20, 5, spinFunc())

	spins := m.Run()
	expected := []machine.Spin{
		{"main", 125, []int{0, 0, 0, 0, 0}},
		{"free", 915, []int{0, 0, 0, 0, 0}},
		{"free", 1620, []int{0, 0, 0, 0, 0}},
	}

	if !reflect.DeepEqual(spins, expected) {
		t.Errorf("got: %v\nexpected: %v\n", spins, expected)
	}
}
