package encode

import (
	"fmt"
	"errors"
	"image_codec/internal/models"
	"image_codec/internal/codec/heap"
	"image_codec/internal/codec/serialization"
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

func buildHaffmanCodes(input []byte) {
	bytesFreq := make(map[byte]int, 256)

	for _, v := range input {
		bytesFreq[v]++
	}

	var result heap.MinHeap
	var node1, node2 models.HeapElement

	for k, v := range bytesFreq {
		result = result.AddNewElement(models.HeapElement{
			Type: "leaf",
			Value: k,
			Freq: v,
		})
	}
	fmt.Println(bytesFreq)

	for len(result) > 1 {
		node1, result = result.GetMinElement()
		node2, result = result.GetMinElement()

		if node1.Freq < node2.Freq {
			result = result.UnionTwoElements(node1, node2)
		} else {
			result = result.UnionTwoElements(node2, node1)
		}
	}

	bytesCodes := make(map[byte]models.HaffmanCode, 256)
	DFSstack := []models.HaffmanTreeUnit{}
	// Проверить корректность указания длины кодов (см результат)
	DFSstack = append(DFSstack, models.HaffmanTreeUnit{
		TreeNode: result[0].LeftChild,
		Code: models.HaffmanCode{
			BitCode: 0,
			CodeLen: 1,
		},
	})
	DFSstack = append(DFSstack, models.HaffmanTreeUnit{
		TreeNode: result[0].RightChild,
		Code: models.HaffmanCode{
			BitCode: 1,
			CodeLen: 1,
		},
	})

	for len(DFSstack) > 0 {
		elem := DFSstack[len(DFSstack)-1]
		DFSstack = DFSstack[:len(DFSstack)-1]
		
		if elem.TreeNode.Type == "leaf" {
			bytesCodes[elem.TreeNode.Value] = elem.Code
			continue
		}

		DFSstack = append(DFSstack, models.HaffmanTreeUnit{
			TreeNode: elem.TreeNode.LeftChild,
			Code: models.HaffmanCode{
				BitCode: elem.Code.BitCode << 1,
				CodeLen: elem.Code.CodeLen + 1,
			},
		})
		DFSstack = append(DFSstack, models.HaffmanTreeUnit{
			TreeNode: elem.TreeNode.RightChild,
			Code: models.HaffmanCode{
				BitCode: (elem.Code.BitCode << 1) | 1,
				CodeLen: elem.Code.CodeLen + 1,
			},
		})
	}

	for k, v := range bytesCodes {
		fmt.Printf("Байт: %d, битовый код: %b, длина кода: %d \n", k, v.BitCode, v.CodeLen)
	}
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
	buildHaffmanCodes(serialized)

	return serialized, nil
}