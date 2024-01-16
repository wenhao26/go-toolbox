package main

import (
	"fmt"

	"toolbox/go-image-magick/convert"
	"toolbox/go-image-magick/utils"
)

func main() {
	if utils.SystemName() != "windows" {
		panic("暂时只能在windows系统使用！")
	}

	filePath := "F:\\test_files\\magick\\"

	// 获取ImageMagick版本信息
	//result, err := convert.MagickVersion()

	// 识别图片信息
	//result, err := convert.ImageIdentify(filePath + "bird.jpg")

	// 转换格式
	//result, err := convert.ImageFormat(filePath + "bird.jpg", filePath + "bird_format.png")

	// 调整图片尺寸
	//result, err := convert.ImageResize("600x900", filePath + "bird.jpg", filePath + "magick\\bird_resize.jpg")

	// 调节图片压缩比
	result, err := convert.ImageQuality("75%", filePath + "bird.jpg", filePath + "magick\\bird_quality_75.jpg")

	if err != nil {
		panic(err)
	}
	fmt.Print(result)
}
