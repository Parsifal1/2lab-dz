package main

import (
	"io/ioutil"
	"math/rand"

	"github.com/fogleman/gg"
	geojson "github.com/paulmach/go.geojson"
)

func main() { // Feature Collection
	var coordinates [][][][]float64
	var err error

	if coordinates, err = getMultyCoordinates(); err != nil {
		return
	}

	dc := gg.NewContext(2366, 2024)
	//рисуем MultyPolygon
	for i := 0; i < len(coordinates); i++ {
		dc.SetRGB(rand.Float64(), rand.Float64(), rand.Float64())
		drawPolygon(dc, coordinates[i][0], 10)
	}

	//рисуем контуры
	for i := 0; i < len(coordinates); i++ {
		dc.SetRGB(rand.Float64(), rand.Float64(), rand.Float64())
		drawLine(dc, coordinates[i][0], 10)
	}

	dc.SavePNG("out.png")
}

func drawPolygon(dc *gg.Context, coordinates [][]float64, scale float64) {
	x0 := coordinates[0][0] * scale
	y0 := coordinates[0][1] * scale * 2.1

	x0 = revertX(x0, scale)
	y0 = float64(dc.Height()) - y0

	dc.MoveTo(x0, y0)

	for index := 1; index < len(coordinates)-1; index++ {
		x := coordinates[index][0] * scale
		y := coordinates[index][1] * scale * 2.1

		x = revertX(x, scale)
		y = float64(dc.Height()) - y

		dc.LineTo(x, y)
	}

	dc.LineTo(x0, y0)
	dc.Fill()
}

func drawLine(dc *gg.Context, coordinates [][]float64, scale float64) {
	x0 := coordinates[0][0] * scale
	y0 := coordinates[0][1] * scale * 2.1

	x0 = revertX(x0, scale)
	y0 = float64(dc.Height()) - y0

	dc.MoveTo(x0, y0)

	for index := 1; index < len(coordinates)-1; index++ {
		x := coordinates[index][0] * scale
		y := coordinates[index][1] * scale * 2.1

		x = revertX(x, scale)
		y = float64(dc.Height()) - y

		dc.LineTo(x, y)
	}

	dc.LineTo(x0, y0)
	dc.SetLineWidth(10)
	dc.Stroke()
}

func getMultyCoordinates() ([][][][]float64, error) {
	var featureCollectionJSON []byte
	var filePath string
	var err error

	filePath = "rf.geojson"

	if featureCollectionJSON, err = ioutil.ReadFile(filePath); err != nil {
		return nil, err
	}
	var featureCollection *geojson.FeatureCollection

	if featureCollection, err = geojson.UnmarshalFeatureCollection(featureCollectionJSON); err != nil {
		return nil, err
	}

	var coordinates = featureCollection.Features[0].Geometry.MultiPolygon

	return coordinates, nil
}

func revertX(x float64, scale float64) float64 {
	if x < 0 {
		x = x / scale
		x = 360 + x
		x = x * scale
	}
	return x
}
