package testing

// import (
// 	"os"
// 	"fmt"
// )

func GenerateGrayLinearGradient(width, height int) []byte {
	var result []byte

	for w := 0; w < width; w++ {
		for h := 0; h < height; h++ {
			grayValue := byte(h * 255 / (width - 1))
			toFill := 3
			for toFill > 0 {
				result = append(result, grayValue)
				toFill--
			}
			
		}
	}
	// saveAsPGM("testGrayLinearGradient.pgm", result, width, height)

	return result
}

// func saveAsPGM(filename string, pixels []byte, width, height int) {
//     file, _ := os.Create(filename)
//     defer file.Close()
    
// 	fmt.Fprintf(file, "P5\n%d %d\n255\n", width, height) // разобраться, что это означает
//     file.Write(pixels)
// }