package telegram

import (
	"github.com/akonovalovdev/server/storage"
	"github.com/akonovalovdev/server/clients/telegram"
)

//Один единственный тип данных, который будет реализовывать оба интерфейса Fetcher и Processor
type Processor struct {
	tg *telegram.Client
	offset int
	storage storage.Storage // используем именно абстрактный интерфейс стторадж, а не конктретную его реализацию
}

//допонительные поля возвращаемые типом Update из telegram
type Meta struct {
	ChatID int
	Username string
}
 
//функция которая создаёт тип процессор
func New(client *telegramm.Client, storage storage.Storage) *Processor{
	//offset у нас дефолтный, поэтому его не указываем
	return &Processor{
	tg: client,
	storage: storage,
	}
}

//метод интерфейса фетчер, извлекать
func (p Processor) Fetch(limit int) ([]events.Event, error) {
	//сначала нужно получить все апдэйты(используя внутренний оффсет и лимит из аргумента)
	updates, err := p.tg.Updates(p.offset, limit) 
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}
	//возвращаем нулевой результат, если апдейтов мы не нашли
	if len(updates) == 0 {
		return nil, nil
	}

	//создаём память под результат
	res := make([]events.Event, 0, len(update)) //??????????????? make создаёт слайс?

	//теперь обходим все апдейты и преобразовать их в эвенты
	for _, u := range update {
		//для преобразования напишием локю функцию event
		res = append(res, event(u))
	}

	//теперь необходимо обновить значение поля offset, поскольку оно изначально дефолтное
	//чтобы при вызове метода fetch, мы получили следующую порцию событий
	//значение offset напрямую связано с ID апдэйта(берём послейдний апдэйт, смотрим его ID и добавляем к этому айди еденичку
	// тогда при следующем запросе мы получим только те апдэйты, у которых ID больше чем из последних уже полученных)
	p.offset = updates[len(updates) - 1].ID + 1

	return res, nil
}

	func event(upd telegram.Update) events.Event {
		//выносим тип события в отдельную переменную
		updType := fetchType(upd)
				
		res := events.Event{
			//так же создаём 2 функции(для получения типа(Type) и текста(Text) соответственно)
			Type: updType,
			Text: fetchText(upd),
		}
		//нельзя просто так взять и добавить ещё 2 поля(username, id) из структуры Message типа Update, поскольку тип Event 
		//является общим для всех возможных мессенджеров. далеко не факт что любому мессенджеру понадобятся эти 2 поля
		//помещаем эти переменные в заранее подготовленную структуру Meta у данного пакета telegram
		if updType == event.Message { //поскольку тип является Message, мы точно знаем что Message не нулевое
			res.Meta = Meta{
				ChatID: upd.Message.Chat.ID,
				Username: upd.Message.From.Username,
			}
		}
		return res
	}

	func fetchText(upd telegram.Update) string {
		//если полученное значение будет nil, то произойдёт паника, поэтому исключаем этот момент
		if upd.Message == nil {
			return ""
		}

		return upd.Message.Text
	}

	func fetchType(upd telegram.Update) events.Type {
		//если полученное значение будет nil, то тип нам не известен и произойдёт паника, поэтому исключаем этот момент
		if upd.Message == nil {
			return events.Unknown
		}
		return events.Message
	}
}
