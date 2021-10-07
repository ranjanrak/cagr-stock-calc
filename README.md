# cagr-stock-calc
Cagr calculator for disciplined investor :smiley:. Calculates [cagr percentage](https://en.wikipedia.org/wiki/Compound_annual_growth_rate#Formula) based on the day's close price of the start and end date requested for any trading stock.

## Installation
```
go get github.com/ranjanrak/cagr-stock-calc
```
## Usage
```go
package main

import (
    cagrcalculator "github.com/ranjanrak/cagr-stock-calc"
)

func main() {
    // start and end date between which you wants to calculator cagr percentage
    startDate := time.Date(2020, time.November, 3, 0, 0, 0, 0, time.UTC)
    endDate := time.Date(2021, time.October, 5, 0, 0, 0, 0, time.UTC)
    // enter trading symbol of the stock
    cagrval := cagrcalculator.CagrCal("RELIANCE", startDate, endDate)
    fmt.Printf("%+v\n", cagrval)
}
```

## Response
```
{symbol:RELIANCE cagrper:44.97}
```

**Disclaimer** : This was done as a hobby project and shouldn't be use for commercial purpose. Go through NSE usage terms [here](https://www.nseindia.com/nse-terms-of-use).