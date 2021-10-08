package cagrcalculator

import (
	"math"
	"time"
)

type cagrR struct {
	symbol  string
	cagrPer float64
}

func CagrCal(symbol string, startTime time.Time, endTime time.Time) cagrR {
	priceDetail := stockData(symbol, startTime, endTime)
	// calculate time delta in years for cagr calculation
	timeDiff := endTime.Sub(startTime).Hours() / (24 * 365)
	cagrDecimal := math.Pow((priceDetail.endPrice/priceDetail.startPrice), (1.0/timeDiff)) - 1
	cargPer := cagrDecimal * 100
	cagrRes := cagrR{
		symbol:  symbol,
		cagrPer: math.Round(cargPer*100) / 100,
	}
	return cagrRes

}
