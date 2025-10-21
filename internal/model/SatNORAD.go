package model

type NORAD struct {
	ObjectName         string  `json:"OBJECT_NAME"`
	ObjectID           string  `json:"OBJECT_ID"`
	Epoch              string  `json:"EPOCH"`
	MeanNotion         float64 `json:"MEAN_MOTION"`
	Eccentricity       float64 `json:"ECCENTRICITY"`
	Inclination        float64 `json:"INCLINATION"`
	Raofascnode        float64 `json:"RA_OF_ASC_NODE"`
	Atgofpercinter     float64 `json:"ARG_OF_PERICENTER"`
	Meananomaly        float64 `json:"MEAN_ANOMALY"`
	Ephemeristype      int64   `json:"EPHEMERIS_TYPE"`
	ClassifiactionType string  `json:"CLASSIFICATION_TYPE"`
	Noradcatid         int64   `json:"NORAD_CAT_ID"`
	ElementSetNo       int64   `json:"ELEMENT_SET_NO"`
	Revatepoche        int64   `json:"REV_AT_EPOCH"`
	Bstar              float64 `json:"BSTAR"`
	MeanMotionDot      float64 `json:"MEAN_MOTION_DOT"`
	MeanMotionDDot     float64 `json:"MEAN_MOTION_DDOT"`
}
