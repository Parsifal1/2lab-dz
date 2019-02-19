package main

import (
	"fmt"
	"io/ioutil"

	"github.com/davvo/mercator"
	"github.com/fogleman/gg"
	geojson "github.com/paulmach/go.geojson"
)

func main() { // Feature Collection
	var featureCollectionJSON []byte
	var filePath string
	var err error

	filePath = "geo.json"

	if featureCollectionJSON, err = ioutil.ReadFile(filePath); err != nil {
		return
	}

	var featureCollection *geojson.FeatureCollection

	if featureCollection, err = geojson.UnmarshalFeatureCollection(featureCollectionJSON); err != nil {
		return
	}

	var coordinates = featureCollection.Features[0].Geometry.Polygon[0]

	dc := gg.NewContext(1000, 1000)
	dc.SetRGB(1, 1, 0)
	drawPolygon(dc, coordinates, 3)

	dc.SetRGB(1, 0, 0)
	dc.DrawCircle(0, 0, 10)
	dc.DrawCircle(float64(dc.Width()), float64(dc.Height()), 10)
	dc.DrawCircle(float64(dc.Width()), 0, 10)
	dc.DrawCircle(0, float64(dc.Height()), 10)
	dc.SetLineWidth(100)
	dc.DrawCircle(0, float64(dc.Height()/2), 10)
	dc.DrawCircle(float64(dc.Width()/2), float64(dc.Height()/2), 10)
	dc.DrawCircle(float64(dc.Width()), float64(dc.Height()/2), 10)
	dc.Fill()

	dc.SavePNG("out.png")
}

func drawPolygon(dc *gg.Context, coordinates [][]float64, scale float64) {
	x0 := coordinates[0][0] * scale
	y0 := coordinates[0][1] * scale * 2

	y0 = float64(dc.Height()/2) - y0
	dc.MoveTo(x0, y0)
	for index := 1; index < len(coordinates)-1; index++ {
		x := coordinates[index][0] * scale
		y := coordinates[index][1] * scale * 2

		y = float64(dc.Height()/2) - y
		dc.LineTo(x, y)
	}
	dc.LineTo(x0, y0)
	dc.Fill()
}

func convert(x, y float64) (float64, float64) {
	var x0, y0, z = x, y, 2
	var px, py = mercator.LatLonToPixels(x, y, z)
	px = px / 10
	py = py / 10
	fmt.Printf("asd: %f, %f\n", x0, y0)
	fmt.Printf("Pixels (zoom %d): %f, %f\n", z, px, py)
	return px, py
}
