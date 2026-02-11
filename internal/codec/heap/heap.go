package heap

import (
	"image_codec/internal/models"
)

type MinHeap []*models.HeapElement

// Добавление нового узла в кучу
func (heap MinHeap) AddNewElement(element models.HeapElement) MinHeap {
	heap = append(heap, &element)

	if len(heap) > 1 {
		heap.RecoverUp(len(heap)-1)
	}
	return heap
}

// Восстановление кучи вверх после добавления нового узла
func (heap MinHeap) RecoverUp(index int) {
	if index == 0 {
		return
	}
	
	newElementIndex := index
	parentElementIndex := (index-1)/2

	if heap[parentElementIndex].Freq > heap[newElementIndex].Freq {
		recoveredIndex := heap.SwapElements(parentElementIndex, newElementIndex)
		heap.RecoverUp(recoveredIndex)
	}
}

// Поменять 2 узла местами
func (heap MinHeap) SwapElements(parentIndex, childIndex int) int {
	parentElement := heap[parentIndex]
	childElement := heap[childIndex]

	heap[parentIndex] = childElement
	heap[childIndex] = parentElement

	return parentIndex // возврат старшего узла
}

// Проверка корректности кучи
func (heap MinHeap) IsValidHeap() bool {
	result := true
	// Временная заглушка
	// Будет проверка, есть ли у корня потомки
	if len(heap) < 3 {
		return result
	}
	
	// Проверяем двух первых потомков
	if heap[0].Freq > heap[1].Freq || heap[0].Freq > heap[2].Freq {
		result = false
		return result
	}

	checkingStack := []int{}
	checkingStack = append(checkingStack, 1, 2)

	for len(checkingStack) > 0 {
		checkingIndex := checkingStack[len(checkingStack)-1]
		checkingStack = checkingStack[:len(checkingStack)-1]

		leftChildIndex := 2*checkingIndex + 1
		rightChildIndex := 2*checkingIndex + 2

		// Если есть левый потомок
		if leftChildIndex <= len(heap)-1 {
			// Проверяем левого потомка
			if heap[checkingIndex].Freq > heap[leftChildIndex].Freq {
				result = false
				break
			}
			// Закидываем левого потомка в очередь на проверку
			checkingStack = append(checkingStack, leftChildIndex)
		}

		// Если есть правый потомок
		if rightChildIndex <= len(heap)-1 {
			// Проверяем правого потомка
			if heap[checkingIndex].Freq > heap[rightChildIndex].Freq {
				result = false
				break
			}
			// Закидываем правого потомка в очередь на проверку
			checkingStack = append(checkingStack, rightChildIndex)
		}
	}

	return result
}

// Получить минимальный узел (корень)
func (heap MinHeap) GetMinElement() (models.HeapElement, MinHeap) {
	minElement := *heap[0]
	heap[0] = heap[len(heap)-1]
	heap = heap[:len(heap)-1]

	// Восстанавливаем кучу
	heap.RecoverDown(0)

	return minElement, heap
}

// Восстановление кучи вниз после удаления корня
func (heap MinHeap) RecoverDown(index int) {
	// Если дошли до конца
	if index == len(heap)-1 {
		return
	}

	elementChilds := "none"
	leftChildIndex := 2*index + 1
	rightChildIndex := 2*index + 2

	// Определяем наличие потомков
	if leftChildIndex <= len(heap)-1 {
		elementChilds = "left"
	}
	if rightChildIndex <= len(heap)-1 {
		if elementChilds == "left" {
			elementChilds = "left + right"
		} else {
			elementChilds = "right"
		}
	}

	switch elementChilds {
	// Нет потомков
	case "none":
		return

	// Только левый потомок
	case "left":
		if heap[index].Freq > heap[leftChildIndex].Freq {
			heap.SwapElements(index, leftChildIndex)
			heap.RecoverDown(leftChildIndex)
		}
		
	// Только правый потомок
	case "right":
		if heap[index].Freq > heap[rightChildIndex].Freq {
			heap.SwapElements(index, rightChildIndex)
			heap.RecoverDown(rightChildIndex)
		}

	// Оба потомка
	case "left + right":
		// Если левый потомок меньше правого
		if heap[leftChildIndex].Freq < heap[rightChildIndex].Freq {
			// Если проверяемый узел больше левого потомка
			if heap[index].Freq > heap[leftChildIndex].Freq {
				heap.SwapElements(index, leftChildIndex)
				heap.RecoverDown(leftChildIndex)
			}
		} else { // Если правый меньше или равен левому потомку
			// Если проверяемый узел больше правого потомка
			if heap[index].Freq > heap[rightChildIndex].Freq {
				heap.SwapElements(index, rightChildIndex)
				heap.RecoverDown(rightChildIndex)
			}
		}

	}
}

// Соединить 2 элемента в узел
func (heap MinHeap) UnionTwoElement(a, b models.HeapElement) MinHeap {
	newElement := models.HeapElement{
		Freq: a.Freq + b.Freq,
		LeftChild: &a,
		RightChild: &b,
	}

	heap = heap.AddNewElement(newElement)
	return heap
}