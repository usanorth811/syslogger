package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/mcuadros/go-syslog.v2"
)

func main() {
	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)

	server := syslog.NewServer()
	server.SetFormat(syslog.RFC5424)
	server.SetHandler(handler)
	server.ListenUDP("0.0.0.0:514")
	server.ListenTCP("0.0.0.0:514")
	fmt.Println("Listening")
	server.Boot()

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			fmt.Println(logParts)
			f, err := os.OpenFile("text.log",
				os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Println(err)
			}
			defer f.Close()

			logger := log.New(f, "prefix", log.LstdFlags)
			logger.Println(logParts)
		}
	}(channel)

	server.Wait()
}