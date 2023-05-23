package e

import "fmt"

//функция которая занимается оборачиванием ошибок
// на вход принимает текст сообщения с подсказкой и саму ошибку
//возвращает ошибку НЕ НУЛЕВУЮ
func Wrap(msg string, err error) error{
	return fmt.Errorf("%s: %w", msg,  err)
}

//функция не возвращающая нулевую ошибку
func WrapIfErr(msg string, err error) error{
	if err == nil {
		return nil
	}
	
	return Wrap(msg, err)
}