package encode

import (
	"errors"
	"fmt"
	"image_codec/internal/codec/serialization"
	"image_codec/internal/models"
)

func deltaEncode(input []models.Pixel) []models.DeltaEncodedElement {
	result := make([]models.DeltaEncodedElement, len(input))

	result[0].R = int16(input[0].R)
	result[0].G = int16(input[0].G)
	result[0].B = int16(input[0].B)

	for i := 1; i < len(input); i++ {
		result[i].R = int16(input[i].R) - int16(input[i-1].R)
		result[i].G = int16(input[i].G) - int16(input[i-1].G)
		result[i].B = int16(input[i].B) - int16(input[i-1].B)
	} 

	return result
}

func rleEncode(input []models.DeltaEncodedElement) []models.RLEEncodedElement {
	var result []models.RLEEncodedElement

	for _, v := range input {
		if len(result) == 0 {
			result = append(result, models.RLEEncodedElement{
				Count: 1,
				Value: v,
			})
			continue
		}
		if result[len(result)-1].Value == v && result[len(result)-1].Count < 255 {
			result[len(result)-1].Count++
			continue
		} else {
			result = append(result, models.RLEEncodedElement{
				Count: 1,
				Value: v,
			})
			continue
		}
	}

	return result
}

func buildHaffmanCodes(input map[byte]int) {
	// Нужно писать кучу и функции для работы с ней

	fmt.Println(input)
}

func Encode(width, height int, input []byte) ([]byte, error) {
	fmt.Println("Входящий массив:", len(input))
	
	// Проверки размеров
	if width > 1000 {
		return nil, errors.New("превышение ширины изображения")
	}
	if height > 1000 {
		return nil, errors.New("превышение высоты изображения")
	}
	if width <= 0 || height <= 0 {
		return nil, errors.New("недопустимые размеры изображения")
	}
	
	pixelsTotal := width * height

	// Проверка количества значений
	if len(input) != pixelsTotal * 3 {
		return nil, errors.New("некорректная длина входящего массива")
	}

	pixelPos := 0
	rawPixels := make([]models.Pixel, pixelsTotal)

	for i := range rawPixels {
		rawPixels[i].R = input[pixelPos]
		rawPixels[i].G = input[pixelPos+1]
		rawPixels[i].B = input[pixelPos+2]
		pixelPos += 3
	}

	// Дельта-кодирование
	deltaEncodedPixels := deltaEncode(rawPixels)

	// RLE кодирование
	rleEncodedPixels := rleEncode(deltaEncodedPixels)

	// Сериализация
	serialized := serialization.Serialize(rleEncodedPixels)
	fmt.Println("Сериализовано данных:", len(serialized))

	// Коды Хаффмана
	// Считаем частоту вхождения байтов
	byteFreq := make(map[byte]int, 256)
	for _, v := range serialized {
		byteFreq[v]++
	}

	buildHaffmanCodes(byteFreq)

	return serialized, nil
}