package main

import (
	"fmt"
	"image/jpeg"
	"os"
	"strconv"

	"github.com/corona10/goimagehash"
)

func main() {
	var imagePath string

	//Check if ImagePath was specified as a first line argument
	if len(os.Args) > 1 {
		imagePath = os.Args[1]
	} else {
		fmt.Println("Path to video frames haven't been specified as first argument")
		os.Exit(1)
	}

	for i := 1; i < 10; i++ {
		fileName := imagePath + "thumb" + "000" + strconv.Itoa(i) + ".jpg"
		file, fileErr := os.Open(fileName)

		if fileErr != nil {
			fmt.Println("Error: can't read file: ", fileErr)
			os.Exit(1)
		}

		img, imgErr := jpeg.Decode(file)

		if imgErr != nil {
			fmt.Println("Error: can't decode image: ", imgErr)
			os.Exit(1)
		}

		hashAverage, _ := goimagehash.AverageHash(img)
		hashDifference, _ := goimagehash.DifferenceHash(img)

		file.Close()
	}

	// file1Path := imagePath + "thumb0001.jpg"
	// file2Path := imagePath + "thumb0002.jpg"

	// file1, file1err := os.Open(file1Path)
	// file2, file2err := os.Open(file2Path)
	// defer file1.Close()
	// defer file2.Close()

	// if file1err != nil {
	// 	fmt.Println("Can't read file1: ", file1err)
	// 	os.Exit(1)
	// }

	// if file2err != nil {
	// 	fmt.Println("Can't read file2: ", file2err)
	// 	os.Exit(1)
	// }

	// img1, _ := jpeg.Decode(file1)
	// img2, _ := jpeg.Decode(file2)
	// hash1, _ := goimagehash.AverageHash(img1)
	// hash2, _ := goimagehash.AverageHash(img2)
	// distance, _ := hash1.Distance(hash2)
	// fmt.Printf("AverageHash Distance between images: %v\n", distance)

	// hash1, _ = goimagehash.DifferenceHash(img1)
	// hash2, _ = goimagehash.DifferenceHash(img2)
	// distance, _ = hash1.Distance(hash2)
	// fmt.Printf("DifferenceHash Distance between images: %v\n", distance)
}
