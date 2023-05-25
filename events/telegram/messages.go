package telegram

	//сообщение вместо команды, будет содержать в себе краткую справку
	const msgHelp = `I can save and keep you pages. Also I can offer you them to read.
	
	In order to save the page, just send me al link to it.
	
	In order to get a random page from you list, send me command /rnd.
	Caution!` + // предупреждаем пользователя о том что ссылка после команды /rnd будет удалена
	//Всё что скидывает нам бо он также будет удалять для своего списка
	`After that, this page will be removed from your list!`


	const msgHello = "Hi there! 🤖\n\n" + msgHelp
	


	//прописываем группу коротких сообщений, которыми бот будет комментировать различные наши действия
	const (
		// неизвестная команда
		msgUnknownCommand = "Unknow command 🤔" 
		// когда пользователь запрашивает ссылку, но у бота уже не осталось сохраннёных ссылок либо не сохранял пока
		msgNoSavedPages = "You have no saved pages 🙊" 
		//это сообщение отправляет бот пользователю, когда тот скидывает новую ссылку и бот успешно её сохранит
		msgSaved = "Saved! 👌" 
		//сообщение на тот случа когда пользователь скидывает ссылку, которую он сохранял ранее и она всё ещё хранится в списке
		msgAlreadyExists = "You have already have this page in you list 🤗" 
	)  
