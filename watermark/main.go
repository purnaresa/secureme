package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"
)

func readImage(fileName string) (image image.Image) {
	baseFile, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("failed to open: %s", err)
	}
	defer baseFile.Close()

	fileNameList := strings.Split(fileName, ".")

	switch fileNameList[1] {
	case "jpg":
		image, err = jpeg.Decode(baseFile)
	case "png":
		image, err = png.Decode(baseFile)
	default:
		err = fmt.Errorf("invalid file type : %s", fileNameList[1])
	}

	if err != nil {
		log.Fatalf("failed to decode for %s: %s", fileName, err)
	}

	return
}

func writeImage(image image.Image, fileName string) (err error) {
	fileOut, err := os.Create(fileName)
	if err != nil {
		return
	}
	defer fileOut.Close()

	err = jpeg.Encode(fileOut, image, &jpeg.Options{Quality: jpeg.DefaultQuality})
	return
}

func createWatermark(base, mark, output string) (err error) {
	baseImage := readImage(base) // step 1 ==> read base image
	markImage := readImage(mark) // step 2 ==> read mark image

	// step 3 ==> calculate position in center
	baseBound := baseImage.Bounds()
	markBound := markImage.Bounds()
	offset := image.Pt(
		(baseBound.Size().X/2)-(markBound.Size().X/2),
		(baseBound.Size().Y/2)-(markBound.Size().Y/2))
	//

	// step 4 ==> put watermark with 50% opacity
	outputImage := image.NewRGBA(baseBound)
	draw.Draw(outputImage, outputImage.Bounds(), baseImage, image.ZP, draw.Src)
	draw.DrawMask(outputImage, markImage.Bounds().Add(offset), markImage, image.ZP, image.NewUniform(color.Alpha{128}), image.ZP, draw.Over)
	//

	err = writeImage(outputImage, output) // step 5 ==> write output to file image
	if err != nil {
		log.Println(err)
	}
	return
}

func main() {
	baseImageFileName := flag.String("base", "-", "base image filename")
	markImageFileName := flag.String("mark", "-", "mark image filename")
	outImageFileName := flag.String("output", "-", "output image filename")
	flag.Parse()

	err := createWatermark(
		*baseImageFileName,
		*markImageFileName,
		*outImageFileName)
	if err != nil {
		log.Fatalln(err)
	}

}
