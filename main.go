package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/massn/ManualAccounter/pkg/chart"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

const defaultAccountFileName = "account.json"

type Entry struct {
	Time      time.Time `json:"time,omitempty"`
	Valuation float64   `json:"valuation,omitempty"`
	Gain      float64   `json:"gain,omitempty"`
}

func main() {
	accountFileName := flag.String("a", defaultAccountFileName, "account json file")
	flag.Parse()

	account := getExistingAcconut(*accountFileName)

	filePath := fmt.Sprintf("charts/%s.html", time.Now().Format("2006-01-02_03:04:05"))

	if flag.NArg() == 0 {
		if err := drawAccount(account, filePath); err != nil {
			panic(err)
		}
		fmt.Println("Drawed the existing account to ", filePath)
		os.Exit(0)
	}
	valArg := flag.Arg(0)
	gainArg := flag.Arg(1)

	newEntry, err := getNewEntry(valArg, gainArg)
	if err != nil {
		panic(err)
	}

	newAccount := append(*account, newEntry)

	if err := writeNewAccount(&newAccount, *accountFileName); err != nil {
		panic(err)
	}
	if err := drawAccount(&newAccount, filePath); err != nil {
		panic(err)
	}
	fmt.Println("Drawed the new account to ", filePath)
}

func drawAccount(account *[]Entry, filePath string) error {
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
		filePath,
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
