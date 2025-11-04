package cmd

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/c-bata/go-prompt"
	"github.com/google/shlex"
	"github.com/marinaF3/app/calendar"
	"github.com/marinaF3/app/events"
	"github.com/marinaF3/app/logger"
)

const (
	AddEvent     = "add-event"
	List         = "list"
	UpdateEvent  = "update-event"
	RemoveEvent  = "remove-event"
	AddRemind    = "add-remind"
	CancelRemind = "cancel-remind"
	Help         = "help"
	Log          = "log"
	Exit         = "exit"
)

type Cmd struct {
	calendar *calendar.Calendar
	// wg       sync.WaitGroup
	logMu sync.Mutex
	log   []string
}

func NewCmd(c *calendar.Calendar) *Cmd {
	return &Cmd{
		calendar: c,
	}
}

func (c *Cmd) logLine(s string) {
	c.logMu.Lock()
	defer c.logMu.Unlock()
	c.log = append(c.log, s)
}

func (c *Cmd) output(s string) {
	fmt.Print(s)
	c.logLine(s)
}

func (c *Cmd) outputLn(s string) {
	c.output(s + "\n")
}

func (c *Cmd) snapshotLog() []string {
	c.logMu.Lock()
	defer c.logMu.Unlock()
	cp := make([]string, len(c.log))
	copy(cp, c.log)
	return cp
}

