package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/massn/ManualAccounter/pkg/chart"
	jsonbin "github.com/massn/ManualAccounter/pkg/json"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

type Entry struct {
	Time      time.Time `json:"time,omitempty"`
	Valuation float64   `json:"valuation,omitempty"`
	Gain      float64   `json:"gain,omitempty"`
}

func main() {
	useRemote := flag.Bool("r", true, "use JSONBin data")
	accountFileName := flag.String("a", "", "account json file")
	key := flag.String("k", "", "JSONBin API-key")
	binId := flag.String("b", "", "JSONBin Bin ID")
	flag.Parse()

	var account *[]Entry
	if *useRemote {
		account, _ = getRemoteAcconut(*binId, *key)
	} else {
		account, _ = getLocalAcconut(*accountFileName)
	}

	nowString := time.Now().Format("2006-01-02_03:04:05")
	filePath := fmt.Sprintf("charts/%s.html", nowString)

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

	if *useRemote {
		if err := writeRemoteNewAccount(&newAccount, nowString, *binId, *key); err != nil {
			panic(err)
		}
	} else {
		if err := writeLocalNewAccount(&newAccount, *accountFileName); err != nil {
			panic(err)
		}
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

func getLocalAcconut(accountFileName string) (*[]Entry, error) {
	account := []Entry{}
	bytes, err := ioutil.ReadFile(accountFileName)
	if err != nil {
		return &account, err
	}
	err = json.Unmarshal(bytes, &account)
	if err != nil {
		return &account, err
	}
	return &account, nil
}

func getRemoteAcconut(binId, key string) (*[]Entry, error) {
	account := []Entry{}
	rp := jsonbin.ReadParam{
		BinId:  binId,
		APIKey: key,
	}
	res, err := jsonbin.Read(rp)
	if err != nil {
		return &account, err
	}
	err = json.Unmarshal([]byte(res.Record), &account)
	if err == nil {
	}
	return &account, nil
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

func writeLocalNewAccount(newAccount *[]Entry, accountFileName string) error {
	newBytes, err := json.MarshalIndent(newAccount, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(accountFileName, newBytes, 0666)
}

func writeRemoteNewAccount(newAccount *[]Entry, binName, binId, key string) error {
	newBytes, err := json.MarshalIndent(newAccount, "", "    ")
	if err != nil {
		return err
	}
	up := jsonbin.UpdateParam{
		BinId:      binId,
		Body:       string(newBytes),
		Versioning: false,
		APIKey:     key,
	}

	res, err := jsonbin.Update(up)
	if err != nil {
		panic(err)
	}

	if res.StatusCode == 200 {
		fmt.Println("Updated json bin")
		return nil
	}
	fmt.Printf("Failed to update. Response:%#v\n", res)

	cp := jsonbin.CreateParam{
		BinName:   binName,
		Body:      string(newBytes),
		IsPrivate: true,
		APIKey:    key,
	}
	res, err = jsonbin.Create(cp)

	fmt.Println("Created json bin")
	return err
}
