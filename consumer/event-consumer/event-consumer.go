package event_consumer

// в данном файле будет прописана реализация интерфейса consumer

import (
	"log"
	"time"

	"github.com/akonovalovdev/server/events"
)

// основной тип
type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int // Размер пачки - говорит нам о том сколько событий мы можем обрабатывать за раз
	done      chan struct{}
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) *Consumer {
	return &Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
		done:      make(chan struct{}),
	}
}

// реализация метода start
func (c *Consumer) Start() {
	updatesChan := make(chan events.Event, c.batchSize)
	for i := 0; i < 10; i++ {
		go c.handleEvents(updatesChan)
	}
	//здесь будет вечный цикл, который будет постоянно ждать новые события и обрабатывать их
	for {

		select {
		case <-c.done: //тут пишем условие если основная программа убита, тогда убиваем это горутину
			log.Printf("consumer finished")
			close(updatesChan)
			return
		default:
			// Ничего не делать
		}
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())

			continue
		}

		//проверяем сколько событий нам удалось получить и если оказалось что их 0 то мы так же пропускаем итерацию
		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		for _, e := range gotEvents {
			updatesChan <- e
		}
	}
}

func (c *Consumer) Stop() {
	close(c.done)

}

// Дополнительная функция для разгрузки метода Start
func (c *Consumer) handleEvents(evnts <-chan events.Event) {
	//перебираем events(события)
	for event := range evnts {
		//здесь полезным будет написать небольшое сообщение в log о том что мы получили новое событие и готовы его обработать
		log.Printf("got new events: %s", event.Text)

		//для обработки событий у нас уже есть процессор. Программа если что-то не так пошло с одним из событий, то она просто пропускает его обработку
		if err := c.processor.Process(event); err != nil {
			log.Printf("can't handle event: %s", err.Error())

			continue
		}

	}

}
