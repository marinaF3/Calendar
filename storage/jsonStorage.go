package storage

import "os"

type JsonStorage struct {
	*Storage
}

func NewJsonStorage(filename string) *JsonStorage { // создаем новый Json сторадж
	return &JsonStorage{
		&Storage{filename: filename}, // немного волшебства композиции
	}
}

func (s *JsonStorage) Save(data []byte) error {
	err := os.WriteFile(s.GetFilename(), data, 0644) // берем имя файла через метод базового типа
	return err
}

func (s *JsonStorage) Load() ([]byte, error) {
	data, err := os.ReadFile(s.GetFilename()) // берем имя файла через метод базового типа
	return data, err
}

// import "os"

// type JsonStorage struct {
// 	*Storage
// }

// func NewJsonStorage(filename string) *JsonStorage {
// 	return &JsonStorage{
// 		&Storage{filename: filename},
// 	}
// }
// func (s *JsonStorage) Save(data []byte) error {
// 	err := os.WriteFile(s.GetFilename(), data, 0644)
// 	return err
// }

// func (s *JsonStorage) Load() ([]byte, error) {
// 	data, err := os.ReadFile(s.GetFilename())
// 	return data, err
// }

// package cmd

// import (
// 	"fmt"
// 	"os"
// 	"strings"

// 	"github.com/c-bata/go-prompt"
// 	"github.com/google/shlex"
// 	"github.com/marinaF3/app/calendar"
// 	"github.com/marinaF3/app/events"
// )

// func (c *Cmd) executor(input string) {
// 	if len(input) == 0 {
// 		return
// 	}
// 	parts, err := shlex.Split(input) // парсим ввод с помощью shlex.Split
// 	// обрабатываем ошибки
// 	if err != nil {
// 		return
// 	}

// 	cmd := strings.ToLower(parts[0]) // берем первую часть как команду

// 	switch cmd {
// 	case "add":
// 		if len(parts) < 4 {
// 			fmt.Println("Формат: add \"название события\" \"дата и время\" \"приоритет\"")
// 			return
// 		}

// 		title := parts[1]
// 		date := parts[2]
// 		priority := events.Priority(parts[3]) // оставляем как было

// 		e, err := c.calendar.AddEvent(title, date, priority)
// 		if err != nil {
// 			fmt.Println("Ошибка добавления:", err) // здесь выведется ошибка из createValidatedEvent
// 			fmt.Println("err", err)
// 			fmt.Println("date", date)
// 		} else {
// 			fmt.Println("Событие:", e.Title, "добавлено")
// 		}

// 	case "update":
// 		if len(parts) < 5 {
// 			fmt.Println("Формат: update \"названия события\" \"новое название\" \"новая дата и время\" \"новый приоритет\"")
// 			return
// 		}
// 		id := parts[1]
// 		newTitle := parts[2]
// 		newDate := parts[3]
// 		newPriority := events.Priority(parts[4]) // оставляем как было

// 		found := false
// 		for key := range c.calendar.CalendarEvents {
// 			if c.calendar.CalendarEvents[key].ID == id {
// 				err := c.calendar.CalendarEvents[key].UpdateEvent(newTitle, newDate, newPriority)
// 				if err != nil {
// 					fmt.Println("Ошибка обновления:", err)
// 				} else {
// 					fmt.Println("Событие обновлено")
// 				}
// 				found = true
// 				break
// 			}
// 		}
// 		if !found {
// 			fmt.Println("Задача не найдена")
// 			return
// 		}

// 	case "help":
// 		fmt.Println("Список всех команд:")
// 		fmt.Println("Update. Формат: update \"названия события\" \"новое название\" \"новая дата и время\" \"новый приоритет\"")
// 		fmt.Println("Add. Формат: add \"названия события\" \"дата и время\" \"приоритет\"")
// 		fmt.Println("Remove. Формат: remove \"названия события\"")
// 		fmt.Println("Help. Формат: help")
// 		fmt.Println("Exit. Формат: exit")

// 	case "exit": // если команда exit
// 		c.calendar.Save()
// 		os.Exit(0) // красивый выход и завершение процесса

// 	default:
// 		fmt.Println("Неизвестная команда:") // подсказки для неизвестных команд
// 		fmt.Println("Введите 'help' для списка команд")
// 	}
// }
// func (c *Cmd) completer(d prompt.Document) []prompt.Suggest {
// 	suggestions := []prompt.Suggest{
// 		{Text: "add", Description: "Добавить событие"},
// 		{Text: "list", Description: "Показать все события"},
// 		{Text: "remove", Description: "Удалить событие"},
// 		{Text: "help", Description: "Показать справку"},
// 		{Text: "exit", Description: "Выйти из программы"},
// 	}

// 	return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
// }

// func (c *Cmd) Run() {
// 	p := prompt.New(
// 		c.executor,
// 		c.completer,
// 		prompt.OptionPrefix("> "),
// 	)
// 	p.Run()
// }

// type Cmd struct { // создадим тип
// 	calendar *calendar.Calendar // который имеет календарь как поле
// }

// func NewCmd(c *calendar.Calendar) *Cmd { // конструктор для Cmd
// 	return &Cmd{
// 		calendar: c,
// 	}
// }
