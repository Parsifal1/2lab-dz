package collectData

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path"
)

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

//walkFiles - обход всех файлов внутри указанной директории
func (f *FileCollector) WalkFiles(filepath string) (err error) {
	var files []os.FileInfo
	var fileReader *os.File

	if files, err = ioutil.ReadDir(filepath); err != nil {
		return
	}

	for i := range files {

		fullPath := path.Join(filepath, files[i].Name())

		if files[i].IsDir() {
			if err = f.WalkFiles(fullPath); err != nil {
				return
			}
			continue
		}
		f.addMeta(fullPath)

		if fileReader, err = os.Open(fullPath); err != nil {
			return
		}

		if err = f.PackFile(fullPath, fileReader); err != nil {
			return
		}
	}

	return
}

func (f *FileCollector) PackFile(filename string, fileReader io.Reader) (err error) {
	var fileWriter io.Writer

	if fileWriter, err = f.Zip.Create(filename); err != nil {
		return
	}

	if _, err = io.Copy(fileWriter, fileReader); err != nil {
		return
	}

	return
}

func (f *FileCollector) Meta2json() (js []byte, err error) {
	return json.Marshal(f.MetaData)
}

func (f *FileCollector) addMeta(fullPath string) {

	f.MetaData = append(f.MetaData, &FileMeta{
		Name: fullPath,
	})

	return
}

func (f *FileCollector) ZipData() (data []byte, err error) {

	if err = f.Zip.Close(); err != nil {
		return
	}

	data = f.ZipBuf.Bytes()
	return
}
