package telegram

import (
	"errors"
	"log"
	"net/url"
	"strings"

	"github.com/akonovalovdev/server/lib/e"
	"github.com/akonovalovdev/server/storage"
)

// выносим ключевые слова по которым будем определять тип команды в отдельные константы
const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

// основной метод doCmd()/ что-то вроде API РОУТЕРА(мы будем смотреть на текст сообщения и по его формату содержания
// будем понимать какая это команда)
// каждую команду будет реализовывать отдельный метод типа процессор
func (p *Processor) doCmd(text string, chatID int, username string) error {
	//для начала мы обработаем текст сообщения, удалив из него лишние пробелы с помощью функции TrimSpace(так как они будут нам мешать)
	text = strings.TrimSpace(text)

	//создаём логи, чтобы изучать как пользуются ботом пользователи и пользуются ли вообще
	//сообщаем что получили новую команду, пишем её содержимое и сообщаем кто автор этого сообщения
	log.Printf("got new command '%s' from '%s", text, username)

	//СПИСОК КОМАНД: сохранить страницу(ссылка без пояснений); RND получить рандомную страницу (/rnd); инструкция бота (/help);
	//начало общения с ботом - автоматически(/start: hi + help)
	//3 команды, у которых есть ключевые слова мы будем определять по конструкции switch

	//для команды добавления страницы мы должны проверить, является ли сообщение ссылкой
	if isAddCmd(text) { //проверку на ссылку выносим в отдельную функцию
		//если проверка выполнилась, то мы имеем дело с командой TODO: AddPage()
		return p.savePage(chatID, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username) //команда воpвращения рандомной ссылки
	case HelpCmd:
		return p.sendHelp(chatID) //команда показать подсказку
	case StartCmd:
		return p.sendHello(chatID) //команда приветствия
	//дефолтный кейс, когда пользователь отправляет нам непонятно что(неизвестная команда или какой-то текст)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand) //пишем что не понимаем пользователя
	}
}

// метод сохраненя страницы; В методе переменную text переименовываем pageURL(чтобы не конфликтовать с пакетом с именем URL)
func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()

	//подготавливаем страницу которую собираемся сохранить
	page := &storage.Page{
		Url:      pageURL,
		UserName: username,
	}

	//Проверяем существует ли такая страница уже
	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}
	/*????????????????????????????????????????????????????????????????????????????????
	ВОТ ТУТ НИХУЯ НЕ ПОНЯТНО
	урок 5 16:25 Если окажется так, что команда уже существует/ мы отправим пользователю сообщение о том что ссылка уже сохранена
	Вопрос: разве в предыдущей проверке не было тоже самое??
	??????????????????????????????????????????????????????????????????????????????????*/
	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists) //выводим сообщение(а все сообщения перенесены в отдельный файл Messages.go в константы)
	}

	//сохраняем страницу
	if err := p.storage.Save(page); err != nil {
		return err
	}

	//если страница корректно сохранилась, мы сообщаем об этом пользователю
	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}
	return nil
}

// метод SendRandom, который будет отправлять пользователю случайную статью
func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send random", err) }()

	//ищем случайную статью
	page, err := p.storage.PickRandom(username)
	//обрабатываем обычную но не особую ошибки
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	//обрабатываем особую ошибку на тот случай если нет сохраннёных ссылок
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPage)
	}

	//если же нам удалось что-то найти, то мы отправляем эту ссылку пользователю
	if err := p.tg.SendMessage(chatID, page.Url); err != nil {
		return err
	}

	//последний шаг. если нам удалось найти и отправить ссылку, то нам обязательно нужно её удалить
	return p.storage.Remove(page)
}

// метод отправки справки ???????????????????????????????????????????????????????????
func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

// метод приветствия
func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

// проверяем является ли текст - ссылкой(так как способов проверки много, создаём отдельную функци. чтобы в случае чего изменить способ)
func isAddCmd(text string) bool {
	return isUrl(text)
}

func isUrl(text string) bool {
	// есть недостатки у этого способа (ссылки типа ya.ru не будут считаться ссылками, то есть те ссылки у которых не указан протокол)
	// чтобы ссылка считалась валидной, у неё всегда должен быть указан подобный префикс(http://...)
	// как будем проверять? Распаросим текст, считая его ссылкой
	u, err := url.Parse(text)

	return err == nil && u.Host != "" //текст считаем ссылкой в том случае если ошибка оказалась нулевая и при этом указан Host
}
