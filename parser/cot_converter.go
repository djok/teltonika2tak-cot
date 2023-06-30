package parser

import (
	"encoding/xml"
	"fmt"
	"time"
)

type CoTConverter struct{}

type Event struct {
	XMLName xml.Name `xml:"event"`
	Type    string   `xml:"type,attr"`
	Time    string   `xml:"time,attr"`
	Start   string   `xml:"start,attr"`
	Stale   string   `xml:"stale,attr"`
	How     string   `xml:"how,attr"`
	Point   Point    `xml:"point"`
	Detail  Detail   `xml:"detail"`
}

type Point struct {
	XMLName xml.Name `xml:"point"`
	Lat     string   `xml:"lat,attr"`
	Lon     string   `xml:"lon,attr"`
	Hae     string   `xml:"hae,attr"`
	Ce      string   `xml:"ce,attr"`
	Le      string   `xml:"le,attr"`
}

type Detail struct {
	XMLName xml.Name `xml:"detail"`
	Contact Contact  `xml:"contact"`
}

type Contact struct {
	XMLName xml.Name `xml:"contact"`
	Phone   string   `xml:"phone,attr"`
}

type CotParsedData struct {
	IMEI    string
	Lat     float64
	Lon     float64
	Speed   float64
	Heading float64
}

func NewCotConverter() *CoTConverter {
	return &CoTConverter{}
}

func (c *CoTConverter) Convert(data *TeltParsedData) []byte {
	event := Event{
		Type:  "a-f-G-U-C",
		Time:  time.Now().Format(time.RFC3339),
		Start: time.Now().Format(time.RFC3339),
		Stale: time.Now().Add(time.Hour * 1).Format(time.RFC3339),
		How:   "m-g",
		Point: Point{
			Lat: fmt.Sprintf("%f", data.Lat),
			Lon: fmt.Sprintf("%f", data.Lon),
			Hae: "9999999.0",
			Ce:  "9999999.0",
			Le:  "9999999.0",
		},
		Detail: Detail{
			Contact: Contact{
				Phone: data.IMEI,
			},
		},
	}

	output, err := xml.MarshalIndent(event, "", "  ")
	if err != nil {
		panic(err)
	}

	return output
}
