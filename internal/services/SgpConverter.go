package services

import (
	"httpfromtcp/rootmod/internal/model"
	"time"

	"github.com/3rdof2/sgp4"
)

func OMMConverter(ommData []sgp4.OMM, satNames []model.SatNames) []ResponseSatData {

	LLA := make(chan ResponseSatData)

	go func() {
		defer close(LLA)
		for i, omm := range ommData {

			var cur = sgp4.OMMinit(&omm)
			eci, _ := sgp4.Propogate(cur, time.Now())

			var tmp = sgp4.ECItoLLA(eci, time.Now())
			LLA <- CreateResponseData(satNames[i].Name, tmp)

		}
	}()

	llaRes := make([]ResponseSatData, 0, len(ommData))

	for lla := range LLA {
		llaRes = append(llaRes, lla)
	}

	return llaRes
}

type ResponseSatData struct {
	Lat     float64 `json:"Lat"`
	Long    float64 `json:"Long"`
	Alt     float64 `json:"Alt"`
	SatName string  `json:"SatName"`
}

func CreateResponseData(satName string, llaData sgp4.LLA) ResponseSatData {

	return ResponseSatData{
		Lat:     llaData.Lat,
		Long:    llaData.Long,
		Alt:     llaData.Alt,
		SatName: satName,
	}
}
