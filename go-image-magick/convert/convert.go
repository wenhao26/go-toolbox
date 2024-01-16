package convert

import (
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
	// -resize "500x300" -strip -quality 75% input.jpg output.jpg
	params := []string{"convert", "-resize", resize, inputImage, outputImage}
	return utils.Magick(params)
}

// 调节图片压缩比
func ImageQuality(quality, inputImage, outputImage string) (string, error) {
	params := []string{"convert", "-quality", quality, inputImage, outputImage}
	return utils.Magick(params)
}

// 缩略图
func ImageThumb() {
	// TODO 取图片中心部分为缩略图
	// convert c:/1.jpg -thumbnail "100x100^" -quality 100 -gravity center -extent 100x100 c:/2.jpg
	// convert c:/1.jpg -thumbnail 200x200 -background white -gravity center -extent 200x200 c:/6.jpg
	// convert -crop 100x200+300+400 c:/1.jpg c:/3.jpg
}
