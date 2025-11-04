package events

import (
	"strings"
	"testing"
)

func TestIsValidTitle(t *testing.T) {
	if isValidTitle("Ab") {
		t.Error("ожидали false для слишком короткого заголовка")
	}

	long := strings.Repeat("A", 51)
	if isValidTitle(long) {
		t.Error("ожидали false для слишком длинного заголовка")
	}

	if isValidTitle("Hi!;*^$%") {
		t.Error("ожидали false для заголовка с запрещёнными символами")
	}

	if !isValidTitle("Заголовок 42") {
		t.Error("ожидали true для валидного заголовка")
	}
}
