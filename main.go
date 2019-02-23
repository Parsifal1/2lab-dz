package main

import (
	"io/ioutil"
	"math/rand"

	"github.com/fogleman/gg"
	geojson "github.com/paulmach/go.geojson"
)

func main() {
	var coordinates [][][][][]float64
	var err error

	if coordinates, err = getMultyCoordinates(); err != nil {
		return
	}

	dc := gg.NewContext(1366, 1024)
	scale := 5.0

	//рисуем полигоны
	forEachPolygon(dc, coordinates, func(polygonCoordinates [][]float64) {
		dc.SetRGB(rand.Float64(), rand.Float64(), rand.Float64())
		drawByPolygonCoordinates(dc, polygonCoordinates, scale, dc.Fill)
	})
	//рисуем контуры полигонов
	forEachPolygon(dc, coordinates, func(polygonCoordinates [][]float64) {
		dc.SetLineWidth(3)
		dc.SetRGB(rand.Float64(), rand.Float64(), rand.Float64())
		drawByPolygonCoordinates(dc, polygonCoordinates, scale, dc.Stroke)
	})

	dc.SavePNG("out.png")
}

func getMultyCoordinates() ([][][][][]float64, error) {
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
	var features = featureCollection.Features
	var coordinates [][][][][]float64
	for i := 0; i < len(features); i++ {
		coordinates = append(coordinates, features[i].Geometry.MultiPolygon)
	}
	return coordinates, nil
}

func forEachPolygon(dc *gg.Context, coordinates [][][][][]float64, callback func([][]float64)) {
	for i := 0; i < len(coordinates); i++ {
		for j := 0; j < len(coordinates[i]); j++ {
			callback(coordinates[i][j][0])
		}
	}
}

func drawByPolygonCoordinates(dc *gg.Context, coordinates [][]float64, scale float64, method func()) {
	x0 := convertNegativeX(coordinates[0][0]) * scale
	y0 := coordinates[0][1] * scale * 2.1
	y0 = float64(dc.Height()) - y0
	dc.MoveTo(x0, y0)
	for index := 1; index < len(coordinates)-1; index++ {
		x := convertNegativeX(coordinates[index][0]) * scale
		y := coordinates[index][1] * scale * 2.1
		y = float64(dc.Height()) - y
		dc.LineTo(x, y)
	}
	dc.LineTo(x0, y0)
	method()
}

func convertNegativeX(x float64) float64 {
	if x < 0 {
		x = 360 + x
	}
	return x
}
