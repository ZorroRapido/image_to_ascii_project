package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"io/ioutil"
	"os"
	_ "strings"

	// Side-effect import.
	_ "image/png"
	"golang.org/x/image/draw"
	"github.com/buger/goterm"
)

var(
	outputFile = flag.String("o", "", "Output file")
	noScale = flag.Bool("noscale", false, "Parameters for image resize")
)

func scale(img image.Image, w int, h int) image.Image {
	dstImg := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.NearestNeighbor.Scale(dstImg, dstImg.Bounds(), img, img.Bounds(), draw.Over, nil)
	return dstImg
}

func convertPixel(c color.Color) rune {
	g := color.GrayModel.Convert(c)
	r, _, _, _ := g.RGBA()
	gc := uint8(r)

	if gc <= 25 {
		return '@'
	} else if gc <= 50 {
		return '%'
	} else if gc <= 75 {
		return '@'
	} else if gc <= 100 {
		return '#'
	} else if gc <= 125 {
		return '+'
	} else if gc <= 156 {
		return '-'
	} else if gc <= 182 {
		return '='
	} else if gc <= 208 {
		return '.'
	} else if gc <= 234 {
		return ':'
	}

 	return ' '
}



func toAscii(img image.Image) [][]rune {
	res := make([][]rune, img.Bounds().Dy())
	for i := range res {
		res[i] = make([]rune, img.Bounds().Dx())
	}
	for i := range res {
		for j := range res[i] {
			res[i][j] = convertPixel(img.At(j, i))
		}
	}
	return res
}

func openImage(imgName string) (image.Image, error) {
	f, err := os.Open(imgName)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(f)
	return img, err
}

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("usage: asciimg <image>")
		os.Exit(0)
	}

	img, err := openImage(flag.Arg(0))
	if err != nil {
		fmt.Println("Error :", err.Error())
		os.Exit(1)
	}

	if *outputFile != "" {  // Если указано имя файла для вывода
		ascii := toAscii(img)

		var outputData string

		for i := range ascii {
			for j := range ascii[i] {
				outputData += string(ascii[i][j])
			}
			outputData += "\n"
		}

		f, err := os.Open(*outputFile)
		if err != nil {
			fmt.Println("Error :", err.Error())
			os.Exit(1)
		}

		fmt.Fprint(f, outputData)
		ioutil.WriteFile(*outputFile, []byte(outputData), 0)

		f.Close()
	} else if *noScale == false {
		img = scale(img, 50, 200)

		ascii := toAscii(img)

		for i := range ascii {
			for j := range ascii[i] {
				fmt.Printf("%c", ascii[i][j])
			}
			fmt.Println()
		}
	} else if *noScale == true {
		terminalHeight := goterm.Height()
		terminalWidth := goterm.Width()

		img = scale(img, terminalWidth, terminalHeight)
	}
}