package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/yupi/slots/machine"
	"github.com/yupi/slots/machine/atkins"
)

type params struct {
	UID   string
	Chips int
	Lines int
	Bet   int
}

type result struct {
	Total int
	Spins []machine.Spin
	JWT   map[string]interface{}
}

type app struct {
	balanceMu sync.Mutex
	balance   map[string]int
	machines  map[string]*machine.MachineSettings
}

func (a *app) updateBalance(uid string, amount int) (int, error) {
	a.balanceMu.Lock()
	defer a.balanceMu.Unlock()

	var balance int
	var ok bool
	if balance, ok = a.balance[uid]; !ok {
		return 0, errors.New("no such user")
	}

	newBalance := balance + amount
	if newBalance < 0 {
		return newBalance, errors.New("insufficient chips")
	}

	a.balance[uid] = newBalance
	return newBalance, nil
}

func (a *app) handler(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)

	var params params
	if err := dec.Decode(&params); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	_, err := a.updateBalance(params.UID, -params.Bet*params.Lines)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	m := machine.New(a.machines["atkins"], params.Lines, params.Bet, nil)
	spins := m.Run()
	total := 0
	for _, s := range spins {
		total += s.Total
	}

	balance, err := a.updateBalance(params.UID, total)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	result := result{
		Total: total,
		Spins: spins,
		JWT: map[string]interface{}{
			"uid":   params.UID,
			"chips": balance,
			"bet":   params.Bet,
			"lines": params.Lines,
		},
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(result); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
}

func main() {
	a := &app{
		balance: map[string]int{
			"vasya": 2000,
			"petya": 5000,
		},
		machines: map[string]*machine.MachineSettings{
			"atkins": atkins.DefaultSettings(),
		},
	}
	http.HandleFunc("/api/machines/atkins-diet/", a.handler)
	http.ListenAndServe(":8080", nil)
}
