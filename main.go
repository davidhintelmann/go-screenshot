package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"strings"
	"time"

	"github.com/kbinani/screenshot"
)

// folder name in this directory
// where screenshots will be saved
const folderName string = "img/"

func main() {
	n := screenshot.NumActiveDisplays()

	for i := 0; i < n; i++ {
		// get entire screen's bounds by pixel WxH
		bounds := screenshot.GetDisplayBounds(i)

		// check folderName folder
		// check if directory already exists
		// else create it
		folderInfo, err := os.Stat(folderName)
		if os.IsNotExist(err) {
			e := os.Mkdir(folderName, 0755)
			if e != nil {
				log.Fatal("Could not create img folder where screenshots are saved.")
			}
		}
		log.Println(folderInfo)

		// create file before saving png
		// convert file name to:
		// d_YYYY-MM-DD_XXHXXMXXS.png
		// where
		// d - display number in int
		// YYYY - year
		// MM - month
		// DD - day
		// XXH - hours
		// XXM - minutes
		// XXS - seconds
		now := time.Now()
		nowString := now.String()[:19]
		nowString = strings.Replace(nowString, " ", "_", 1)
		nowString = strings.Replace(nowString, ":", "H", 1)
		nowString = strings.Replace(nowString, ":", "M", 1)
		nowString += "S"
		filePath := folderName
		fileName := fmt.Sprintf(filePath+"%d_%s.png", i, nowString)
		file, _ := os.Create(fileName)
		defer file.Close()

		// take the screenshot
		img, err := screenShot(false, bounds)
		if err != nil {
			log.Fatalln("Could not take a screenshot.")
		}

		// save image
		png.Encode(file, img)
	}
}

func screenShot(timestamp bool, captureBox image.Rectangle) (*image.RGBA, error) {
	if timestamp {
		img, err := screenshot.CaptureRect(captureBox)
		if err != nil {
			panic(err)
		}
		return img, err
	} else {
		// create newBounds for screenshot
		// reduce number of y pixels to remove
		// operating system clock
		// subtract 50 from height on Windows 11
		newBounds := image.Rectangle{
			image.Point{0, 0}, // minimum pixels
			image.Point{captureBox.Dx(), captureBox.Dy() - 50}, // maximum pixels
		}
		img, err := screenshot.CaptureRect(newBounds)
		if err != nil {
			panic(err)
		}
		return img, err
	}
}
