package codec

import (
	"fmt"
	"image_codec/internal/codec/decode"
	"image_codec/internal/codec/encode"
	// "image_codec/testing"
	"image_codec/utils"
	"os"
)

func Run(command string, width, height int, inputFile, outputFile string) error {

	switch command {
	case "encode":
		// Генерируем тестовое изображение
		// testInput := testing.GenerateGrayLinearGradient(width, height)

		// Читаем указанный файл
		input, err := os.ReadFile(inputFile)

		// Кодирование
		bytes, haffmanCodes, err := encode.Encode(width, height, input)
		if err != nil {
			return err
		}

		// Сохраняем файл
		err = utils.SaveFile(outputFile, width, height, haffmanCodes, bytes)
		if err != nil {
			return err
		}
	case "decode":
		// Чтение файла
		w, h, fileData, err := utils.ReadFile(outputFile)
		if err != nil {
			return err
		}

		// Декодирование
		decodedInput, err := decode.Decode(w, h, fileData)
		if err != nil {
			return err
		}
		fmt.Println(decodedInput[0:10])
	}

	return nil
}