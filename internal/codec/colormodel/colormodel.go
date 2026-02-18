package colormodel

import "image_codec/internal/models"

func RGBToYCbCr(input []models.Pixel) ([]byte, []byte, []byte) {
	resultY := make([]byte, len(input))
	resultCb := make([]byte, len(input))
	resultCr := make([]byte, len(input))

	for i := 0; i < len(input); i++ {
		yChannelValue := models.RedCoef * float32(input[i].R) + models.GreenCoef * float32(input[i].G) + models.BlueCoef * float32(input[i].B)
		cbChannelValue := 0.564 * (float32(input[i].B) - yChannelValue) // 0.564 * (B - Y)
		crChannelValue := 0.713 * (float32(input[i].R) - yChannelValue) // 0.713 * (R - Y)

		resultY[i] = byte(yChannelValue)
		resultCb[i] = byte(cbChannelValue) + 128
		resultCr[i] = byte(crChannelValue) + 128
	}

	return resultY, resultCb, resultCr
}

func YCbCrToRGB(yChannel []int16, cbChannel []int16, crChannel []int16) []models.Pixel {
	result := make([]models.Pixel, len(yChannel))

	for i := 0; i < len(yChannel); i++ {
		result[i].R = 0
		result[i].G = 0
		result[i].B = 0
	}

	return result
}