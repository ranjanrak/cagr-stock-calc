package cagrcalculator

import (
	"math"
	"time"
)

type CagrR struct {
	symbol  string
	cagrper float64
}

func CagrCal(symbol string, startTime time.Time, endTime time.Time) CagrR {
	priceDetail := StockData(symbol, startTime, endTime)
	// calculate time delta in years for cagr calculation
	time_diff := endTime.Sub(startTime).Hours() / (24 * 365)
	cagr_decimal := math.Pow((priceDetail.EndPrice/priceDetail.StartPrice), (1.0/time_diff)) - 1
	carg_per := cagr_decimal * 100
	cagrRes := CagrR{
		symbol:  symbol,
		cagrper: math.Round(carg_per*100) / 100,
	}
	return cagrRes

}
