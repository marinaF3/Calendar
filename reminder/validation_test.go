package reminder

import (
	"strings"
	"testing"
)

func TestValidationMessage(t *testing.T) {
	_, err := validateMessage(" ")

	if err != nil {
		t.Error("ошибка: пустое сообщение")
	}

	_, err = validateMessage(".")
	if err != nil {
		t.Error("ошибка: слишком короткое название")
	}

	_, err = validateMessage(strings.Repeat("A", 101))
	if err != nil {
		t.Error("ошибка: слишком длинное название")
	}

}
