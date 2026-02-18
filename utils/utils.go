package utils

import (
	"os"
	"fmt"
	"errors"
	"encoding/binary"
	"image_codec/internal/models"
)

func SaveFile(fileName string, width, height int, haffmanCodeTable map[byte]models.HaffmanCode, data []byte) error {
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
	err = binary.Write(file, binary.LittleEndian, byte(2))
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
	// Коды Хаффмана
	// Длина таблицы
	err = binary.Write(file, binary.LittleEndian, uint16(len(haffmanCodeTable)))
	if err != nil {
		return err
	}
	// Сама таблица
	for k, v := range haffmanCodeTable {
		binary.Write(file, binary.LittleEndian, k) // 1 байт
		binary.Write(file, binary.LittleEndian, v.BitCode) // 4 байта
		binary.Write(file, binary.LittleEndian, v.CodeLen) // 4 байта
	}

	// Пишем данные
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	fmt.Println("Файл успешно сохранен!")
	return nil
}

func ReadFile(fileName string) (uint16, uint16, []byte, map[byte]models.HaffmanCode, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return 0, 0, nil, nil, err
	}
	defer file.Close()

	fileData, _ := os.Stat(fileName)

	data := make([]byte, fileData.Size())
	n, err := file.Read(data)
	if err != nil {
		return 0, 0, nil, nil, err
	}

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

	// Читаем таблицу кодов Хаффмана
	haffmanTableLen := uint16(data[13] & 0xFF) | uint16(data[14]) << 8
	haffmanTableByteLen := haffmanTableLen * models.HaffmanCodesTableBytesLen
	dataStartPosition := models.HaffmanCodesTableOffset + haffmanTableByteLen

	// Восстанавливаем таблицу кодов
	haffmanCodesTable := make(map[byte]models.HaffmanCode, 256)
	for i := 15; haffmanTableLen > 0; haffmanTableLen-- {
		bitCodeValue := uint32(data[i+1] & 0xFF) | uint32(data[i+2]) << 8 | uint32(data[i+3]) << 16 | uint32(data[i+4]) << 24
		codeLenValue := uint32(data[i+5] & 0xFF) | uint32(data[i+6]) << 8 | uint32(data[i+7]) << 16 | uint32(data[i+8]) << 24
		
		haffmanCodesTable[data[i]] = models.HaffmanCode{
			BitCode: bitCodeValue,
			CodeLen: codeLenValue,
		}
		i += 9
	}
	
	// fmt.Println("Таблица на чтении!")
	// for k, v := range haffmanCodesTable {
	// 	fmt.Printf("Байт: %d, битовый код: %0*b, длина кода: %d \n", k, v.CodeLen, v.BitCode, v.CodeLen)
	// }

	// Проверка сигнатуры
	if magic != [4]byte{'R', 'M', 'Z', 0} {
		fmt.Println("сигнатура файла некорректна")
		return 0, 0, nil, nil, errors.New("сигнатура файла некорректна")
	} else {
		fmt.Println("Сигнатура файла корректна")
	}
	// Проверка версии
	if version != byte(2) {
		fmt.Println("версия файла некорректна")
		return 0, 0, nil, nil, errors.New("версия файла некорректна")
	} else {
		fmt.Println("Версия файла корректна")
	}

	fmt.Println("Размеры изображения:", widthValue, heightValue)

	data = data[dataStartPosition:]

	return widthValue, heightValue, data, haffmanCodesTable, nil
}