package brightmqtt

// This file defines the Go structs for parsing an MQTT message received from
// the Glow ethernet dongle. This is for messages sent to /your/ MQTT broker
// it will not work for messages sent to Glow's own broker, they are in the
// more verbose format.
//
// Example electricity meter message:
//
// {
// 	"electricitymeter": {
// 		"timestamp": "2022-08-25T06:16:59Z",
// 		"energy": {
// 			"export": {
// 				"cumulative": 0,
// 				"units": "kWh"
// 			},
// 			"import": {
// 				"cumulative": 4896.645,
// 				"day": 0.003,
// 				"week": 0.035,
// 				"month": 0.257,
// 				"units": "kWh",
// 				"mpan": "1012400931394",
// 				"supplier": "British Gas",
// 				"price": {
// 					"unitrate": 0.2924,
// 					"standingcharge": 0.3792
// 				}
// 			}
// 		},
// 		"power": {
// 			"value": 0.481,
// 			"units": "kW"
// 		}
// 	}
// }
//
// Example gas meter message:

// {
//     "gasmeter": {
//         "timestamp": "2022-08-25T06:27:51Z",
//         "energy": {
//             "import": {
//                 "cumulative": 12491.78,
//                 "day": 0,
//                 "week": 14.334,
//                 "month": 109.153,
//                 "units": "kWh",
//                 "cumulativevol": 1107.678,
//                 "cumulativevolunits": "m3",
//                 "dayvol": 0,
//                 "weekvol": 14.334,
//                 "monthvol": 109.153,
//                 "dayweekmonthvolunits": "kWh",
//                 "mprn": "3342241002",
//                 "supplier": "---",
//                 "price": {
//                     "unitrate": 0.07344,
//                     "standingcharge": 0.2722
//                 }
//             }
//         }
//     }
// }
//
import (
	"time"
)

type ElectricitysMsg struct {
	Electricitymeter ElectricityMeter `json:"electricitymeter"`
}

type GasMsg struct {
	Gasmeter GasMeter `json:"gasmeter"`
}

type GasMeter struct {
	Energy    GasEnergy `json:"energy"`
	Timestamp time.Time `json:"timestamp"`
}

type GasEnergy struct {
	Import GasImport `json:"import"`
}

type GasImport struct {
	Cumulative           float64 `json:"cumulative"`
	Day                  float64 `json:"day"`
	Week                 float64 `json:"week"`
	Month                float64 `json:"month"`
	Units                string  `json:"units"`
	Cumulativevol        float64 `json:"cumulativevol"`
	Cumulativevolunits   string  `json:"cumulativevolunits"`
	Dayvol               float64 `json:"dayvol"`
	Weekvol              float64 `json:"weekvol"`
	Monthvol             float64 `json:"monthvol"`
	Dayweekmonthvolunits string  `json:"dayweekmonthvolunits"`
	Mprn                 string  `json:"mprn"`
	Supplier             string  `json:"supplier"`
	Price                Price   `json:"price"`
}

type ElectricityMeter struct {
	Energy    ElectricityEnergy `json:"energy"`
	Power     Power             `json:"power"`
	Timestamp time.Time         `json:"timestamp"`
}

type ElectricityEnergy struct {
	Export ElectricityExport `json:"export"`
	Import ElectricityImport `json:"import"`
}

type ElectricityImport struct {
	Cumulative float64 `json:"cumulative"`
	Day        float64 `json:"day"`
	Month      float64 `json:"month"`
	Mpan       string  `json:"mpan"`
	Price      Price   `json:"price"`
	Supplier   string  `json:"supplier"`
	Units      string  `json:"units"`
	Week       float64 `json:"week"`
}

type ElectricityExport struct {
	Cummulative float64 `json:"cummulative"`
	Units       string  `json:"units"`
}

type Power struct {
	Value float64 `json:"value"`
	Units string  `json:"units"`
}

type Price struct {
	StandingCharge float64 `json:"standingcharge"`
	Unitrate       float64 `json:"unitrate"`
}
