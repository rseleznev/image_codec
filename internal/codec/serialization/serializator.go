package serialization

import "image_codec/internal/models"

func Serialize(input []models.RLEEncodedElement) []byte {
	result := make([]byte, len(input)*7)
	offset := 0

	for _, v := range input {
		result[offset] = v.Count
		
		// Дробим каждый int16 на 2 отдельных байта
		rLittle := byte(v.Value.R & 0xFF)
		rBig := byte(v.Value.R >> 8)

		gLittle := byte(v.Value.G & 0xFF)
		gBig := byte(v.Value.G >> 8)

		bLittle := byte(v.Value.B & 0xFF)
		bBig := byte(v.Value.B >> 8)

		result[offset+1] = rLittle
		result[offset+2] = rBig
		result[offset+3] = gLittle
		result[offset+4] = gBig
		result[offset+5] = bLittle
		result[offset+6] = bBig

		offset += 7
	}

	return result
}

func Deserialize(input []byte) []models.RLEEncodedElement {
	result := make([]models.RLEEncodedElement, 0, len(input)/7)

	// Здесь может быть ошибка, если данные записаны некорректно!
	// Сделать тут или ранее проверку количества элементов?
	for i := 0; i < len(input); i += 7 {
		var element models.RLEEncodedElement
		
		element.Count = input[i]

		rLittleUint16 := uint16(input[i+1])
		rBigUint16 := uint16(input[i+2])
		r := rLittleUint16 | (rBigUint16 << 8)

		gLittleUint16 := uint16(input[i+3])
		gBigUint16 := uint16(input[i+4])
		g := gLittleUint16 | (gBigUint16 << 8)

		bLittleUint16 := uint16(input[i+5])
		bBigUint16 := uint16(input[i+6])
		b := bLittleUint16 | (bBigUint16 << 8)

		element.Value.R = int16(r)
		element.Value.G = int16(g)
		element.Value.B = int16(b)

		result = append(result, element) // переделать на вставку по индексу?
	}

	return result
}