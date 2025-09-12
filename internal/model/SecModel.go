package model

import (
	"encoding/xml"
	"fmt"
	"log"
)

type SecModel struct {
	XMLName xml.Name `xml:"document"`
	Data    Data     `xml:"data"`
}

type Data struct {
	Rows Rows `xml:"rows"`
}

type Rows struct {
	Items []Row `xml:"row"`
}

type Row struct {
	SecID              string `xml:"secid,attr"`
	ShortName          string `xml:"shortname,attr"`
	RegNumber          string `xml:"regnumber,attr"`
	Name               string `xml:"name,attr"`
	ISIN               string `xml:"isin,attr"`
	IsTraded           string `xml:"is_traded,attr"`
	EmitentID          string `xml:"emitent_id,attr"`
	EmitentTitle       string `xml:"emitent_title,attr"`
	EmitentINN         string `xml:"emitent_inn,attr"`
	EmitentOKPO        string `xml:"emitent_okpo,attr"`
	Type               string `xml:"type,attr"`
	Group              string `xml:"group,attr"`
	PrimaryBoardID     string `xml:"primary_boardid,attr"`
	MarketPriceBoardID string `xml:"marketprice_boardid,attr"`
}

func (r *Rows) Validate() error {
	if len(r.Items) == 0 {
		log.Println("Empty rows recieved end of process")
		return fmt.Errorf("End of requests")
	}
	return nil
}
