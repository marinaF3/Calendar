package events

import (
	"errors"
	"fmt"
	"strings"
	"time" // импорт встроенного пакета для даты

	"github.com/araddon/dateparse"
	"github.com/google/uuid"
	"github.com/marinaF3/app/dateFormat"
	"github.com/marinaF3/app/reminder"
)

type Event struct {
	Title    string             `json:"title"`
	StartAt  time.Time          `json:"start_at"`
	ID       string             `json:"id"`
	Priority Priority           `json:"priority"`
	Reminder *reminder.Reminder `json:"reminder"`
}

func NewEvent(title string, dateStr string, priority Priority) (*Event, error) {

	date, err := dateparse.ParseLocal(dateStr)
	if err != nil {
		return &Event{}, errors.New("неверный формат даты")
	}

	if date.Before(time.Now()) {
		return &Event{}, errors.New("дата не может быть в прошлом")
	}

	isTitleValid := isValidTitle(title)

	if !isTitleValid {
		return &Event{}, errors.New("неверное название")
	}

	isValidPriority := priority.Validate()

	if isValidPriority != nil {
		return &Event{}, errors.New("невалидный приоритет")
	}
	return &Event{
		ID:       uuid.New().String(),
		Title:    title,
		StartAt:  date,
		Priority: priority,
		Reminder: nil,
	}, nil

}

func (e *Event) Update(title string, date string, priority Priority) error {
	time, err := dateparse.ParseAny(date)

	if err != nil {
		return errors.New("невалидная дата")
	}

	isTitleValid := isValidTitle(title)

	if !isTitleValid {
		return errors.New("невалидное название")

	}

	isValidPriority := priority.Validate()

	if isValidPriority != nil {
		return errors.New("invalid priority")
	}

	e.Title = title
	e.StartAt = time
	e.Priority = Priority(priority)

	return nil

}

func (e *Event) Print() {
	fmt.Println(e)
}

func (e *Event) AddReminder(message string, at string, duration string, notify func(string)) error {
	at = strings.TrimSpace(at)
	if at == "" {
		return errors.New("время напоминания не может быть пустым")
	}

	var t time.Time

	if d, err := time.ParseDuration(duration); err == nil {
		if d <= 0 {
			return errors.New("время должно быть больше нуля")
		}
		t = time.Now().Add(d)
	} else {
		tt, err2 := dateFormat.ParseLocal(at)
		fmt.Println("at", at)
		if err2 != nil {
			return errors.New("неверный формат даты")
		}
		t = tt
	}

	r, err := reminder.NewReminder(message, t, notify)
	if err != nil {
		return err
	}
	e.Reminder = r
	e.Reminder.Start()
	return nil
}

func (e *Event) RemoveReminder() {
	if e.Reminder != nil {
		e.Reminder.Stop()
		e.Reminder = nil
	}
}
