package main

import (
	"errors"
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"os"
)

const COLS = 120
const ROWS = 35
var CHARACTERS = []string{ " ", ".", ":", "-", "=", "+", "*", "%", "░", "▒", "▓", "█" }
var AMOUNT_CHARACTERS = len(CHARACTERS)

func min(a, b int) int {
	if a > b { return b } else { return a }
}

func colorToCharacter(value int) string {
	return CHARACTERS[min(value * AMOUNT_CHARACTERS/ 255, AMOUNT_CHARACTERS - 1)]
}

func toAsciiCharacter(img gocv.Mat, pixelWidth, pixelHeight int) (buffer string) {
	for i := 0; i < ROWS; i++ {
		for j := 0; j < COLS; j++ {
			// Compute the mean color for the current pixel we are printing
			rec := image.Rectangle{
				Min: image.Point{ X: j * pixelWidth, Y: i * pixelHeight },
				Max: image.Point{ X: (j + 1) * pixelWidth, Y: (i + 1) * pixelHeight },
			}
			region := img.Region(rec)
			gray := region.Mean()
			buffer += colorToCharacter(int(gray.Val1))
		}
		buffer += "\n"
	}
	return buffer
}

func main() {
	webcam, _ := gocv.OpenVideoCapture(0)
	defer webcam.Close()
	img := gocv.NewMat()
	defer img.Close()

	// Read the camera once to pre-compute the pixel's width & height
	webcam.Read(&img)

	width := img.Cols()
	height := img.Rows()
	pixelWidth := width / COLS
	pixelHeight := height / ROWS

	if COLS > width || ROWS > height {
		fmt.Println(errors.New("invalid amount of row or columns"))
		os.Exit(1)
	}

	for {
		webcam.Read(&img)
		gocv.CvtColor(img, &img, gocv.ColorBGRToGray)
		gocv.Flip(img, &img, 1)

		fmt.Println(toAsciiCharacter(img, pixelWidth, pixelHeight))
	}
}
