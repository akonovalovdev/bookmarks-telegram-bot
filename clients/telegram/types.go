package telegram

// Каждый запрос будет получать некоторые объекты ответов, которые будут содержать требуемые абдэйты, кроме прочей информации 
// дополнительная(прочая) информация, которая будет добавляться на каждый запрос

type UpdatesResponse struct {
	Ok bool `json:"ok"` // булевое значение, есть ли апдэйты
	Result []Update `json:"result"` //поле в котором лежат апдэйты
}

type Update struct {
	ID int `json:"update_id"`  //update_id  из сервера будут приходить в формате json и стандартный парсер будет искать поле
	Message *IncomingMessage `json:"message"`
}

//отдельный тип для структуры message для входящих(Incoming) сообщений
type IncomingMessage struct {
	Text string `json:"text"`
	From From `json:"from"`
	Chat Chat `json:"chat"`
}

//отдельная часть структуры IncomingMessage
type From struct {
	Username string `json:"username"`
}

//отдельная часть структуры IncomingMessage
type Chat struct {
	ID int `json:"id"`
}
