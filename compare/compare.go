package compare

import (
	"errors"
	"os"

	"gopkg.in/gographics/imagick.v2/imagick"
)

func ImageComp(fromPath, toPath string) (ssim float64, psnr float64, err error) {
	imagick.Initialize()
	defer imagick.Terminate()

	from, err := loadImage(fromPath)
	if err != nil {
		return 0, 0, err
	}
	to, err := loadImage(toPath)
	if err != nil {
		return 0, 0, err
	}

	if from.GetImageWidth() != to.GetImageWidth() || from.GetImageHeight() != to.GetImageHeight() {
		return 0, 0, errors.New("image size not match")
	}

	return Ssim(from, to), Psnr(from, to), nil

}
func loadImage(path string) (*imagick.MagickWand, error) {
	img := imagick.NewMagickWand()
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	err = img.ReadImage(path)
	if err != nil {
		return nil, err
	}
	return img, nil
}
