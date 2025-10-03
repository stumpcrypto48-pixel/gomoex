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
	SecID              string `xml:"secid,attr" gorm:"column:sec_id"`
	ShortName          string `xml:"shortname,attr" gorm:"column:short_name"`
	RegNumber          string `xml:"regnumber,attr" gorm:"column:reg_number"`
	Name               string `xml:"name,attr" gorm:"column:name"`
	ISIN               string `xml:"isin,attr" gorm:"column:isin"`
	IsTraded           int    `xml:"is_traded,attr" gorm:"column:is_traded"`
	EmitentID          int    `xml:"emitent_id,attr" gorm:"column:emitent_id"`
	EmitentTitle       string `xml:"emitent_title,attr" gorm:"column:emitent_title"`
	EmitentINN         string `xml:"emitent_inn,attr" gorm:"column:emitent_inn"`
	EmitentOKPO        string `xml:"emitent_okpo,attr" gorm:"column:emitent_okpo"`
	Type               string `xml:"type,attr" gorm:"column:type"`
	Group              string `xml:"group,attr" gorm:"column:group_name"`
	PrimaryBoardID     string `xml:"primary_boardid,attr" gorm:"column:primary_board_id"`
	MarketPriceBoardID string `xml:"marketprice_boardid,attr" gorm:"column:market_price_board_id"`
}

func (Row) TableName() string {
	return "moex.secs"
}

func (r *Rows) Validate() error {
	if len(r.Items) == 0 {
		log.Println("Empty rows recieved end of process")
		return fmt.Errorf("End of requests")
	}
	return nil
}
