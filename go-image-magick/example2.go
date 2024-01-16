package main

import (
	"fmt"

	"gopkg.in/gographics/imagick.v2/imagick"
)

func main() {
	imagick.Initialize()
	defer imagick.Terminate()

	var err error
	mw := imagick.NewMagickWand()
	err = mw.ReadImage("F:\\test_files\\magick\\bird.jpg")
	if err != nil {
		panic(err)
	}

	width := mw.GetImageWidth()
	height := mw.GetImageHeight()

	// 调整图片尺寸
	reWidth := width / 2
	reHeight := height / 2
	err = mw.ResizeImage(reWidth, reHeight, imagick.FILTER_LANCZOS, 1)
	if err != nil {
		panic(err)
	}

	// 调整图片质量
	err = mw.SetImageCompressionQuality(75)
	if err != nil {
		panic(err)
	}

	// 导出图片
	_ = mw.WriteImage("F:\\test_files\\magick\\bird_magick_1.jpg")
	fmt.Println("OK")
}
