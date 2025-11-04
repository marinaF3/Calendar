package reminder

import (
	"errors"
	"strings"
	"time"
)

const empty = ""
const maxMessageLen = 50

func validateMessage(message string) (string, error) {
	message = strings.TrimSpace(message)

	if message == empty {
		return empty, errors.New("сообщение пустое")
	}

	if len([]rune(message)) > maxMessageLen {
		return empty, errors.New("количество символов превышает 50")

	}

	return message, nil
}

func validateAt(at time.Time) error {
	if at.IsZero() {
		return errors.New("время не может быть нулевым")
	}

	if at.Before(time.Now()) {
		return errors.New("время не может быть в прошлом")
	}

	return nil
}
