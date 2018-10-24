package zipData

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
)

//ZipFiles - создаем сборщик метаданных, передаем его методу перебора указанной директории,
//в котором также упаковываем файлы в архив, выводим метаданные в json, закрываем writer,
//сохраняем архив
func ZipFiles() (err error) {

	collector := NewFileCollector()

	if err = walkFiles(collector, "./filesDir"); err != nil {
		return
	}

	var js []byte

	if js, err = collector.meta2json(); err != nil {
		return
	}

	fmt.Println("Метаданные: ", js)

	var zipData []byte

	if zipData, err = collector.zipData(); err != nil {
		return
	}

	if err = ioutil.WriteFile("archive.zip", zipData, 0644); err != nil {
		return
	}

	return
}

func walkFiles(collector *FileCollector, filepath string) (err error) {
	var files []os.FileInfo

	if files, err = ioutil.ReadDir(filepath); err != nil {
		return
	}

	for i := range files {

		fullPath := path.Join(filepath, files[i].Name())

		if files[i].IsDir() {
			if err = walkFiles(collector, fullPath); err != nil {
				return
			}
			continue
		}
		collector.addMeta(fullPath)

		var fileReader *os.File

		if fileReader, err = os.Open(fullPath); err != nil {
			return
		}

		if err = collector.packFile(fullPath, fileReader); err != nil {
			return
		}
	}

	return
}

//FileMeta - единица передачи метаданных файла
type FileMeta struct {
	Name string // `json: "filename"`
}

//FileCollector - для сбора итогового файла
type FileCollector struct {
	ZipBuf   *bytes.Buffer
	Zip      *zip.Writer
	MetaData []*FileMeta
}

//NewFileCollector - конструктор FileCollector
func NewFileCollector() *FileCollector {
	buf := new(bytes.Buffer)
	return &FileCollector{
		ZipBuf:   buf,
		Zip:      zip.NewWriter(buf),
		MetaData: make([]*FileMeta, 0, 100),
	}
}

func (f *FileCollector) meta2json() (js []byte, err error) {
	return json.Marshal(f.MetaData)
}

func (f *FileCollector) addMeta(fullPath string) {

	f.MetaData = append(f.MetaData, &FileMeta{
		Name: fullPath,
	})

	return
}

func (f *FileCollector) packFile(filename string, fileReader io.Reader) (err error) {
	var fileWriter io.Writer

	if fileWriter, err = f.Zip.Create(filename); err != nil {
		return
	}

	if _, err = io.Copy(fileWriter, fileReader); err != nil {
		return
	}

	return
}

func (f *FileCollector) zipData() (data []byte, err error) {

	if err = f.Zip.Close(); err != nil {
		return
	}

	data = f.ZipBuf.Bytes()
	return
}
