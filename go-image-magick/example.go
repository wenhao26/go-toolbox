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
	//result, err := convert.ImageResize("600x900", filePath + "bird.jpg", filePath + "bird_resize.jpg")

	// 调节图片压缩比
	//result, err := convert.ImageQuality("75%", filePath + "bird.jpg", filePath + "bird_quality_75.jpg")

	// 裁剪图片
	//result, err := convert.ImageCrop(300, 200, 750, 750, filePath+"bird.jpg", filePath+"bird_crop_750.jpg")

	// 九宫格图片
	//result, err := convert.ImageGrid9(filePath+"bird.jpg", filePath+"grids\\bird_grid_%d.png")

	// 将图片切成N份格子图片
	//result, err := convert.ImageGrids(2, filePath+"2mb-image.jpg", filePath+"grids\\2mb-image\\2mb-image_%d.png")

	// 调整图片大小和质量压缩
	result, err := convert.ImageResizeAndQuality("800x700","85%", filePath + "20mb-image.jpg", filePath + "20mb-image__resize_quality_85.jpg")

	if err != nil {
		panic(err)
	}
	fmt.Print(result)
}
