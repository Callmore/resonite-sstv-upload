package main

import (
	"flag"
	"fmt"
	"image"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

var (
	imagePath   = flag.String("i", "", "Requred, Path to image to upload. Supports png, jpeg and gif.")
	userID      = flag.String("u", "", "Requied, Resonite User ID to post under.")
	machineID   = flag.String("m", "", "Resonite Machine ID to use.")
	isGrayscale = flag.Bool("g", false, "Causes the uploaded image to be grayscale instead of in colour.")
)

const postTargetURL = "http://sstv.foohy.net/upload"

func main() {
	flag.Parse()

	if *userID == "" {
		fmt.Printf("No user ID set! Use -u to set one.")
		return
	}

	img, err := loadImage(*imagePath)
	if err != nil {
		panic(err)
	}

	// Read in input image
	imgString := getPixelsAsString(img)

	if len(imgString) != img.Bounds().Dx()*img.Bounds().Dy()*6 {
		panic("output image is incorrect, length wrong.")
	}

	textToPost := fmt.Sprintf("%d %d\n%s", img.Bounds().Dx(), img.Bounds().Dy(), imgString)

	_, offset := time.Now().Zone()

	uri, err := url.Parse(postTargetURL)
	if err != nil {
		panic(err)
	}

	query := url.Values{}
	if *isGrayscale {
		query.Set("algo", "bw")
	} else {
		query.Set("algo", "color")
	}
	query.Set("neos_id", *userID)
	if *machineID != "" {
		query.Set("neos_mid", *machineID)
	}
	query.Set("utc_off", strconv.Itoa(offset))

	uri.RawQuery = query.Encode()

	resp, err := http.Post(uri.String(), "text/plain", strings.NewReader(textToPost))
	if err != nil {
		panic(err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		panic(fmt.Errorf("bad status code: %s", resp.Status))
	}

	fmt.Println("Upload complete! Image can be viewed globally at http://sstv.foohy.net/")
}

func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)

	return img, err
}

func getPixelsAsString(img image.Image) string {
	outputSlice := []string{}

	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			outputSlice = append(outputSlice, fmt.Sprintf("%02X%02X%02X", r>>8, g>>8, b>>8))
		}
	}

	return strings.Join(outputSlice, "")
}
