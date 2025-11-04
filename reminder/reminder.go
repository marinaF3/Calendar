package reminder

import (
	"errors"
	"fmt"
	"time"
)

type Reminder struct {
	Message string
	At      time.Time
	Sent    bool
	Timer   *time.Timer  `json:"-"`
	Notify  func(string) `json:"-"`
}

func NewReminder(message string, at time.Time, notify func(string)) (*Reminder, error) {
	now := time.Now().In(at.Location())
	if at.Before(now) {
		return nil, errors.New("время не может быть в прошлом")
	}
	msg, err := validateMessage(message)

	if err != nil {
		return nil, err

	}

	errValidateAt := validateAt(at)

	if errValidateAt != nil {
		return nil, errValidateAt
	}

	return &Reminder{
		Message: msg,
		At:      at,
		Sent:    false,
		Notify:  notify,
	}, nil
}

func (r *Reminder) Start() {
	r.Timer = time.AfterFunc(time.Until(r.At), r.Send)
}

func (r *Reminder) Send() {
	if r.Sent {
		return
	}
	fmt.Println("Напоминение:", r.Message)
	r.Sent = true
}
func (r *Reminder) Stop() {

	if r.Timer != nil {
		r.Timer.Stop()
		r.Timer = nil
	}

	// Здесь будет логика остановки напоминания
}

// import (
// 	"errors"
// 	"fmt"
// 	"time"
// )

// type Reminder struct {
// 	Message string
// 	At      time.Time
// 	Sent    bool
// 	timer   *time.Timer
// }

// func NewReminder(message string, at time.Time) (*Reminder, error) {
// 	now := time.Now().In(at.Location())
// 	if at.Before(now) {
// 		return &Reminder{}, errors.New("напоминание не может быть в прошлое")
// 	}

// 	return &Reminder{
// 		Message: message,
// 		At:      at,
// 		Sent:    false,
// 	}, nil
// }

// func (r *Reminder) Start() {
// 	fmt.Println("стартовал")
// 	delay := time.Until(r.At)
// 	fmt.Printf("Напоминание сработает через: %v\n", delay)

// 	r.timer = time.AfterFunc(delay, r.Send)

// }

// func (r *Reminder) Send() {
// 	if r.Sent {
// 		return
// 	}
// 	fmt.Println("Reminder!", r.Message)
// 	r.Sent = true
// }

// func (r *Reminder) Stop() error {
// 	isStop := r.timer.Stop()
// 	if !isStop {
// 		return errors.New("таймер уже истек или был ранее остановлен.")
// 	}
// 	return nil
// }
