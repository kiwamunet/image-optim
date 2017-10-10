package compare

import (
	"math"

	"gopkg.in/gographics/imagick.v2/imagick"
)

func Psnr(from, to *imagick.MagickWand) float64 {
	width := int(from.GetImageWidth())
	height := int(from.GetImageHeight())
	colorMse := [3]float64{}
	for x := 0; x < width; x++ {
		for y := 0; y < 50; y++ {
			f, _ := from.GetImagePixelColor(x, y)
			t, _ := to.GetImagePixelColor(x, y)
			colorMse[0] += math.Pow(f.GetRed()-t.GetRed(), 2.0)
			colorMse[1] += math.Pow(f.GetGreen()-t.GetGreen(), 2.0)
			colorMse[2] += math.Pow(f.GetBlue()-t.GetBlue(), 2.0)
			f.Destroy()
			t.Destroy()
		}
	}
	mse := 0.0
	for i := 0; i < 3; i++ {
		mse += colorMse[i] / float64(width*height)
	}
	mse /= 3
	if mse == 0.0 {
		return 0.0
	}
	psnr := -10 * math.Log10(mse)
	return psnr
}
