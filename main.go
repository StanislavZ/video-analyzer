package main

import (
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/corona10/goimagehash"
)

var imagesToProcess int
var imagePath string
var hashAverageSlice []*goimagehash.ImageHash
var hashDifferenceSlice []*goimagehash.ImageHash
var imageDistances []int

func main() {

	start := time.Now()

	processArguments()

	parseImagesDirectory()

	generateImageHashes()

	calculateImageDistances()

	analyzeVideoFreeze()

	elapsed := time.Since(start)

	log.Printf("Time took %s", elapsed)
}

func processArguments() {
	//Check if ImagePath was specified as a first line argument
	if len(os.Args) > 1 {
		imagePath = os.Args[1]
	} else {
		fmt.Println("Path to video frames haven't been specified as first argument")
		os.Exit(1)
	}
}

func generateImageHashes() {
	for i := 1; i < imagesToProcess+1; i++ {
		imagePath := generateImagePath(i)
		file, fileErr := os.Open(imagePath)

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

		hashAverageSlice = append(hashAverageSlice, hashAverage)
		hashDifferenceSlice = append(hashDifferenceSlice, hashDifference)

		file.Close()
	}
}

func generateImagePath(index int) string {
	indexString := strconv.Itoa(index)

	var prefix string

	switch len(indexString) {
	case 1:
		prefix = "00000"
	case 2:
		prefix = "0000"
	case 3:
		prefix = "000"
	case 4:
		prefix = "00"
	case 5:
		prefix = "0"
	}

	return imagePath + "thumb" + prefix + indexString + ".jpg"
}

func calculateImageDistances() {
	for i := 0; i < len(hashAverageSlice)-1; i++ {
		hash1 := hashAverageSlice[i]
		hash2 := hashAverageSlice[i+1]

		distance, _ := hash1.Distance(hash2)

		imageDistances = append(imageDistances, distance)
	}

	for i := 0; i < len(hashDifferenceSlice)-1; i++ {
		hash1 := hashDifferenceSlice[i]
		hash2 := hashDifferenceSlice[i+1]

		distance, _ := hash1.Distance(hash2)

		imageDistances[i] += distance
	}

	fmt.Println("Image Distance length: ", len(imageDistances))
	fmt.Println("Image Distance array: ", imageDistances)
}

func analyzeVideoFreeze() {

	freezeHappened := false
	freezeFrameStarted := 0

	for i := 0; i < len(imageDistances); i++ {

		distance := imageDistances[i]

		if distance == 0 {
			freezeFrameStarted = i
			j := i
			//Check the length of the freeze
			for imageDistances[j] == 0 {
				//If index is at the end of slice then stop
				if j < len(imageDistances)-1 {
					j++
				} else {
					break
				}
			}
			freezedFramesAmount := j - i
			i = j

			//calculate distance on the edges of the freeze, if its non zero then no freeze
			//edgeFramesDistance := calculateTwoImagesDistance(freezeFrameStarted, i)

			freezeHappened = true
			log.Printf("Freeze started at frame %s and took %s frame(s):", strconv.Itoa(freezeFrameStarted), strconv.Itoa(freezedFramesAmount))
		}
	}

	if !freezeHappened {
		fmt.Println("No video freezing detected")
	}
}

func calculateTwoImagesDistance(index1, index2 int) int {

	distance1, _ := hashAverageSlice[index1].Distance(hashAverageSlice[index2+1])
	distance2, _ := hashDifferenceSlice[index1].Distance(hashDifferenceSlice[index2+1])

	return distance1 + distance2
}

func parseImagesDirectory() {
	files, err := ioutil.ReadDir(imagePath)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.Contains(file.Name(), ".jpg") {
			imagesToProcess++
		}
	}

	fmt.Println("imagesToProcess: ", imagesToProcess)
}
