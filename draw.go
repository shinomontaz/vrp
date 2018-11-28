package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	"github.com/shinomontaz/vrp/types"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func drawOrders(filename string, orders []*types.Order, wh *types.LatLng) {
	length := 500.0
	myImage := image.NewRGBA(image.Rect(0, 0, int(length), int(length)))
	outputFile, _ := os.Create(filename)
	defer outputFile.Close()
	myred := color.RGBA{200, 0, 0, 255}
	myblue := color.RGBA{0, 0, 200, 255}
	black := color.RGBA{0, 100, 0, 255}

	Wh := image.Rect(int(length*wh.Lat/100.0-2), int(length*wh.Lng/100.0-2), int(length*wh.Lat/100.0+2), int(length*wh.Lng/100.0+2))
	draw.Draw(myImage, Wh, &image.Uniform{myblue}, image.ZP, draw.Src)

	AxisLat := image.Rect(int(length*wh.Lat/100.0), int(length*wh.Lng/100.0), 0, int(length*wh.Lng/100.0+1))
	AxisLng := image.Rect(int(length*wh.Lat/100.0), int(length*wh.Lng/100.0), int(length*wh.Lat/100.0+1), 0)

	draw.Draw(myImage, AxisLat, &image.Uniform{black}, image.ZP, draw.Src)
	draw.Draw(myImage, AxisLng, &image.Uniform{black}, image.ZP, draw.Src)

	addLabel(myImage, 0, int(length*wh.Lng/100.0+1), fmt.Sprintf("Lat"))
	addLabel(myImage, int(length*wh.Lat/100.0+1), 0, fmt.Sprintf("Lng"))

	i := 0
	for _, order := range orders {
		x1 := length * order.Coords.Lat / 100.0
		y1 := length * order.Coords.Lng / 100.0

		redRect := image.Rect(int(x1-2), int(y1-2), int(x1+2), int(y1+2))

		// create a red rectangle atop the green surface
		draw.Draw(myImage, redRect, &image.Uniform{myred}, image.ZP, draw.Src)

		addLabel(myImage, int(x1+2), int(y1-2), fmt.Sprintf("%d - %d", order.ID, i))
		i++
	}

	//	addLine(myImage, x1, y1, x2, y2)

	png.Encode(outputFile, myImage)
}

func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{200, 100, 0, 255}
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}
