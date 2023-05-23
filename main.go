package main

import (
	"flag"
	"log"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {
	//создаём телеграм клиент(клиент который общаяется с телеграмом): (тип Client и его методы реализованы в файле telegram.go в папке client)
	// получает сообщения которые ему пишут и отправляет собственные
	tgClient = telegram.New(tgBotHost, mustToken())
	
	//fetcher = fetcher.New()

	//prcessor = prcessor.New()

	//consumer.Strart(fetcher, processor)
}

func mustToken() string {
	token := flag.String(
		"token-bot-token",
		"",
		"token for access to telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}