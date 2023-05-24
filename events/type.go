package events
//абстрактный пакет который не привязан ни к одному мессенджеру

type Fetcher interface {
	Fetch(limit int) ([]Event, error)  // offset перенесли внутрь Fetcher (на вход получает только определенную пачку событий)
}

type Processor interface {
	Process(e Event) error
}

//тип события
type Type int 

//список событий кторые мы будем использовать
const (
	Unknown Type = iota //неизвестный тип, для случая когда мы не смогли определить что за тип у нашего события. Значение константы 0
	Message  // значение 1
)

type Event struct {
	Type Type //тип события
	Text string // текст события
	Meta interface{} //поле позволяющее гибко подгонять дополнительные получаемые поля для разных мессенджеров(у телеграм chatID и username)
}