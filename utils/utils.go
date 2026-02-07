package utils

import (
	"os"
	"fmt"
	"errors"
	"encoding/binary"
	"image_codec/internal/models"
)

func SaveFile(fileName string, width, height int, data []byte) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer file.Close()

	// Пишем заголовок
	// Magic
	err = binary.Write(file, binary.LittleEndian, [4]byte{'R', 'M', 'Z', 0})
	if err != nil {
		return err
	}
	// Version
	err = binary.Write(file, binary.LittleEndian, byte(1))
	if err != nil {
		return err
	}
	// Width
	err = binary.Write(file, binary.LittleEndian, uint16(width))
	if err != nil {
		return err
	}
	// Height
	err = binary.Write(file, binary.LittleEndian, uint16(height))
	if err != nil {
		return err
	}
	// DataSize
	err = binary.Write(file, binary.LittleEndian, uint32(len(data)))
	if err != nil {
		return err
	}

	// Пишем данные
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	fmt.Println("Файл успешно сохранен!")
	return nil
}

func ReadFile(fileName string) (uint16, uint16, []byte, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return 0, 0, nil, err
	}
	defer file.Close()

	fileData, _ := os.Stat(fileName)

	data := make([]byte, fileData.Size())
	n, _ := file.Read(data)
	data = data[:n]

	magic := [4]byte{}
	magic[0] = data[0]
	magic[1] = data[1]
	magic[2] = data[2]
	magic[3] = data[3]

	version := data[models.VersionOffset]
	widthBytes := data[models.WidthOffset:models.HeightOffset]
	widthValue := uint16(widthBytes[0] & 0xFF) | uint16(widthBytes[1]) << 8

	heightBytes := data[models.HeightOffset:models.DataSizeOffset]
	heightValue := uint16(heightBytes[0] & 0xFF) | uint16(heightBytes[1]) << 8

	// Проверка сигнатуры
	if magic != [4]byte{'R', 'M', 'Z', 0} {
		fmt.Println("сигнатура файла некорректна")
		return 0, 0, nil, errors.New("сигнатура файла некорректна")
	} else {
		fmt.Println("Сигнатура файла корректна")
	}
	// Проверка версии
	if version != byte(1) {
		fmt.Println("версия файла некорректна")
		return 0, 0, nil, errors.New("версия файла некорректна")
	} else {
		fmt.Println("Версия файла корректна")
	}

	fmt.Println("Размеры изображения:", widthValue, heightValue)

	data = data[13:]

	return widthValue, heightValue, data, nil
}

func ParseFile(input []byte) (uint16, uint16, []byte, error) {
	magic := [4]byte{}
	magic[0] = input[0]
	magic[1] = input[1]
	magic[2] = input[2]
	magic[3] = input[3]

	version := input[models.VersionOffset]
	widthBytes := input[models.WidthOffset:models.HeightOffset]
	widthValue := uint16(widthBytes[0] & 0xFF) | uint16(widthBytes[1]) << 8

	heightBytes := input[models.HeightOffset:models.DataSizeOffset]
	heightValue := uint16(heightBytes[0] & 0xFF) | uint16(heightBytes[1]) << 8

	// Проверка сигнатуры
	if magic != [4]byte{'R', 'M', 'Z', 0} {
		fmt.Println("сигнатура файла некорректна")
		return 0, 0, nil, errors.New("сигнатура файла некорректна")
	}

	// Проверка версии
	if version != byte(1) {
		fmt.Println("версия файла некорректна")
		return 0, 0, nil, errors.New("версия файла некорректна")
	}

	data := input[13:]

	return widthValue, heightValue, data, nil
}