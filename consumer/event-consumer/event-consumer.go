package event_consumer
// в данном файле будет прописана реализация интерфейса consumer

import (
	"github.com/akonovalovdev/server/events"
)

//основной тип
type Consumer struct {
	fetcher events.Fetcher
	processor events.Processor
	batchSize int //размер пачки - говорит нам о том сколько событий мы можем обрабатывать за раз
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
	fetcher: fetcher,
	processor: processor,
	batchSize: batchSize,
	}
}

//реализация метода start
func (c Consumer) Start() error {
	//здесь будет вечный цикл, который будет постоянно ждать новые события и обрабатывать их
	for{
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())
		
			continue
		}

		//проверяем сколько событий нам удалось получитьж  и если оказалось что их 0 то мы так же пропускаем итерацию
		if len(gotEvents) == 0{
			time.Sleep(1 * time.Second)

			continue
		}
		
	}
}