package main

import (
	"encoding/json"
	"fmt"
	"github.com/massn/ManualAccounter/pkg/chart"
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
	account := getExistingAcconut(accountFileName)
	if len(os.Args) != 3 {
		if err := drawAccount(account); err != nil {
			panic(err)
		}
		fmt.Println("Drawed the existing account.")
		os.Exit(0)
	}
	valArg := os.Args[1]
	gainArg := os.Args[2]

	newEntry, err := getNewEntry(valArg, gainArg)
	if err != nil {
		panic(err)
	}

	newAccount := append(*account, newEntry)

	if err := writeNewAccount(&newAccount, accountFileName); err != nil {
		panic(err)
	}
	if err := drawAccount(&newAccount); err != nil {
		panic(err)
	}
}

func drawAccount(account *[]Entry) error {
	gainData := []chart.PointData{}
	valuationData := []chart.PointData{}
	for _, entry := range *account {
		date := entry.Time.Format("2006-01-02")
		gainData = append(gainData, chart.PointData{Date: date, Value: entry.Gain})
		valuationData = append(valuationData, chart.PointData{Date: date, Value: entry.Valuation})
	}
	return chart.Render(
		chart.SeriesData{Name: "Gain", ChartData: gainData},
		chart.SeriesData{Name: "Valuation", ChartData: valuationData},
	)
}

func getExistingAcconut(accountFileName string) *[]Entry {
	account := []Entry{}
	bytes, err := ioutil.ReadFile(accountFileName)
	if err == nil {
		_ = json.Unmarshal(bytes, &account)
	}
	return &account
}

func getNewEntry(valArg, gainArg string) (Entry, error) {
	val, err := strconv.ParseFloat(valArg, 64)
	if err != nil {
		return Entry{}, err
	}
	gain, err := strconv.ParseFloat(gainArg, 64)
	if err != nil {
		return Entry{}, err
	}
	now := time.Now()
	return Entry{Time: now, Valuation: val, Gain: gain}, nil
}

func writeNewAccount(newAccount *[]Entry, accountFileName string) error {
	newBytes, err := json.MarshalIndent(newAccount, "", "    ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(accountFileName, newBytes, 0666); err != nil {
		return err
	}
	return nil
}
