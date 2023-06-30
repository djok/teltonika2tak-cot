package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"

	// "wintak-gw/gui"
	"teltonika2tak-cot/parser"
	"teltonika2tak-cot/server"
)

type Config struct {
	WintakServer struct {
		Address string `json:"address"`
	} `json:"wintak_server"`
}

func main() {
	// Load configuration
	var ProcessWinTak = true
	configFile, err := os.Open("config.json")
	if err != nil {
		ProcessWinTak = false
		log.Println(err)
	}
	defer configFile.Close()

	var config Config
	if ProcessWinTak {
		decoder := json.NewDecoder(configFile)
		err = decoder.Decode(&config)
		if err != nil {
			log.Println(err)
		}
	}

	// Initialize Teltonika server
	teltonikaServer := server.NewTeltonikaServer()

	// Initialize WINTAK server

	wintakServer := server.NewWintakServer(config.WintakServer.Address)

	// Initialize Teltonika parser
	teltonikaParser := parser.NewTeltonikaParser()

	// Initialize CoT converter
	cotConverter := parser.NewCotConverter()

	// Initialize GUI
	// gui := gui.NewGui()

	addr := ":7908"
	// Start listening for incoming connections
	go func() {
		teltonikaServer.ServeTCP(addr)
		log.Println("Running")
	}()

	// Start data conversion and transmission process
	go func() {
		for data := range teltonikaServer.DataChan {
			log.Println(".")
			parsedData := teltonikaParser.Parse(data)
			cotData := cotConverter.Convert(parsedData)
			if ProcessWinTak {
				wintakServer.Send(cotData)
			}
			// gui.Update(parsedData)
		}
	}()

	// Wait until Ctrl+C is pressed
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("Ctrl+C pressed, exiting...")

	// Start GUI
	// gui.Run()
}
