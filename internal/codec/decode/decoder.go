package decode

import (
	"errors"
	"image_codec/internal/codec/serialization"
	"image_codec/internal/models"
)

func deltaDecode(input []models.DeltaEncodedElement) []models.Pixel {
	result := make([]models.Pixel, len(input))

	result[0].R = byte(input[0].R)
	result[0].G = byte(input[0].G)
	result[0].B = byte(input[0].B)

	for i := 1; i < len(input); i++ {
		result[i].R = result[i-1].R + byte(input[i].R)
		result[i].G = result[i-1].G + byte(input[i].G)
		result[i].B = result[i-1].B + byte(input[i].B)
	}

	return result
}

func rleDecode(input []models.RLEEncodedElement) []models.DeltaEncodedElement {
	var result []models.DeltaEncodedElement
	var element models.DeltaEncodedElement
	var toFill byte

	for _, v := range input {
		toFill = v.Count
		
		element.R = v.Value.R
		element.G = v.Value.G
		element.B = v.Value.B

		for toFill > 0 {
			result = append(result, element)
			toFill--
		}
	}

	return result
}

func haffmanDecode(input []byte, haffmanCodesTable map[byte]models.HaffmanCode) []byte {
	var result []byte
	var bitCodeValue, bitCodeLen uint32
	var usedBits byte

	// Переделываем список кодов для ускорения работы
	haffmanCodesTableOptimized := map[models.HaffmanCode]byte{}

	for k, v := range haffmanCodesTable {
		haffmanCodesTableOptimized[v] = k
	}

	bitCodeLen = 1

	for _, v := range input {
		
		// Собираем искомый код
		// В текущем байте еще есть неиспользованные биты
outer:
		for usedBits < 8 {
			// Выковыриваем нужные биты
			decodingByte := v << usedBits
			usedBits++
			decodingByte = decodingByte >> 7 // здесь может быть проблема, когда длина кода больше 8
			bitCodeValue |= uint32(decodingByte)
			
			// Ищем код в таблице
			if b, ok := haffmanCodesTableOptimized[models.HaffmanCode{
				BitCode: bitCodeValue,
				CodeLen: bitCodeLen,
			}]; ok {
				result = append(result, b)
				bitCodeLen = 1
				bitCodeValue = 0
				continue outer
			}
			// Код не найден
			bitCodeLen++
			bitCodeValue = bitCodeValue << 1
		}
		usedBits = 0
	}

	return result
}

func Decode(width, height uint16, input []byte, haffmanCodesTable map[byte]models.HaffmanCode) ([]byte, error) {
	// Декодирование кодов Хаффмана
	haffmanEncodedBytes := haffmanDecode(input, haffmanCodesTable)
	
	// Десериализация
	inputRLEEncodedPixels := serialization.Deserialize(haffmanEncodedBytes)

	// RLE-декодирование
	inputDeltaEncodedPixels := rleDecode(inputRLEEncodedPixels)

	// Проверка первого пикселя
	if inputDeltaEncodedPixels[0].R < 0 || inputDeltaEncodedPixels[0].R > 255 {
		return nil, errors.New("некорректный первый пиксель канала R")
	}
	if inputDeltaEncodedPixels[0].G < 0 || inputDeltaEncodedPixels[0].G > 255 {
		return nil, errors.New("некорректный первый пиксель канала G")
	}
	if inputDeltaEncodedPixels[0].B < 0 || inputDeltaEncodedPixels[0].B > 255 {
		return nil, errors.New("некорректный первый пиксель канала B")
	}

	// Дельта-декодирование
	inputRawPixels := deltaDecode(inputDeltaEncodedPixels)

	valuesTotal := len(inputRawPixels)*3
	offset := 0
	result := make([]byte, valuesTotal)

	for _, v := range inputRawPixels {
		result[offset] = v.R
		result[offset+1] = v.G
		result[offset+2] = v.B

		offset += 3
	}

	// Проверка кол-ва данных
	if count := int(width)*int(height)*3; count != len(result) {
		return nil, errors.New("некорректное кол-во элементов")
	}

	return result, nil
}