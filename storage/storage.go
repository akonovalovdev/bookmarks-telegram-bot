package storage

import (
	"fmt"
	"io"
	"crypto/sha1"

	"github.com/akonovalovdev/server/lib/e"
)

//интерфейс который момжет работать с любой файловой системой, но в нашем случае это файловая система хранения
type Storage interface {  //storage-место хранения
	Save(p *Page) error
	PickRandom(userName string) (*Page, error) // принимает имя пользователя, чтобы понимать чьи ссылки искать и возвращает страницу
	Remove(p *Page) error
	IsExists(p *Page) (bool, error) // существует ли та или иная страница(говорит либо да либо нет) и логическая ошибка если он не смог узнать будевой параметр
}

//основной тип данных с которым будут работать Storage(страница на которую ведёт ссылка, которую мы скинули боту)
type Page struct { //page - страница
	 Url string // сам адрес непосредственно
	 UserName string // пользователь, которому отдавать
//	 Created time.Time // - дополнительное поле если потребуется возвращать самые старые или самые новые ссылки
}

//пишем метод для типа Page чтобы генерировать имя файла с помощью ХЭШа
//Возвращает текстовое представление Хэша и ошибку
func (p Page) Hash() (string, error) {
	h := sha1.New()

	//так как оба параметра типа стринг используем метод из пакета io WriteString
	if _, err := io.WriteString(h,p.Url); err != nil{
		return "", e.Wrap("can't calculate hash", err)
	
	}
	if _, err := io.WriteString(h,p.UserName); err != nil{
		return "", e.Wrap("can't calculate hash", err)
	}

	// возвращаем всё что получилось, переводя байты в строку
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}