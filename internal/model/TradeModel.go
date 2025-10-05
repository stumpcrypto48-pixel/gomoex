package model

import (
	"encoding/xml"
)

type TradeModel struct {
	XMLName xml.Name  `xml:"document"`
	Data    TradeData `xml:"data"`
}

type TradeData struct {
	Rows TradeRows `xml:"rows"`
}

type TradeRows struct {
	Items []TradeRow `xml:"row"`
}

type TradeRow struct {
	OpenPrice  float64 `xml:"open,attr"`
	ClosePrice float64 `xml:"close,attr"`
	LowPrice   float64 `xml:"low,attr"`
	HighPrice  float64 `xml:"high,attr"`
	Value      float64 `xml:"value,attr"`
	Volume     uint64  `xml:"volume,attr"`
	StartDate  string  `xml:"begin,attr"`
	EndDate    string  `xml:"end,attr"`
}

// func (Row) TableName() string {
// 	return "moex.secs"
// }

// func (r *Rows) Validate() error {
// 	if len(r.Items) == 0 {
// 		log.Println("Empty rows recieved end of process")
// 		return fmt.Errorf("End of requests")
// 	}
// 	return nil
// }
