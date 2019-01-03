package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"math/rand"
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

	drawAxis(myImage, wh, length)
	drawWh(myImage, wh, length)
	drawRoute(myImage, orders, myred, length)

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

func drawAxis(myImage *image.RGBA, wh *types.LatLng, length float64) {
	black := color.RGBA{0, 100, 0, 255}

	drawLine(myImage, length*wh.Lat/100.0, length*wh.Lng/100.0, 0, length*wh.Lng/100.0+1, black)
	drawLine(myImage, length*wh.Lat/100.0, length*wh.Lng/100.0, length*wh.Lat/100.0+1, 0, black)

	addLabel(myImage, 0, int(length*wh.Lng/100.0+1), fmt.Sprintf("Lat"))
	addLabel(myImage, int(length*wh.Lat/100.0+1), 0, fmt.Sprintf("Lng"))
}

func drawWh(myImage *image.RGBA, wh *types.LatLng, length float64) {
	myblue := color.RGBA{0, 0, 200, 255}

	Wh := image.Rect(int(length*wh.Lat/100.0-2), int(length*wh.Lng/100.0-2), int(length*wh.Lat/100.0+2), int(length*wh.Lng/100.0+2))
	draw.Draw(myImage, Wh, &image.Uniform{myblue}, image.ZP, draw.Src)
}

func drawRoute(myImage *image.RGBA, orders []*types.Order, col color.RGBA, length float64) {
	i := 0
	for _, order := range orders {
		x1 := length * order.Coords.Lat / 100.0
		y1 := length * order.Coords.Lng / 100.0

		redRect := image.Rect(int(x1-2), int(y1-2), int(x1+2), int(y1+2))

		draw.Draw(myImage, redRect, &image.Uniform{col}, image.ZP, draw.Src)
		addLabel(myImage, int(x1+2), int(y1-2), fmt.Sprintf("%d - %d", order.ID, i))
		i++
	}
}

func drawRoute2(myImage *image.RGBA, route *Route, col color.RGBA, length float64) {
	route.SolveTsp()
	i := 0
	for _, order := range route.List {
		x1 := length * order.Coords.Lat / 100.0
		y1 := length * order.Coords.Lng / 100.0

		redRect := image.Rect(int(x1-2), int(y1-2), int(x1+2), int(y1+2))

		draw.Draw(myImage, redRect, &image.Uniform{col}, image.ZP, draw.Src)
		addLabel(myImage, int(x1+2), int(y1-2), fmt.Sprintf("%d - %d", order.ID, i))
		if i > 0 {
			drawLine(myImage, length*route.List[i-1].Coords.Lat/100.0, length*route.List[i-1].Coords.Lng/100.0, x1, y1, col)
		}
		i++
	}
}

func drawRouteSet(filename string, routeset RouteSet) {
	length := 500.0
	myImage := image.NewRGBA(image.Rect(0, 0, int(length), int(length)))
	outputFile, _ := os.Create(filename)
	defer outputFile.Close()

	wh := routeset.Wareheouse

	drawAxis(myImage, wh, length)
	drawWh(myImage, wh, length)
	for _, route := range routeset.List {
		randColor := color.RGBA{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), 255}
		drawRoute2(myImage, route, randColor, length)
	}

	png.Encode(outputFile, myImage)
}

func drawLine(myImage *image.RGBA, x1, y1, x2, y2 float64, col color.RGBA) {
	y := y1

	start := math.Min(x1, x2)
	end := math.Max(x1, x2)

	for x := int(start); x < int(end); x++ {
		y = (float64(x)-x1)*(y2-y1)/(x2-x1) + y1
		myImage.Set(x, int(y), &image.Uniform{col})
	}
}
