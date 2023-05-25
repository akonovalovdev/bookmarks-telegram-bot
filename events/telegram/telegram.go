package telegram

import (
	"errors"

	"github.com/akonovalovdev/server/storage"
	"github.com/akonovalovdev/server/clients/telegram"
	"github.com/akonovalovdev/server/lib/e"
	"github.com/akonovalovdev/server/events"
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

var (
	 ErrUnknownEventType = errors.New("unknown event type")
	 ErrUnknownMetaType = errors.New("unknown meta type")
	)

//функция которая создаёт тип процессор
func New(client *telegram.Client, storage storage.Storage) *Processor{
	//offset у нас дефолтный, поэтому его не указываем
	return &Processor{
	tg: client,
	storage: storage,
	}
}

/*_______________________________________________________________________________________
	Явное отличае между апдэйтами(Updates) и эвентами(Events) заключается в следующем:
	Updates - понятия телеграмма и они относятся только к нему(в другом месенджере термина апдэйт возможно даже не будет)
	Events - более общая сущность, в неё мы можем преобразовывать всё что получим от других мессенджеров, что бы они нам не предоставляли
_________________________________________________________________________________________*/

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
	res := make([]events.Event, 0, len(updates)) //??????????????? make создаёт слайс?

	//теперь обходим все апдейты и преобразовать их в эвенты
	for _, u := range updates {
		//для преобразования напишием локальную функцию event
		res = append(res, event(u))
	}

	//теперь необходимо обновить значение поля offset, поскольку оно изначально дефолтное
	//чтобы при вызове метода fetch, мы получили следующую порцию событий
	//значение offset напрямую связано с ID апдэйта(берём послейдний апдэйт, смотрим его ID и добавляем к этому айди еденичку
	// тогда при следующем запросе мы получим только те апдэйты, у которых ID больше чем из последних уже полученных)
	p.offset = updates[len(updates) - 1].ID + 1

	return res, nil
}

//Метод будет выполнять различные действия в зависимости от типа эвента
func (p Processor) Process(event events.Event) error {
	//будет всего 2 возможных кейса(если в будущем придётся работать с другими апдэйтами телеги, добавим другой кейс)
	switch event.Type {
	case events.Message: //работаем с сообщением
		return p.processMessage(event) //выносим всю логику работы с сообщением в отдельную функцию принимающую на вход Event
	default: //когда не знаем с чем работаем
		return e.Wrap("can't process message", ErrUnknownEventType)
	}
}

func (p Processor) processMessage(event events.Event) error {
	//для работы с этим методом необходимо получить meta
	meta, err := meta(event) //процесс получения meta выносим в отдельную функцию
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	//в зависисмости от типа сообщения выбираем определённое действие с ним
	//если пользователь скинул ссылку - сохраняем её
	//если пользователь отправил копанду RND - то мы должны найти ссылку из сохранённых и вернуть ему
	//если пользователь отправит нам команду help - мы должны ему отправить краткую справку по боту
	//все эти группы действий назовём - КОМАНДАМИ(comands) и весь код, который будет к нему относиться, вынесем в отдельны файл
	//вызываем функцию-роутер
	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return e.Wrap("can't process massage", err)
	}
	return nil
}

	
func meta(event events.Event) (Meta, error) {
	//делаем typeresepsion
	res, ok := event.Meta.(Meta)
	if !ok{
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}

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
	if updType == events.Message { //поскольку тип является Message, мы точно знаем что Message не нулевое
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
