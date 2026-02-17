package main

import (
	"os"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image_codec/internal/codec"
	"image_codec/internal/codec/decode"
	"image_codec/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "не указано ни одной команды")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "encode":
		set := flag.NewFlagSet(command, flag.ExitOnError)

		width := set.Int("width", 0, "ширина изображения")
		height := set.Int("height", 0, "высота изображения")
		input := set.String("inputFile", "default_input.raw", "название исходного файла с данными")
		output := set.String("outputFile", "default_output.rmz", "название сохраненного файла")

		set.Parse(os.Args[2:])

		if *width <= 0 || *width > 1000 {
			fmt.Fprintln(os.Stderr, "ширина изображения должна быть больше 0 и меньше или равна 1000")
			os.Exit(1)
		}
		if *height <= 0 || *height > 1000 {
			fmt.Fprintln(os.Stderr, "высота изображения должна быть больше 0 и меньше или равна 1000")
			os.Exit(1)
		}

		err := codec.Run(command, *width, *height, *input, *output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "decode":
		set := flag.NewFlagSet(command, flag.ExitOnError)

		sourceFile := set.String("sourceFile", "default_output.rmz", "название сохраненного файла")

		set.Parse(os.Args[2:])

		err := codec.Run(command, 0, 0, "", *sourceFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "view":
		set := flag.NewFlagSet(command, flag.ExitOnError)
		sourceFile := set.String("sourceFile", "default_output.rmz", "файл, который нужно открыть")

		set.Parse(os.Args[2:])

		ViewImage(*sourceFile)
		return

	case "info":
		fmt.Fprintln(os.Stderr, "encode - закодировать массив байт и записать в файл")
		fmt.Fprintln(os.Stderr, "decode - прочитать указанный файл и декодировать")

	default:
		fmt.Fprintln(os.Stderr, "неизвестная команда:", command)
		os.Exit(2)
	}
}

func ViewImage(fileName string) {
	// Чтение файла
	w, h, fileData, haffmanCodes, err := utils.ReadFile(fileName)
	if err != nil {
		fmt.Println("VIEW: ошибка чтения файла", err)
		return
	}

	// Декодирование
	decodedInput, err := decode.Decode(w, h, fileData, haffmanCodes)
	if err != nil {
		fmt.Println("VIEW: ошибка декодирования данных")
	}

	// Собираем картинку
	image := buildImage(int(w), int(h), decodedInput)

	// Запускаем просмотрщика
	myApp := app.New()
    window := myApp.NewWindow("MyCodec Viewer - " + fileName)

	imageCanvas := canvas.NewImageFromImage(image)
	imageCanvas.FillMode = canvas.ImageFillOriginal
    
    // Добавляем в окно
    window.SetContent(container.NewScroll(imageCanvas))
    
    // Запускаем
    window.Resize(fyne.NewSize(800, 600))
    window.ShowAndRun()
}

func buildImage(width, height int, input []byte) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	pixelPos := 0

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{
				R: input[pixelPos],
				G: input[pixelPos+1],
				B: input[pixelPos+2],
				A: 255,
			})
			pixelPos += 3
		}
	}

	return img
}