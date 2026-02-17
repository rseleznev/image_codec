package encode

import (
	"errors"
	"fmt"
	"image_codec/internal/codec/heap"
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

func buildHaffmanCodes(input []byte) map[byte]models.HaffmanCode {
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
	DFSstack = append(DFSstack, models.HaffmanTreeUnit{
		TreeNode: result[0],
		Code: models.HaffmanCode{
			BitCode: 0,
			CodeLen: 0,
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

	// fmt.Println("Таблица на запись!")
	// for k, v := range bytesCodes {
	// 	fmt.Printf("Байт: %d, битовый код: %0*b, длина кода: %d \n", k, v.CodeLen, v.BitCode, v.CodeLen)
	// }

	return bytesCodes
}

func haffmanEncode(data []byte, codes map[byte]models.HaffmanCode) ([]byte, byte) {
	result := []byte{}
	var byteBuffer byte
	var bitsToFill byte = 8
	var filledBytes byte

	for i, v := range data {
		// fmt.Println()
		// fmt.Println("!НАЧАЛО ВНЕШНЕГО ЦИКЛА!")
		// fmt.Printf("Байт: %d, битовый код: %0*b, длина кода: %d \n", v, codes[v].CodeLen, codes[v].BitCode, codes[v].CodeLen)
		// fmt.Println("Доступно битов в начале внеш цикла: ", bitsToFill)
		// fmt.Printf("Буфер: %08b \n", byteBuffer)

		// Код влезает в буфер
		if codes[v].CodeLen <= uint32(bitsToFill) {
			// fmt.Println("Код влезает в буфер")
			byteBuffer = byteBuffer << codes[v].CodeLen | byte(codes[v].BitCode)
			bitsToFill -= byte(codes[v].CodeLen)
			// fmt.Printf("Буфер: %08b \n", byteBuffer)
			if bitsToFill == 0 {
				result = append(result, byteBuffer)
				bitsToFill = 8
				byteBuffer = 0
			}
			// Если дошли до последнего кодируемого байта
			if i == len(data)-1 {
				if bitsToFill > 0 && bitsToFill < 8 {
					byteBuffer <<= bitsToFill
					filledBytes = bitsToFill
					result = append(result, byteBuffer)
				}
				break
			}
			continue
		}

		// Код не влезает в буфер
		// Заполняем текущий буфер насколько можем
		byteBuffer <<= bitsToFill // освобождаем свободные младшие биты
		addingBits := byte(codes[v].BitCode >> (codes[v].CodeLen-uint32(bitsToFill)) )
		byteBuffer |= addingBits
		bitsLeft := codes[v].CodeLen-uint32(bitsToFill)
		// fmt.Printf("Буфер внеш цикла: %08b \n", byteBuffer)

		result = append(result, byteBuffer)

		byteBuffer = 0
		bitsToFill = 8
		
		// Докидываем оставшиеся биты
		for bitsLeft > 0 {
			// fmt.Println("!НАЧАЛО ВНУТРЕННЕГО ЦИКЛА!")
			// fmt.Println("Осталось битов в коде: ", bitsLeft)
			// Все оставшиеся биты влезают в буфер
			if bitsLeft <= uint32(bitsToFill) {
				// fmt.Println("Докидываем все оставшиеся биты")
				bitCode := byte(codes[v].BitCode << (8 - bitsLeft)) // отсекаем использованные биты слева
				addingBits = bitCode >> (8 - bitsLeft) // ставим в позицию младших битов
				// fmt.Printf("Добавляемые биты: %08b \n", addingBits)
				byteBuffer |= addingBits
				// fmt.Printf("Буфер внутр цикл: %08b \n", byteBuffer)
				bitsToFill -= byte(bitsLeft)
				bitsLeft = 0
				// fmt.Println("Доступно битов во внутр цикле: ", bitsToFill)
				// Буфер заполнен
				if bitsToFill == 0 {
					result = append(result, byteBuffer)
					bitsToFill = 8
					byteBuffer = 0
				}
				break
			}

			// fmt.Println("Докидываем биты, сколько влезет")
			// Оставшиеся биты не влезают в буфер
			bitCode := codes[v].BitCode << (codes[v].CodeLen - bitsLeft) // отсекаем использованные биты слева
			addingBits = byte(bitCode >> (codes[v].CodeLen - 8)) // отсекаем биты справа
			// fmt.Printf("Добавляемые биты: %08b \n", addingBits)
			byteBuffer |= addingBits
			// fmt.Printf("Буфер внутр цикл: %08b \n", byteBuffer)
			result = append(result, byteBuffer)
			bitsLeft -= 8
			bitsToFill = 8
			byteBuffer = 0
		}
		// Дошли до последнего кодируемого байта и остались незаполненные биты
		if i == len(data)-1 {
			if bitsToFill > 0 && bitsToFill < 8 {
				filledBytes = bitsToFill
				result = append(result, byteBuffer)
			}
		}
	}

	return result, filledBytes
}

func Encode(width, height int, input []byte) ([]byte, map[byte]models.HaffmanCode, error) {
	fmt.Println("Входящий массив:", len(input))
	
	// Проверки размеров
	if width > 1000 {
		return nil, nil, errors.New("превышение ширины изображения")
	}
	if height > 1000 {
		return nil, nil, errors.New("превышение высоты изображения")
	}
	if width <= 0 || height <= 0 {
		return nil, nil, errors.New("недопустимые размеры изображения")
	}
	
	pixelsTotal := width * height

	// Проверка количества значений
	if len(input) != pixelsTotal * 3 {
		return nil, nil, errors.New("некорректная длина входящего массива")
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
	haffmanCodes := buildHaffmanCodes(serialized)
	haffmanEncoded, _ := haffmanEncode(serialized, haffmanCodes)
	fmt.Println("Длина закодированного набора: ", len(haffmanEncoded))

	return haffmanEncoded, haffmanCodes, nil
}