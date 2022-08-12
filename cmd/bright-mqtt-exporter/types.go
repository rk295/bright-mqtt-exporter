package main
import (
	"strconv"
)

type BrightFloat64 float64

func (b *BrightFloat64) UnmarshalJSON(data []byte) error {

        s := string(data)
        s = s[1 : len(s)-1] // Remove quotes
		
		i, err := strconv.ParseFloat(s, 64)
        if err != nil {
                return err
        }
        *b = BrightFloat64(i)
        return nil
}

type BrightElectricitysMsg struct {
	Electricitymeter Meter `json:"electricitymeter"`
}

type BrightGasMsg struct {
	Gasmeter Meter `json:"gasmeter"`
}

type Meter struct {
	Timestamp string `json:"timestamp"`
	Energy    Energy `json:"energy"`
	Power     BrightFloat64 `json:"power"`
	Mpan      string `json:"mpan,omitempty"`
	Mprn      string `json:"mprn,omitempty"`
	Supplier  string `json:"supplier,omitempty"`
	Price     Price  `json:"price"`
}

type Energy struct {
	Export string `json:"export"`
	Units  string `json:"units"`
	Import Import `json:"import"`
}

type Import struct {
	Cummulative BrightFloat64 `json:"cummulative"`
	Day         string `json:"day"`
	Week        string `json:"week"`
	Month       string `json:"month"`
}

type Price struct {
	Unitrate       BrightFloat64 `json:"unitrate"`
	Standingcharge BrightFloat64 `json:"standingcharge"`
}

