package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"

	"github.com/davvo/mercator"
	"github.com/fogleman/gg"
	geojson "github.com/paulmach/go.geojson"
)

func main() {
	var featureCollectionJSON []byte
	var filePath = "rf.geojson"
	var err error

	if featureCollectionJSON, err = ioutil.ReadFile(filePath); err != nil {
		fmt.Println(err.Error())
	}

	if _, err = getPNG(featureCollectionJSON); err != nil {
		fmt.Println(err.Error())
	}
}

func getPNG(featureCollectionJSON []byte) (string, error) {
	var coordinates [][][][][]float64
	var err error

	if coordinates, err = getMultyCoordinates(featureCollectionJSON); err != nil {
		return "", err
	}

	dc := gg.NewContext(width, height)
	scale := 5.0

	dc.InvertY()

	//рисуем полигоны
	forEachPolygon(dc, coordinates, func(polygonCoordinates [][]float64) {
		dc.SetRGB(rand.Float64(), rand.Float64(), rand.Float64())
		drawByPolygonCoordinates(dc, polygonCoordinates, scale, dc.Fill)
	})
	//рисуем контуры полигонов
	dc.SetLineWidth(3)
	forEachPolygon(dc, coordinates, func(polygonCoordinates [][]float64) {
		dc.SetRGB(rand.Float64(), rand.Float64(), rand.Float64())
		drawByPolygonCoordinates(dc, polygonCoordinates, scale, dc.Stroke)
	})

	var out = strconv.Itoa(rand.Intn(10000)) + ".png"

	dc.SavePNG(out)

	return out, nil
}

func getMultyCoordinates(featureCollectionJSON []byte) ([][][][][]float64, error) {
	var featureCollection *geojson.FeatureCollection
	var err error

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

const width, height = 256, 256
const mercatorMaxValue float64 = 20037508.342789244

const mercatorToCanvasScaleFactorX = float64(width) / (mercatorMaxValue)
const mercatorToCanvasScaleFactorY = float64(height) / (mercatorMaxValue)

func drawByPolygonCoordinates(dc *gg.Context, coordinates [][]float64, scale float64, method func()) {
	for index := 0; index < len(coordinates)-1; index++ {
		x, y := mercator.LatLonToMeters(coordinates[index][1], coordinates[index][0])

		x, y = centerRussia(x, y)

		x *= mercatorToCanvasScaleFactorX
		y *= mercatorToCanvasScaleFactorY

		dc.LineTo(x, y)
	}
	dc.ClosePath()
	method()
}

func centerRussia(x float64, y float64) (float64, float64) {
	var west = float64(1635093.15883866)

	if x > 0 {
		x -= west
	} else {
		x += 2*mercatorMaxValue - west
	}

	return x, y
}
