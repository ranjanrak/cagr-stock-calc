package cagrcalculator

import (
	"archive/zip"
	"encoding/csv"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type StockClose struct {
	Symbol     string
	StartPrice float64
	EndPrice   float64
}

const (
	NseBase string = "https://archives.nseindia.com/content/historical/EQUITIES/"
)

func StockData(symbol string, startTime time.Time, endTime time.Time) StockClose {
	// fetch close price for the start date
	StartUrl := BhavUrl(startTime)
	StartPrice := ReqBhav(symbol, StartUrl)

	// fetch close price for the end date
	EndUrl := BhavUrl(endTime)
	EndPrice := ReqBhav(symbol, EndUrl)

	return StockClose{
		Symbol:     symbol,
		StartPrice: StartPrice,
		EndPrice:   EndPrice,
	}

}

func ReqBhav(symbol string, BhavUrl string) float64 {
	client := &http.Client{}
	// stream=True
	req, err := http.NewRequest("GET", BhavUrl, nil)
	if err != nil {
		log.Fatal("Request preparation failed", err)
	}
	req.Header.Set("User-Agent",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36")
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal("Request failed", err)
	}
	defer resp.Body.Close()

	// to-do : Unzip in memory without downloading

	// Create a zip file
	BhavZip := "bhavdata.csv.zip"
	out, err := os.Create(BhavZip)
	if err != nil {
		log.Fatal("Zip file not created", err)
	}
	defer out.Close()

	// write downloaded zip to created zip file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatal("Error reading zip file", err)
	}

	// read zip file
	r, e := zip.OpenReader(BhavZip)
	if e != nil {
		panic(e)
	}
	defer r.Close()

	rc, err := r.File[0].Open()
	if err != nil {
		log.Fatal(err)
	}
	// create and copy value to csv file
	outCsv, err := os.Create(r.File[0].Name)
	_, err = io.Copy(outCsv, rc)
	if err != nil {
		log.Fatal(err)
	}
	rc.Close()

	// read csv file
	csvFile, err := os.Open(r.File[0].Name)
	if err != nil {
		log.Fatal("Unable to read input file ", err)
	}
	defer csvFile.Close()

	// filter out request symbol close value
	csvReader := csv.NewReader(csvFile)
	records, err := csvReader.ReadAll()
	var closeValue float64
	for _, value := range records {
		if value[0] == symbol {
			closeValue, err = strconv.ParseFloat(value[6], 8)
			if err != nil {
				log.Fatal("Unable to parse str to float ", err)
			}
			break
		}
	}
	return closeValue

}

func BhavUrl(ReqstTime time.Time) string {
	// Create Bhav request url
	DateFmt := FormatTime(ReqstTime)
	UrlFmt := NseBase + DateFmt[5:] + "/" + strings.ToUpper(DateFmt[2:5]) + "/" + "cm" +
		strings.ToUpper(DateFmt) + "bhav.csv.zip"
	return UrlFmt
}

func FormatTime(ReqTime time.Time) string {
	// Format time obj to required url arrangement
	DayMonthStr := ReqTime.Format("02Jan")
	YearStr := ReqTime.Format("2006")
	DateStr := DayMonthStr + YearStr
	return DateStr
}
