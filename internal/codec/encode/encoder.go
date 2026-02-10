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

type minHeap []models.HeapElement

func (mH minHeap) AddNewElement(element models.HeapElement) minHeap {
	mH = append(mH, element)

	if len(mH) > 1 {
		mH.RecoverUp(len(mH)-1)
	}
	return mH
}

func (mH minHeap) RecoverUp(index int) {
	if index == 0 {
		return
	}
	
	newElementIndex := index
	parentElementIndex := (index-1)/2

	if mH[parentElementIndex].Freq > mH[newElementIndex].Freq {
		recoveredIndex := mH.SwapElements(parentElementIndex, newElementIndex)
		mH.RecoverUp(recoveredIndex)
	}
}

func (mH minHeap) SwapElements(parentIndex, childIndex int) int {
	parentElement := mH[parentIndex]
	childElement := mH[childIndex]

	mH[parentIndex] = childElement
	mH[childIndex] = parentElement

	return parentIndex
}

func (mH minHeap) CheckIfHeapIsValid() bool {
	// Временная заглушка
	if len(mH) < 3 {
		return true
	}
	
	if mH[0].Freq > mH[1].Freq || mH[0].Freq > mH[2].Freq {
		return false
	}

	checkingStack := []int{}
	checkingStack = append(checkingStack, 1)
	result := true

	for len(checkingStack) > 1 {
		checkingIndex := checkingStack[len(checkingStack)]
		checkingStack = checkingStack[:0]

		if checkingIndex == len(mH)-1 {
			// Последний элемент кучи, нет потомков
			break
		}

		leftChildIndex := 2*checkingIndex + 1
		rightChildIndex := 2*checkingIndex + 2

		// Если есть левый потомок
		if leftChildIndex <= len(mH)-1 {
			// Проверяем левого потомка
			if mH[checkingIndex].Freq > mH[leftChildIndex].Freq {
				result = false
				break
			}
			// Закидываем левого потомка в очередь на проверку
			checkingStack = append(checkingStack, leftChildIndex)
		}

		// Если есть правый потомок
		if rightChildIndex <= len(mH)-1 {
			// Проверяем правого потомка
			if mH[checkingIndex].Freq > mH[rightChildIndex].Freq {
				result = false
				break
			}
			// Закидываем правого потомка в очередь на проверку
			checkingStack = append(checkingStack, rightChildIndex)
		}
	}

	return result
}

func buildHaffmanCodes(input []byte) {
	bytesFreq := make(map[byte]int, 256)

	for _, v := range input {
		bytesFreq[v]++
	}

	var result minHeap

	for k, v := range bytesFreq {
		result = result.AddNewElement(models.HeapElement{
			Value: k,
			Freq: v,
		})
	}

	check := result.CheckIfHeapIsValid()

	fmt.Println(result)
	fmt.Println(check)
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