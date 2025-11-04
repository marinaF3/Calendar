package storage

import (
	"archive/zip"
	"errors"
	"io"
	"os"
)

type ZipStorage struct {
	*Storage
}

func NewZipStorage(filename string) *ZipStorage { // создаем новый Json сторадж
	return &ZipStorage{
		&Storage{filename: filename}, // немного волшебства композиции
	}
}

func (z *ZipStorage) Save(data []byte) error { // принимаем массив байт (можно тот же самый JSON)
	f, err := os.Create(z.filename)
	if err != nil {
		return err
	}
	defer f.Close()

	zw := zip.NewWriter(f) // создадим архив в файле f
	defer zw.Close()

	w, err := zw.Create("data")
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func (z *ZipStorage) Load() ([]byte, error) {
	r, err := zip.OpenReader(z.filename) // создадим чтение из архива

	if err != nil {
		return nil, err
	}
	defer r.Close() // не забудем закрыть

	if len(r.File) == 0 {
		return nil, errors.New("архив пуст")
	} // проверка на пустоту

	file := r.File[0] // возьмем первый файл из архива (будем считать что он всегда один

	rc, err := file.Open() // откроем его
	if err != nil {
		return nil, err
	}
	defer rc.Close() // не забудем закрыть

	return io.ReadAll(rc) // прочитаем все данные и вернем их
}

// import (
// 	"archive/zip"
// 	"errors"
// 	"io"
// 	"os"
// )

// type ZipStorage struct {
// 	*Storage
// }

// func NewZipStorage(filename string) *ZipStorage {
// 	return &ZipStorage{
// 		&Storage{filename: filename},
// 	}
// }

// func (z *ZipStorage) Save(data []byte) error {
// 	f, err := os.Create(z.GetFilename())
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

// 	zw := zip.NewWriter(f)
// 	defer zw.Close()

// 	w, err := zw.Create("data")
// 	if err != nil {
// 		return err
// 	}
// 	_, err = w.Write(data)
// 	return err
// }

// func (z *ZipStorage) Load() ([]byte, error) {
// 	r, err := zip.OpenReader(z.GetFilename())
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer r.Close()

// 	if len(r.File) == 0 {
// 		return nil, errors.New("архив пуст")
// 	}
// 	file := r.File[0]
// 	rc, err := file.Open()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rc.Close()

// 	return io.ReadAll(rc)
// }
