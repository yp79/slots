package atkins

import (
	"reflect"
	"testing"

	"github.com/yupi/slots/machine"
)

func TestSpin(t *testing.T) {
	sm := New()
	expected := []machine.Symbol{
		Cheese, Mayonnaise, Cheese, Mayonnaise, BuffaloWings,
		Eggs, Ham, Mayonnaise, Bacon, Bacon,
		Atkins, Atkins, Butter, Ham, Scale,
	}

	stops := []int{18, 12, 5, 31, 0}
	result := sm.Spin(stops)

	if !reflect.DeepEqual(result, expected) {
		t.Fail()
		t.Log("got:\n")
		sm.PrintLines(result)
		t.Log("expected:\n")
		sm.PrintLines(expected)
	}
}

func TestCheckPayline(t *testing.T) {
	sm := New()
	tests := []struct {
		lines    []machine.Symbol
		paylines []int
		symbol   []machine.Symbol
		count    []int
	}{
		{
			lines:    sm.Spin([]int{19, 31, 18, 18, 26}),
			paylines: []int{5, 11, 16},
			symbol:   []machine.Symbol{Bacon, Mayonnaise, Bacon},
			count:    []int{3, 3, 3},
		},
		{
			lines:    sm.Spin([]int{10, 8, 1, 21, 5}),
			paylines: []int{8},
			symbol:   []machine.Symbol{Butter},
			count:    []int{5},
		},
		{
			lines:    sm.Spin([]int{14, 13, 25, 0, 27}),
			paylines: []int{12, 1, 5, 16, 17, 13, 15},
			symbol: []machine.Symbol{Mayonnaise, Steak, BuffaloWings, Steak, Steak,
				BuffaloWings, BuffaloWings},
			count: []int{3, 3, 2, 2, 2, 3, 3},
		},
	}

	for _, test := range tests {
		for i, p := range test.paylines {
			sym, count := sm.CheckPayline(test.lines, p-1)
			if sym != test.symbol[i] || count != test.count[i] {
				t.Errorf("payline %d, got %d %s, expected %d %s", p,
					count, sm.SymbolNames[sym],
					test.count[i], sm.SymbolNames[test.symbol[i]])
				t.Log("lines:\n" + sm.PrintLines(test.lines))
			}
		}
	}
}

func TestCalcRPT(t *testing.T) {
	sm := New()
	tests := []struct {
		lines []machine.Symbol
		total int
	}{
		{
			lines: sm.Spin([]int{14, 13, 25, 0, 27}),
			total: 540,
		},
		{
			lines: []machine.Symbol{
				Scale, EmptySym, EmptySym, EmptySym, Scale,
				EmptySym, EmptySym, Scale, EmptySym, EmptySym,
				EmptySym, EmptySym, EmptySym, Scale, EmptySym,
			},
			total: 25 * 5,
		},
	}

	for _, test := range tests {
		total := sm.CalcRTP(test.lines, 20, 5)
		if total != test.total {
			t.Errorf("got %d, expected %d", total, test.total)
		}
	}
}
