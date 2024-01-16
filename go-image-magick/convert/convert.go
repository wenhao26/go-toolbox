package convert

import (
	"fmt"

	"toolbox/go-image-magick/utils"
)

// https://github.com/gographics/imagick

// 获取版本信息
func MagickVersion() (string, error) {
	params := []string{"convert", "-version"}
	return utils.Magick(params)
}

// 查看图片信息
func ImageIdentify(filename string) (string, error) {
	params := []string{"identify", filename}
	return utils.Magick(params)
}

// 转换格式
func ImageFormat(inputImage, outputImage string) (string, error) {
	// convert -density 300 -quality 100 input.pdf  output.png
	params := []string{"convert", inputImage, outputImage}
	return utils.Magick(params)
}

// 调整图片大小
func ImageResize(resize, inputImage, outputImage string) (string, error) {
	// convert -resize "500x300" -strip -quality 75% input.jpg output.jpg
	params := []string{"convert", "-resize", resize, inputImage, outputImage}
	return utils.Magick(params)
}

// 调节图片压缩比
func ImageQuality(quality, inputImage, outputImage string) (string, error) {
	params := []string{"convert", "-quality", quality, inputImage, outputImage}
	return utils.Magick(params)
}

// 调整图片大小和质量压缩值
func ImageResizeAndQuality(resize, quality, inputImage, outputImage string) (string, error) {
	// convert -resize "500x300" -strip -quality 75% input.jpg output.jpg
	params := []string{"convert", "-resize", resize, "-strip", "-quality", quality, inputImage, outputImage}
	return utils.Magick(params)
}

// 裁剪图片
func ImageCrop(x, y, cropWidth, cropHeight int, inputImage, outputImage string) (string, error) {
	// magick convert -crop 600x700+100+100 input.jpg  output.png
	cropArea := fmt.Sprintf("%dx%d+%d+%d", cropWidth, cropHeight, x, y)
	params := []string{"convert", "-crop", cropArea, inputImage, outputImage}
	return utils.Magick(params)
}

// 九宫格图
func ImageGrid9(inputImage, outputImage string) (string, error) {
	// magick.exe convert -crop 33.333%x33.333% +repage 2mb-image.jpg grids/image_%d.png
	params := []string{"convert", "-crop", "33.333%x33.333%", "+repage", inputImage, outputImage}
	return utils.Magick(params)
}

// 裁剪成不同等分的格子图片
func ImageGrids(number int, inputImage, outputImage string) (string, error) {
	// 计算宽高占比值
	scale := float64(1) / float64(number)
	whScale := utils.RemoveTrailingZeros(scale * 100)
	whScaleParam := whScale + "%x" + whScale + "%"

	params := []string{"convert", "-crop", whScaleParam, "+repage", inputImage, outputImage}
	return utils.Magick(params)
}

// 缩略图
func ImageThumb() {
	// TODO 取图片中心部分为缩略图
	// convert c:/1.jpg -thumbnail "100x100^" -quality 100 -gravity center -extent 100x100 c:/2.jpg
	// convert c:/1.jpg -thumbnail 200x200 -background white -gravity center -extent 200x200 c:/6.jpg
	// convert -crop 100x200+300+400 c:/1.jpg c:/3.jpg
}
