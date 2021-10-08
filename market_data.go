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

type stockClose struct {
	symbol     string
	startPrice float64
	endPrice   float64
}

const (
	nseBase string = "https://archives.nseindia.com/content/historical/EQUITIES/"
)

func stockData(symbol string, startTime time.Time, endTime time.Time) stockClose {
	// fetch close price for the start date
	startUrl := bhavUrl(startTime)
	startPrice := reqBhav(symbol, startUrl)

	// fetch close price for the end date
	endUrl := bhavUrl(endTime)
	endPrice := reqBhav(symbol, endUrl)

	return stockClose{
		symbol:     symbol,
		startPrice: startPrice,
		endPrice:   endPrice,
	}

}

func reqBhav(symbol string, bhavUrl string) float64 {
	client := &http.Client{}
	// stream=True
	req, err := http.NewRequest("GET", bhavUrl, nil)
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
	bhavZip := "bhavdata.csv.zip"
	out, err := os.Create(bhavZip)
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
	r, e := zip.OpenReader(bhavZip)
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

func bhavUrl(reqstTime time.Time) string {
	// Create Bhav request url
	dateFmt := formatTime(reqstTime)
	urlFmt := nseBase + dateFmt[5:] + "/" + strings.ToUpper(dateFmt[2:5]) + "/" + "cm" +
		strings.ToUpper(dateFmt) + "bhav.csv.zip"
	return urlFmt
}

func formatTime(reqTime time.Time) string {
	// Format time obj to required url arrangement
	dayMonthStr := reqTime.Format("02Jan")
	yearStr := reqTime.Format("2006")
	dateStr := dayMonthStr + yearStr
	return dateStr
}
