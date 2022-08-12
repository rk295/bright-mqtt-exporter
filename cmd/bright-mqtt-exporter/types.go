package main

// This file defines the Go structs for parsing an MQTT message received from
// the Glow ethernet dongle. This is for messages sent to /your/ MQTT broker
// it will not work for messages sent to Glow's own broker, they are in the
// more verbose format.
//
// Example message:
//
// {
//     "gasmeter": {
//         "timestamp": "2022-08-12T09:50:59 +00",
//         "energy": {
//             "export": "0.00",
//             "units": "kWh",
//             "import": {
//                 "cummulative": "12458.51",
//                 "day": "5.42",
//                 "week": "29.76",
//                 "month": "75.89"
//             }
//         },
//         "power": "0.00",
//         "mprn": "111111111",
//         "price": {
//             "unitrate": "0.07",
//             "standingcharge": "0.27"
//         }
//     }
// }
//
// It makes an attempt to parse the timestamp and any float fields up into
// native Go types.
//

import (
	"strconv"
	"time"
)

const (
	brightTimeFormat = "2006-01-02T03:04:05 +00"
)

type BrightFloat64 float64

func (b *BrightFloat64) UnmarshalJSON(data []byte) error {

	i, err := strconv.ParseFloat(parseJSONString(data), 64)
	if err != nil {
		return err
	}
	*b = BrightFloat64(i)
	return nil
}

type BrightTime time.Time

func (t *BrightTime) UnmarshalJSON(data []byte) error {

	parsedTime, err := time.Parse(brightTimeFormat, parseJSONString(data))
	if err != nil {
		return err
	}
	*t = BrightTime(parsedTime)
	return nil
}

func parseJSONString(data []byte) string {
	s := string(data)
	return s[1 : len(s)-1] // Remove quotes
}

type BrightElectricitysMsg struct {
	Electricitymeter Meter `json:"electricitymeter"`
}

type BrightGasMsg struct {
	Gasmeter Meter `json:"gasmeter"`
}

type Meter struct {
	Timestamp BrightTime    `json:"timestamp"`
	Energy    Energy        `json:"energy"`
	Power     BrightFloat64 `json:"power"`
	Mpan      string        `json:"mpan,omitempty"`
	Mprn      string        `json:"mprn,omitempty"`
	Supplier  string        `json:"supplier,omitempty"`
	Price     Price         `json:"price"`
}

type Energy struct {
	Export BrightFloat64 `json:"export"`
	Units  string        `json:"units"`
	Import Import        `json:"import"`
}

type Import struct {
	Cummulative BrightFloat64 `json:"cummulative"`
	Day         BrightFloat64 `json:"day"`
	Week        BrightFloat64 `json:"week"`
	Month       BrightFloat64 `json:"month"`
}

type Price struct {
	Unitrate       BrightFloat64 `json:"unitrate"`
	Standingcharge BrightFloat64 `json:"standingcharge"`
}
