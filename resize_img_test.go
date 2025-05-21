package main

import (
	"image"
	"image/jpeg"
	"os"
	"testing"

	"github.com/nfnt/resize"
)

func TestX(t *testing.T) {
	// 1. 打开原始图片
	file, err := os.Open("./img/bar2.jpg")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 2. 解码图片
	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	// 3. 缩放到宽度 200px（高度自动计算，保持比例）
	thumbnail := resize.Resize(200, 0, img, resize.Bilinear)

	// 4. 保存缩略图
	out, err := os.Create("thumbnail.jpg")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	jpeg.Encode(out, thumbnail, &jpeg.Options{Quality: 80})
}
