package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/yupi/slots/machine"
	"github.com/yupi/slots/machine/atkins"
)

func TestAtkinsMachine(t *testing.T) {
	rand.Seed(45)

	a := &app{
		balance: map[string]int{
			"vasya": 2000,
			"petya": 5000,
		},
		machines: map[string]*machine.MachineSettings{
			"atkins": atkins.DefaultSettings(),
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.handler(w, r)
	}))
	defer ts.Close()

	p := params{
		UID:   "vasya",
		Chips: 2000,
		Lines: 20,
		Bet:   5,
	}

	b, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post(ts.URL, "application/json", bytes.NewBuffer(b))
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var got result
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatal(err)
	}

	expected := result{
		Total: 145,
		Spins: []machine.Spin{
			{
				Type:  "main",
				Total: 145,
				Stops: []int{13, 12, 1, 2, 24}},
		},
		JWT: map[string]interface{}{
			"uid":   "vasya",
			"bet":   5.0,
			"chips": 2045.0,
			"lines": 20.0,
		},
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fail()
	}
}
