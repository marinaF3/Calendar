package calendar

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/marinaF3/app/dateFormat"
	"github.com/marinaF3/app/events"
	"github.com/marinaF3/app/storage"
)

type Calendar struct {
	calendarEvents map[string]*events.Event
	storage        storage.Store
	Notification   chan string
}

var (
	ErrEventNotFound    = errors.New("событие не найдено")
	ErrReminderNotFound = errors.New("у события нет напоминания")
)

func NewCalendar(s storage.Store) *Calendar {
	return &Calendar{
		calendarEvents: make(map[string]*events.Event),
		storage:        s,
	}
}

func (c *Calendar) AddEvent(title string, date string, priority events.Priority) (*events.Event, error) {
	e, err := events.NewEvent(title, date, priority)
	if err != nil {
		return &events.Event{}, err
	}
	c.calendarEvents[e.ID] = e
	return e, nil
}

func (c *Calendar) EditEvent(id string, title string, dateStr string, priority events.Priority) (string, string, error) {
	e, exists := c.calendarEvents[id]
	if !exists {
		return "", "", fmt.Errorf("id=%q: %w", id, ErrEventNotFound)
	}

	oldTitle := e.Title

	if title == "_" {
		title = e.Title
	}
	if dateStr == "_" {
		dateStr = dateFormat.FormatLocal(e.StartAt)
	}
	if priority == "_" {
		priority = e.Priority
	}

	err := e.Update(title, dateStr, priority)
	if err != nil {
		return "", "", err
	}
	return oldTitle, e.Title, nil
}

func (c *Calendar) DeleteEvent(ID string) (*events.Event, error) {

	e, exist := c.calendarEvents[ID]

	if !exist {
		return nil, ErrEventNotFound
	}
	if e.Reminder != nil {
		e.RemoveReminder()
	}

	delete(c.calendarEvents, ID)

	return e, nil

}

func (c *Calendar) ShowEvents() {
	for _, value := range c.calendarEvents {
		fmt.Println("Coбытие", value.Title, "-", value.StartAt.Format("2006-01-02 15:04:05"))
	}

}

func (c *Calendar) Save() error {
	data, err := json.Marshal(c.calendarEvents)
	if err != nil {
		return fmt.Errorf("ошибка при сериализации JSON: %w", err)
	}
	err = c.storage.Save(data)
	return err
}

func (c *Calendar) Load() error {
	data, err := c.storage.Load()
	if err != nil {
		return fmt.Errorf("ошибка загрузки из стораджа: %w", err)
	}
	err = json.Unmarshal(data, &c.calendarEvents)

	if err != nil {
		return fmt.Errorf("ошибка парсинга JSON: %w", err)
	}
	return nil
}

func (c *Calendar) GetEvents() map[string]*events.Event {
	events := make(map[string]*events.Event)
	for k, v := range c.calendarEvents {
		events[k] = v
	}
	return events
}

func (c *Calendar) SetEventReminder(ID string, message string, at string, duration string) error {
	e, exist := c.calendarEvents[ID]

	if !exist {
		return ErrEventNotFound
	}

	fmt.Println(at)
	return e.AddReminder(message, at, duration, c.Notify)
}

func (c *Calendar) CancelReminder(ID string) error {
	e, exist := c.calendarEvents[ID]
	if !exist {
		return ErrEventNotFound

	}

	if e.Reminder == nil {
		return ErrReminderNotFound

	}

	e.RemoveReminder()

	return nil
}

func (c *Calendar) Notify(msg string) {
	c.Notification <- msg
}

func (c *Calendar) Close() {
	close(c.Notification)
}
