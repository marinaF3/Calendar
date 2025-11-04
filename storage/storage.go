package storage

type Store interface { // создает интерфейс с именем Store
	Save(data []byte) error // метод сохранения данных
	Load() ([]byte, error)  // метод загрузки данных
	GetFilename() string    // метод получения имени файла
}

type Storage struct { // тип Storage
	filename string // с полем filename
}

func (s *Storage) GetFilename() string { // метод получения имени
	return s.filename // вернет строку с именем
}

// type Store interface { // создает интерфейс с именем Store
// 	Save(data []byte) error // метод сохранения данных
// 	Load() ([]byte, error)  // метод загрузки данных
// 	GetFilename() string    // метод получения имени файла
// }

// type Storage struct { // тип Storage
// 	filename string // с полем filename
// }

// func (s *Storage) Save(data []byte) error {
// 	return os.WriteFile(s.filename, data, 0644)
// }

// func (s *Storage) Load() ([]byte, error) {
// 	data, err := os.ReadFile(s.filename)
// 	return data, err
// }

// func (s *Storage) GetFilename() string { // метод получения имени
// 	return s.filename // вернет строку с именем
// }
