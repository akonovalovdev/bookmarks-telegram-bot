package telegram

import (
	// "github.com/akonovalovdev/server/lib/e"
	"github.com/akonovalovdev/server/clients/telegram"
)

//Один единственный тип данных, который будет реализовывать оба интерфейса Fetcher и Processor
type Processor struct {
	tg *telegram.Client
	offset int
	// storage
}

//функция которая создаёт тип процессор
func New(client *telegramm.Client, storage) {
	
}