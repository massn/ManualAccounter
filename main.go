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
	Time      time.Time `json:"time"`
	Valuation float64   `json:"valuation"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("[ERROR] input VALUATION as the argument.")
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
	now := time.Now()
	newEntry := Entry{Time: now, Valuation: val}
	newBytes, err := json.MarshalIndent(append(account, newEntry), "", "    ")
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(accountFileName, newBytes, 0666); err != nil {
		panic(err)
	}
}
