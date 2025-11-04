package main

import (
	"fmt"

	"github.com/marinaF3/app/calendar"
	"github.com/marinaF3/app/cmd"
	"github.com/marinaF3/app/logger"
	"github.com/marinaF3/app/storage"
)

func main() {

	if err := logger.Init("app.log"); err != nil {
		fmt.Printf("Ошибка инициализации логгера: %v\n", err)
		return
	}
	defer logger.Close()

	logger.Info("Приложение запущено")

	s := storage.NewJsonStorage("calendar.json")
	c := calendar.NewCalendar(s)

	fmt.Println("Загрузка данных")

	err := c.Load()
	if err != nil {
		fmt.Printf("Ошибка загрузки данных: %s\n", err)
		return

	}
	fmt.Printf("Данные успешно загружены\n")

	cli := cmd.NewCmd(c)
	cli.Run()

}
