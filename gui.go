package main 

import(
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"github.com/google/gxui"
	"github.com/google/gxui/drivers/gl"
	"github.com/google/gxui/samples/flags"
)

func (n Node) appMain(driver gxui.Driver) {
	args := flag.Args()
	if len(args) != 1 {
		fmt.Print("usage: image_viewer image-path\n")
		os.Exit(1)
	}

	fileName := args[0]
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Failed to open image '%s': %v\n", fileName, err)
		os.Exit(1)
	}

	source, _, err := image.Decode(f)
	if err != nil {
		fmt.Printf("Failed to read image '%s': %v\n", fileName, err)
		os.Exit(1)
	}

	theme := flags.CreateTheme(driver)
	img := theme.CreateImage()

	mx := source.Bounds().Max
	window := theme.CreateWindow(mx.X, mx.Y, "Photo Editor")
	window.SetScale(flags.DefaultScaleFactor)
	window.AddChild(img)

	// Copy the image to a RGBA format before handing to a gxui.Texture
	rgba := image.NewRGBA(source.Bounds())
	draw.Draw(rgba, source.Bounds(), source, image.ZP, draw.Src)
	texture := driver.CreateTexture(rgba, 1)
	img.SetTexture(texture)

	ticker := time.NewTicker(time.Millisecond * 500)
	go func() {
		for _ = range ticker.C{
			img.editImage()
			img.Redraw()
		}
	}

	window.OnClose(driver.Terminate)
}

flag.Parse()
gl.StartDriver(appMain)