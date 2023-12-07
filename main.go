package main

import (
	"flag"
	"log"
	tgClient "read-adviser-bot/clients/telegram"
	eventConsumer "read-adviser-bot/consumer/event-consumer"
	"read-adviser-bot/events/telegram"
	"read-adviser-bot/storage/files"
)

const tgBotHost = "api.telegram.org"
const storagePath = "file_storage"
const batchSize = 100

// 6731230463:AAHB2QiO-q70FaNAP78UExVxnEncG-TNses
func main() {
	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
	)

	log.Print("service started")

	consumer := eventConsumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
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
