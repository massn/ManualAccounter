package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

const accountFileName = "account.json"

type Entry struct {
	Time      time.Time `json:"time,omitempty"`
	Valuation float64   `json:"valuation,omitempty"`
	Gain      float64   `json:"gain,omitempty"`
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("[ERROR] input VALUATION and GAIN as the argument.")
		os.Exit(1)
	}

	account := []Entry{}
	bytes, err := ioutil.ReadFile(accountFileName)
	if err == nil {
		_ = json.Unmarshal(bytes, &account)
	}

	val, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		panic(err)
	}
	gain, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil {
		panic(err)
	}
	now := time.Now()
	newEntry := Entry{Time: now, Valuation: val, Gain: gain}
	newBytes, err := json.MarshalIndent(append(account, newEntry), "", "    ")
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(accountFileName, newBytes, 0666); err != nil {
		panic(err)
	}
}
