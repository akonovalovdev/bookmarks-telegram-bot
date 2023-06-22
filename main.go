package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	tgClient "github.com/akonovalovdev/server/clients/telegram"
	"github.com/akonovalovdev/server/consumer/event-consumer"
	"github.com/akonovalovdev/server/events/telegram"
	"github.com/akonovalovdev/server/storage/files"
)

// read_adviser_akonovalovdev_bot

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "files_storage"
	batchSize   = 100
)

func main() {
	//создаём телеграм клиент(клиент который общаяется с телеграмом):
	//(тип Client и его методы реализованы в файле telegram.go в папке client)
	// получает сообщения которые ему пишут и отправляет собственные

	//создаём объект реализующий интерфейсы Processor и Fetcher
	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
	)

	log.Print("service started")

	//запускаем консьюмер consumer.Start(fetcher, processor)
	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	go func() {
		sig := []os.Signal{syscall.SIGTERM, syscall.SIGINT, os.Interrupt}
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, sig...)
		defer signal.Reset(sig...)
		<-sigChan
		log.Printf("Program critical stop")
		consumer.Stop()
	}()

	consumer.Start()
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
