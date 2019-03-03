package main

import (
	"fmt"           // пакет для форматированного ввода вывода
	"html/template" // пакет для логирования
	"math/rand"
	"net/http" // пакет для поддержки HTTP протокола
	"strconv"

	"github.com/fogleman/gg"
	geojson "github.com/paulmach/go.geojson"
	// пакет для работы с  UTF-8 строками
)

var imgSrc string

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./index.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	fmt.Println(imgSrc)

	t.ExecuteTemplate(w, "index", imgSrc)

}

func drawHandler(w http.ResponseWriter, r *http.Request) {
	content := r.FormValue("content")

	var featureCollectionJSON []byte
	var err error

	featureCollectionJSON = []byte(content)

	if imgSrc, err = getPNG(featureCollectionJSON); err != nil {
		fmt.Println(err.Error())
	}

	http.Redirect(w, r, "/", 302)
}

func main() {

	imgSrc = ""

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/draw", drawHandler)

	http.ListenAndServe(":3000", nil)
}

func getPNG(featureCollectionJSON []byte) (string, error) {
	var coordinates [][][][][]float64
	var err error

	if coordinates, err = getMultyCoordinates(featureCollectionJSON); err != nil {
		return "", err
	}

	dc := gg.NewContext(1366, 1024)
	scale := 5.0

	//рисуем полигоны
	forEachPolygon(dc, coordinates, func(polygonCoordinates [][]float64, i int, j int) {
		dc.SetRGB(rand.Float64(), rand.Float64(), rand.Float64())
		drawByPolygonCoordinates(dc, polygonCoordinates, scale, dc.Fill)
	})
	//рисуем контуры полигонов
	forEachPolygon(dc, coordinates, func(polygonCoordinates [][]float64, i int, j int) {
		dc.SetRGB(rand.Float64(), rand.Float64(), rand.Float64())
		dc.SetLineWidth(3)
		drawByPolygonCoordinates(dc, polygonCoordinates, scale, dc.Stroke)
	})

	var out = strconv.Itoa(rand.Intn(10000)) + ".png"
	out = "assets/" + out
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

func forEachPolygon(dc *gg.Context, coordinates [][][][][]float64, callback func([][]float64, int, int)) {
	for i := 0; i < len(coordinates); i++ {
		for j := 0; j < len(coordinates[i]); j++ {
			callback(coordinates[i][j][0], i, j)
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
