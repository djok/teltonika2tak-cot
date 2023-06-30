package parser

import (
	"encoding/binary"
	"log"
)

type TeltonikaParser struct {
}

type TeltParsedData struct {
	IMEI    string
	Time    uint64
	Lat     float64
	Lon     float64
	Alt     float64
	Heading float64
	Speed   float64
}

func NewTeltonikaParser() *TeltonikaParser {
	return &TeltonikaParser{}
}

func (p *TeltonikaParser) Parse(data []byte) *TeltParsedData {
	if len(data) < 33 {
		log.Println("Invalid data length")
		return nil
	}

	imei := string(data[:15])
	time := binary.BigEndian.Uint64(data[15:23])
	lat := float64(binary.BigEndian.Uint32(data[23:27])) / 10000000
	lon := float64(binary.BigEndian.Uint32(data[27:31])) / 10000000
	alt := float64(binary.BigEndian.Uint16(data[31:33])) / 100
	heading := float64(binary.BigEndian.Uint16(data[33:35])) / 100
	speed := float64(binary.BigEndian.Uint16(data[35:37])) / 100

	return &TeltParsedData{
		IMEI:    imei,
		Time:    time,
		Lat:     lat,
		Lon:     lon,
		Alt:     alt,
		Heading: heading,
		Speed:   speed,
	}
}
