---
title: Screenshots with Go
author: David Hintelmann
date: 2023-08-20
description: "Blog post for taking screenshots of your desktop with go programming language. Developed on Windows 11 but works on my Linux Ubuntu as well."
tags: ["go","screenshots","programming","windows11","linux","guide"]
---

## Introduction

I was looking for different ways to automate capturing a screenshot of my desktop.
It would be possible to use python but using go may be faster since it is a compiled language.
Then I will need to figure out how to capture portions of my screen instead of the entire thing incase I want to avoid capturing any timestamp from the operating system's clock.

## Environment

**Language:** [Go](https://go.dev) version 1.21

**Operating System:** Windows 11 (developed) - Linux Ubuntu (tested)

## kbinani's repo for taking screenshots

In this exercise we will be importing [kbinani's](https://github.com/kbinani/screenshot) screenshot library to capture desktop screens to a `.png` image.

Brief overview of this libraries features:
> - Go library to capture desktop screen.
> - Multiple display supported.
> - Supported GOOS: windows, darwin, linux, freebsd, openbsd, and netbsd.
> - cgo free except for GOOS=darwin.

In this blog post I will be discussing the first two points above, since I am developing this program on Windows 11. I have tested capturing a few screenshots on Linux Ubuntu but I haven't modified the code to detect the OS and crop the image's timestamp out. This feature currently only works on Windows 11.

kbinani's library is relatively simple since there are only five functions:

* Capture
  * Capture returns screen capture of specified desktop region. x and y represent distance from the upper-left corner of primary display. Y-axis is downward direction. This means coordinates system is similar to Windows OS.
* CaptureDisplay
  * CaptureDisplay captures whole region of displayIndex'th display, starts at 0 for primary display.
* CaptureRect
  * CaptureRect captures specified region of desktop.
* GetDisplayBounds
  * GetDisplayBounds returns the bounds of displayIndex'th display. The main display is displayIndex = 0.
* NumActiveDisplays
  * NumActiveDisplays returns the number of active displays.

More information about this librarie's functions at its [pkg.go.dev](https://pkg.go.dev/github.com/kbinani/screenshot#pkg-functions) repo.

## Explanation of Code

Before diving into the code I will give a brief overview on the folder structure for this project. As this is a simple program, we also have a simple folder structure.

### Working Directory

Your working directory, called screenshot, should look like:
 - screenshot
   - main.go
   - go.mod
   - go.sum
 - img

img directory does not need to be created, this will be done automatically. The file `main.go` is where all our go code will be since this is a simple program. The other two files `go.mod` may need to be cleaned up using the `go mod tidy` command in the working directory to manage any dependencies, and `go.sum` file is used to validate the checksum of each direct and indirect dependency to confirm that none of them have been modified.

### Imports and Package Name

Below we name the package `package main` and after this import all the necessary packages for this code to run correctly.

```go
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
```

### Folder Name for Screenshots

As shown in the [Working Directory](#working-directory) section, our folder name for saving images is called `img` which is declared below.

```go
// folder name in this directory
// where screenshots will be saved
const folderName string = "img/"
```

### Main Function

The main function for each go program is where the program begins. First we use the `NumActiveDisplays()` function from kbinani's library to return the number of displays so we can loop through all of them and take a screenshot. Once we start looping over all the displays available we use the `GetDisplayBounds()` function to find the pixel size of the current display that we are looping over.

```go
func main() {
	n := screenshot.NumActiveDisplays()

	for i := 0; i < n; i++ {
		// get entire screen's bounds by pixel WxH
		bounds := screenshot.GetDisplayBounds(i)
        // ...
    }
}
```

### Check if img Folder Exists

Next in the main function, we check if the directory exists which we declared above as `folderName`. If it does not exist it will be created.

```go
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
```

### Create File

Now we will create a file in the `folderName` folder to save the `.png` data to. This file name will be based on the display number followed by the date and time. The format is:

    d_YYYY-MM-DD_XXHXXMXXS.png

    where:
        d - display number in int
        YYYY - year
        MM - month
        DD - day
        XXH - hours
        XXM - minutes
        XXS - seconds
    
    example: 0_2023-08-20_06H29M12S.png
        display - 0
        year - 2023
        MM - 08 (August)
        DD - 20
        XXH - 6 AM
        XXM - 29 minutes
        XXS - 12 seconds

go code

```go
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
```

### Take Screenshot

Now we will finally take a screenshot with the `screenShot()` that I created which takes two parameters. The first parameter is of type `bool` and is either `true` or `false`. If set to `false` then the screenshot will not include the timestamp of Windows 11 operating system in the bottom right corner. If the first parameter is set to `true` then the timestamp will be included in the screenshot, i.e. the entire screen will be captured. The second parameters passes the bounds of the screen to the function.

{{< alert >}}
**Warning!** if you are running this code on Linux it should work (worked on my Ubuntu set up) but remember to set the first parameter to true since the operating system's clock location is different for Linux and Windows. The screenshot function will have to modified if you plan on using this code on anything other than Windows 11.
{{< /alert >}}

After the screenshot is taken it will be saved as a `.png`

```go
        // take the screenshot
		img, err := screenShot(false, bounds)
		if err != nil {
			log.Fatalln("Could not take a screenshot.")
		}

		// save image
		png.Encode(file, img)
```

### screenShot Function

I created a function which is described above in the [Take Screenshot](#take-screenshot) section.

```go
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
```

### Install and run this code

In order to run this code you will have two options. It is also assumed you have go version 1.21 installed.

1. Create a new file called `main.go` and paste in the code from the [Full Code](#full-code) section below. 
   - Navigate to the directory, in terminal, where your `main.go` file is
   - Run the commands:
	> go mod init

	> go mod tidy
   - After you have initialized and tidied up the `go.mod` file you can run:
	> go run main.go
   - This will execute the program and take your screenshot!
2. Another method would be to copy `main.go` from my [GitHub](https://github.com/davidhintelmann/davidhintelmann.github.io/tree/main/content/posts/third)
   - Navigate to the directory, in terminal, where your `main.go` file is
   - Run the command:
	> go mod tidy
   - After you have tidied up the `go.mod` file you can run:
	> go run main.go
   - This will execute the program and take your screenshot!

You will find your screenshot in the `img` folder.

## Full Code

```go
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

```