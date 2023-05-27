package telegram

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/akonovalovdev/server/lib/e"
)

// для осуществления запросов нам потребуются несколько параметров
type Client struct {
	host      string      // хост api сервиса телеграмма
	basePatch string      // базовый путь(префикс с которого начинаются все запросы) (tg-bot.com/bot<token>)
	client    http.Client // http client храним тут чтобы его не создавать для каждого  запроса отдельно
}

// для упрощения поиска переменных, переводим названия методов в константы (если ребята из телеги переименуют метод, мы с лёгкостью поймём где его искать)
const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

// функция которая будет создавать клиент
func New(host string, token string) *Client {
	return &Client{
		host:      host,
		basePatch: newBasePath(token), //реализация базового пути скрыта в отдельную функцию, так как мало интересует пользователя
		client:    http.Client{},      // стандартный
	}
}

// Так же если нам потребуется создавать токен в разных частях программы, отдельная функция нам поможет это сделать
func newBasePath(token string) string {
	return "bot" + token
}

// Метод получения апдэйтов(новых сообщений) Возвращает структуру в будет содержаться то что нам нужно знать об апдэйте
// Типы(структуры) для удобства работы с ними вынесены в отдельный файл types в этом же пакете
// Входные параметры: limit - количество апдэйтов, которок мы будем получать за один запрос;
//
//	offset- "смещение" Внутри api все апдэйты будут копиться ввиде очереди, с помощью offset api вернёт нам все
//	апдэйты начиная пачками от 1 до 100( без параметра не видны границы пачек получаемых из очереди)
func (c *Client) Updates(offset int, limit int) (updates []Update, err error) {
	defer func() { err = e.WrapIfErr("can't do request", err) }()
	//формируем параметры запроса
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset)) // метод add помогает добавить указанный параметр к запросу
	q.Add("limit", strconv.Itoa(limit))
	//Необходимо отправить запрос и так как код отправки запроса будет выглядеть одинаково для всех наших методов,
	//выносим его в отдельную функцию doRequest
	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}
	//сохраняем полученные значения в  переменную
	var res UpdatesResponse
	// распарсиваем (aнмаршиливаем json)
	//в первом аргументе что именно мы будем парсить, а второй куда, именно по ссылке!!(иначе функция анмаршалинг не сможет ничего туда добавить)
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, err
}

// метод для отправки сообщений
func (c *Client) SendMessage(chatId int, text string) error {
	//подготавливаем параметры запроса
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatId))
	q.Add("text", text)

	//выполняем запрос, тело ответа нам не понадобится + передаём параметры
	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

// Локальный метод функция используемая для отправки всех запросов
// на вход она получает некий метод в виде строки(полученный из документации) и сформированный запрос со всеми добавленными параметрами и возвращать слайс байт и возможно ошибку
func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	defer func() { err = e.WrapIfErr("can't do request", err) }()

	//сформируем url, на который будет отправляться запрос с помощью одноимённого пакета url
	u := url.URL{
		Scheme: "https",                        // протокол
		Host:   c.host,                         // хост из структуры клиент
		Path:   path.Join(c.basePatch, method), // путь состоящий из двух частей(базовая часть пути из клиента и метод полученный в аргументе)
	}
	//ФОРМИРУЕМ объект запроса(без отправки на сервер) только подготавливаем
	//передаём http метод GET, url в виде строки и тело запроса( в нашем случае пустое так как всё необходимое у нас уже есть виде параметров)
	//плюс в методе GET обычно отсутствует тело
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	//Передаём в свежеиспечённый объект реквест, параметры query, которые мы получили в аргументе
	//Приводим параметры query с помощью метода Encode к тому виду в котором допустимо отправлять на сервер
	req.URL.RawQuery = query.Encode()

	//Отправляем получившийся запрос! для отправки используем тот клиент который заранее подготовили
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	//закрываем тело ответа в дефё с игнорированием ошибки
	defer func() { _ = resp.Body.Close() }()

	//получаем содержимое ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