func (c *Cmd) executor(input string) {
	input = strings.TrimSpace(input)
	if input == "" {
		return
	}
	c.logLine("> " + input + "\n")

	parts, err := shlex.Split(input)
	if err != nil {
		c.outputLn(err.Error())
		return
	}
	if len(parts) == 0 {
		return
	}
	cmd := strings.ToLower((parts[0])) // берем первую часть как команду

	switch cmd {
	case AddEvent:
		{
			logger.Info("Обработка команды add-event")
			if len(parts) < 4 {
				fmt.Println("Формат: add-event \"название события\" \"дата и время в формате 2025/01/31 00:01\" \"приоритет\"")
				logger.Error("Неверный формат команды add-event")
				return

			}

			title := parts[1]
			date := parts[2]
			priority := events.Priority(parts[3])

			e, err := c.calendar.AddEvent(title, date, priority)

			if err != nil {
				fmt.Println("Ошибка добавления:", err)
				logger.Error("Ошибка добавления" + err.Error())
			} else {
				fmt.Println("Событие: \"", e.Title, "\" добавлено")
				logger.Info("Событие: \"" + e.Title + "\" добавлено")
			}
		}

	case UpdateEvent:
		{
			logger.Info("Обработка команды update-event")
			if len(parts) < 4 {
				fmt.Println("Формат: update-event <ID> \"новое название\" \"дата и время в формате 2025-01-31 00:01\" \"приоритет\"")
				logger.Error("Неверный формат команды update-event")
				return

			}

			ID := parts[1]
			title := parts[2]
			date := parts[3]
			priority := events.Priority(parts[4])

			oldTitle, newTitle, err := c.calendar.EditEvent(ID, title, date, priority)

			if err != nil {
				fmt.Println("Ошибка обновления события:", err) // выводим ошибки
				logger.Error("Ошибка обновления события: " + err.Error())
			} else {
				fmt.Printf("Изменено событие: \"%s\" на \"%s\"", oldTitle, newTitle)
				logger.Info(fmt.Sprintf("Событие обновлено: ID=%s, OldTitle=%s, NewTitle=%s", ID, oldTitle, newTitle))
			}
		}
	case RemoveEvent:
		{
			logger.Info("Обработка команды remove")
			if len(parts) < 2 {
				//TODO добавить формат
				logger.Error("Неверный формат команды remove")
				fmt.Printf("Формат: %s\n", "remove-event ID")
				return
			}

			ID := parts[1]

			deleteEvent, err := c.calendar.DeleteEvent(ID)

			if err != nil {
				fmt.Printf("Ошибка при удалении \"%s\"\n", err)
				logger.Error("Ошибка удаления события: " + err.Error())
				return
			}

			fmt.Println("Удалено событие \"", deleteEvent.Title, "\"")
			logger.Info(fmt.Sprintf("Событие удалено: ID=%s, Title=%s", deleteEvent.ID, deleteEvent.Title))

		}

	case List:
		logger.Info("Обработка команды list")
		eventsList := c.calendar.GetEvents()
		if len(eventsList) == 0 {
			logger.Info("Календарь пуст")
			fmt.Println("Календарь пуст!")
			return
		}
		logger.Info(fmt.Sprintf("Выведено %d событий", len(eventsList)))
		for _, e := range eventsList {
			fmt.Printf("▶ %s\n", e.Title)
			fmt.Printf("  🆔 %s\n", e.ID) // Добавлена строка с ID
			fmt.Printf("  📅 %s  🏷️ %s\n",
				e.StartAt.Format("02 Jan 15:04"),
				e.Priority)
			if e.Reminder == nil {
				fmt.Printf("  ⏰ Напоминание:  %s\n", "-")
			} else {
				fmt.Printf("  ⏰ Напоминание: %s\n", e.Reminder.At.Format("02 Jan 15:04"))
			}

			fmt.Println()
		}

	case AddRemind:
		{
			logger.Info("Обработка команды add-remind")
			if len(parts) < 4 {
				fmt.Printf("Формат: %s\n", "add-remind <ID> \"сообщение\" \"дата и время\" \"длительность\"")
				logger.Error("Неверный формат команды add-remind")
				return
			}

			id := parts[1]
			message := strings.TrimSpace(parts[2])
			at := parts[3]
			duration := parts[4]

			fmt.Println("duration", duration)

			if err := c.calendar.SetEventReminder(id, message, at, duration); err != nil {
				fmt.Printf("Ошибка: %s\n", err)
				logger.Error("Ошибка добавления напоминания: " + err.Error())
				return
			}
			fmt.Printf("Напоминание добавлено: \"%s\"\n", message)
			logger.Info(fmt.Sprintf("Добавлено напоминание: ID=%s, Message=%s", id, message))
		}

	case CancelRemind:
		{
			logger.Info("Обработка команды cancel-remind")
			if len(parts) < 2 {
				fmt.Printf("Формат: %s\n", "cancel-remind <ID>")
				logger.Error("Неверный формат команды cancel-remind")
				return
			}

			id := parts[1]
			if err := c.calendar.CancelReminder(id); err != nil {
				fmt.Printf("Ошибка: %s\n", err)
				logger.Error("Ошибка отмены напоминания: " + err.Error())
				return
			}
			fmt.Printf("Напоминание отменено\n")
			logger.Info(fmt.Sprintf("Напоминание отменено: ID=%s", id))

		}
	case Log:
		logger.Info("Обработка команды log")
		lines := c.snapshotLog()
		if len(lines) == 0 {
			fmt.Printf("Лог пуст\n")
			logger.Info("Лог пуст")
			return
		}
		for _, line := range lines {
			fmt.Print(line)
			logger.Info(fmt.Sprintf("Выведено %d строк лога", len(lines)))
		}

	case Exit:
		logger.Info("Обработка команды exit")
		fmt.Printf("Сохранение данных...\n")

		for _, e := range c.calendar.GetEvents() {
			if e.Reminder != nil {
				e.Reminder.Stop()
			}
		}
		err := c.calendar.Save()

		if err != nil {
			fmt.Printf("Ошибка сохранения данных: %s\n", err)
			logger.Error("Ошибка сохранения данных: " + err.Error())
			return
		}
		// c.calendar.Close()
		fmt.Printf("Данные успешно сохранены.Приложение завершило работу")
		logger.Info("Данные успешно сохранены.Приложение завершило работу")
		os.Exit(0)

	default:
		fmt.Println("Неизвестная команда:")
		fmt.Println("Введите 'help' для списка команд")
	}
	fmt.Println(">>", input)
}

func (c *Cmd) completer(d prompt.Document) []prompt.Suggest {
	suggestions := []prompt.Suggest{
		{Text: AddEvent, Description: "Добавить событие"},
		{Text: List, Description: "Показать все события"},
		{Text: RemoveEvent, Description: "Удалить событие"},
		{Text: Help, Description: "Показать справку"},
		{Text: Exit, Description: "Выйти из программы"},
		{Text: AddRemind, Description: "Добавить напоминание"},
		{Text: CancelRemind, Description: "Отменить напоминание"},
		{Text: UpdateEvent, Description: "Обновить событие"},
		{Text: Log, Description: "Посмотреть логи"},
	}

	return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
}

func (c *Cmd) Run() {
	p := prompt.New(
		c.executor,
		c.completer,
		prompt.OptionPrefix("> "),
		prompt.OptionMaxSuggestion(3),
	)
	go func() {
		for msg := range c.calendar.Notification {
			fmt.Printf("%s\n", msg)
		}
	}()
	p.Run()
}
