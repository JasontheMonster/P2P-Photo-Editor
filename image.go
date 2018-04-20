package main

import(
	"image"
	"image/png"
	"os"
	"log"
)

func getImage(name string) *image.Image{
	img, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer img.Close()

	imgData, imgType, err := image.Decode(img)
	if err != nil {
		log.Fatal(err)
	}

	return imgData
}

func conv(img *image.Image, matrice [][]int) *image.NRGBA {
    imageRGBA := image.NewNRGBA((*img).Bounds())
    w := (*img).Bounds().Dx()
    h := (*img).Bounds().Dy()
    sumR := 0
    sumB := 0
    sumG := 0
    var r uint32
    var g uint32
    var b uint32
    for y := 0; y < h; y++ {
        for x := 0; x < w; x++ {

            for i := -1; i <= 1; i++ {
                for j := -1; j <= 1; j++ {

                    var imageX int
                    var imageY int

                    imageX = x + i
                    imageY = y + j

                    r, g, b, _ = (*img).At(imageX, imageY).RGBA()
                    sumG = (sumG + (int(g) * matrice[i+1][j+1]))
                    sumR = (sumR + (int(r) * matrice[i+1][j+1]))
                    sumB = (sumB + (int(b) * matrice[i+1][j+1]))
                }
            }

            imageRGBA.Set(x, y, color.NRGBA{
                uint8(min(sumR/9, 0xffff) >> 8),
                uint8(min(sumG/9, 0xffff) >> 8),
                uint8(min(sumB/9, 0xffff) >> 8),
                255,
            })

            sumR = 0
            sumB = 0
            sumG = 0

        }
    }

    return imageRGBA
}

func saveImg(img *image.Image){
	f,err := os.Create("newImg.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	png.Encode(f, img)
}

func show(img *image.Image) {
	pic.ShowImage(img)
}
