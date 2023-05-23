package files

import (
	"fmt"
	"encoding/gob"
	"os"
	"errors"
	"path/filepath"
	"math/rand"
	"time"

	"github.com/akonovalovdev/server/lib/e"
	"github.com/akonovalovdev/server/storage"
)

// тип который будет реализовывать интерфейс Storage
type Storage struct {
	basePath string //базовый путь, Хранит информацию о том в какой папке мы будем хранить данные
}

const defaultPerm = 0774 //Значение по умолчанию - perm(разрешение)

var ErrNoSavedPages = errors.New("no saved pages") // Специальная переменная ошибки, если файлов нет(пользователь пока ничего не сохранил)

//функция которая будет создавать объекты типа Storage
func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

//метод сохраняющий ссылку в файл
func (s Storage) Save(page *storage.Page) (err error) {
	defer func() {err = e.WrapIfErr("can't save page", err)}()

	//определяемся, куда будем сохранять наши данные, в какую директорию
	//этот путь будет состоять из нескольких частей
	// path.Join()(собирает путь из кусочков, но нам не подходит так как в качестве разделителя использует СЛЭШ), а в винде разделитель обратный СЛЭШ
	//пакет filepath решает эту проблему, для винды в функции .Join() используется обратный слэш
	fPath := filepath.Join(s.basePath, page.UserName) // все файлы каждого пользователя мы будем складывать одноимённую папку UserName
	//Для того чтобы посмотреть что сохранил себе пользователь, достаточно посмотреть содержимое его папки

	//создаём путь с помощью пакета os
	//функция MkdirAll создаст все директории, переданные в этот путь
	//необходимо лишь указать параметры доступа для созданной директории
	if err := os.MkdirAll(fPath, defaultPerm); err!=nil{
		return err
	}

	//необходимо определиться с названием файла(все файлы должны иметь уникалное имя,нельзя сохранить более одного с одинаковым названием в 1 папке)
	//мы будем смотреть на параметры текущеq страницы и в зависимости от их значений генерировать какой-то ХЭШ
	fName, err := fileName(page)
	if err != nil {
		return err
	}

	//Добавим к пути до файла собственно имя этого файла
	fPath = filepath.Join(fPath, fName)

	// Создаём непосредственно файл
	file,  err := os.Create(fPath)
	if err != nil {
		return err
	}

	//Закрываем файл в дефё, игнорируя ошибку
	defer func() { _ = file.Close() } ()

	//теперь нам осталось преобразовать страницу, т.е. привести к формату, чтобы можно было записать в файл 
	//и по нему можно было бы восстановить исходную структуру
	// авторы GO подготовили для нас формат который подходит для этих целей идеально: gob
	//с помощью функции NewEncoder, мы создаём сам инкодер и ему сразу же передаём в аргумент файл file в который будет записываться результат и
	//у самого инкодера вызываем метод Encode, передав ему нашу страницу page
	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil
}


//Метод принимает имя пользователя, чтобы понимать чьи ссылки искать и возвращает страницу
func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() {err = e.WrapIfErr("can't do request", err)}()

	path := filepath.Join(s.basePath, userName) //аналогично предыдущему методу получаем путь до директории с файлами
	
	//получаем список файлов
	files, err := os.ReadDir(path)
	if err != nil{
		return nil, err
	}

	//проверяем наличие файлов, если файлов нет, мы вернём заранее определённую ошибку на такие случаи
	if len(files) == 0 {
		return nil, ErrNoSavedPages
	}

	//???????????????????????????????????????? что его не устраивает ниже(скрин сделал)
	//ТЕПЕРЬ НАМ НУЖНО ПОЛУЧИТЬ СЛУЧАЙНО ЧИСЛО ОТ 0 ДО НОМЕРА ПОСЛЕДНЕГО ФАЙЛА
	rand.Seed(time.Now().UnixNano()) // (7 10 1; 7 10 1 при перезапуске всегда одна и таже последовательнось (ПСЕВДОСЛУЧАЙНОСТЬ)) ->
	//-> чтобы этого избежать мы используем не константу, а текущее время и всегда будет разный Seed

	//теперь получаем само число, указывая верхнюю границу, которая будет совпадать с числом файлов
	n := rand.Intn(len(files))

	//получаем случанйый файл с тем номером, которы только что сгенерировали
	file := files[n]

	//декодируем файл и возвращаем его содержимое
	//для этого нам нужно открыть файл(open) и вызвать для него декодер(decode) вынесем эти действия в отдельный локальный метод
	
	return s.decodePage(filepath.Join(path, file.Name())) // ???????????????????? откуда метод Name
}


//метод удаления ссылок
func (s Storage) Remove(p *storage.Page) error {
	//получаем имя файла
	fileName, err := fileName(p)
	if err != nil {
		return e.Wrap("not available fileName, can't remove", err)
	}
	//Собираем полный путь до файла
	path := filepath.Join(s.basePath, p.UserName, fileName)

	//удаляем файл
	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("can't remove file %s", path)

		return e.Wrap(msg, err)
	} 

	return nil
}

//метод возвращает логический параметр существует данная страница или нет(сохранял ли пользователь её ранее)
func (s Storage) IsExists(p *storage.Page) (bool, error) {
	//получаем имя файла
	fileName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("not available fileName, file not exist", err)
	}
	//Собираем полный путь до файла
	path := filepath.Join(s.basePath, p.UserName, fileName)

	//проверяем существование файла
	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist): //?????????????????????????????????????? откуда данный метод ErrNotExist
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if file %s exists", path)

		return false, e.Wrap(msg, err)
	} 
	return true, nil
}

//функция декодирующая случайно выбранный файл
func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	//открываем файл
	f, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("can,t decode page", err)
	}
	defer func() { _=f.Close()}() // Закрываем файл

	//создаём переменную в которую файл будет декодирован
	var p storage.Page
	
	//осуществляем декодирование и в методе Decode мы передаём ссылку на страницу
	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Wrap("can,t decode page", err)
	}

	return &p, nil
}



//создаём функцию для определения  имени файла/ Она будет получать на вход string или error
//Отдельная функция для одного единиственного метода используется для удобства поиска в коде места где возможно придётся изменить порядок хэширования
func fileName(p *storage.Page) (string,error) {
	return p.Hash()
}