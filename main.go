package main

import (
	"flag"
	"log"

	tgClient "github.com/akonovalovdev/server/clients/telegram"
	"github.com/akonovalovdev/server/events/telegram"
	"github.com/akonovalovdev/server/storage/files"
	"github.com/akonovalovdev/server/consumer/event-consumer"
)

const (
	tgBotHost = "api.telegram.org"
	storagePath = "files_storage"
	batchSize = 100
)

//токен 6143760943:AAEEJBrZPzkkSOh7ESj-RL6ms4ikFJF0cBI

func main() {
	//создаём телеграм клиент(клиент который общаяется с телеграмом): (тип Client и его методы реализованы в файле telegram.go в папке client)
	// получает сообщения которые ему пишут и отправляет собственные
	//tgClient = telegram.New(tgBotHost, mustToken())
	
	//создаём объект реализующий интерфейсы Processor и Fetcher
	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
	)

	log.Print("service started") //?????????????????????? для чего сообщения передаются через log (6й урокб 10я минута)

	//запускаем консьюмер consumer.Strart(fetcher, processor)
	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err) // записываем в лог сообщение об ошибке и останавливаем программу
	}
}	

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"token for access to telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}