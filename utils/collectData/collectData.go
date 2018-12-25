package collectData

import (
	"archive/zip"
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path"
)

//FileMeta - единица передачи метаданных файла
type FileMeta struct {
	Name           string   `json: "filename"`
	OriginalSize   uint64   `json:"original_size"`
	CompressedSize uint64   `json:"compressed_size"`
	ModTime        string   `json:"mod_time"`
	Sha1Hash       [20]byte `json:"sha1_hash"`
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
	var header *zip.FileHeader
	var fileBytes []byte

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
		} else {
			if header, err = zip.FileInfoHeader(files[i]); err != nil {
				return
			}

			if fileBytes, err = ioutil.ReadFile(fullPath); err != nil {
				return
			}

			f.addMeta(header, fullPath, fileBytes)

			if fileReader, err = os.Open(fullPath); err != nil {
				return
			}

			if err = f.PackFile(fullPath, fileReader); err != nil {
				return
			}
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

func (f *FileCollector) addMeta(header *zip.FileHeader, fullPath string, fileBytes []byte) {

	f.MetaData = append(f.MetaData, &FileMeta{
		Name:           fullPath,
		OriginalSize:   header.UncompressedSize64,
		CompressedSize: header.CompressedSize64,
		ModTime:        header.Modified.Format("Mon Jan 2 15:04:05 MST 2006"),
		Sha1Hash:       sha1.Sum(fileBytes)})

	return
}

func (f *FileCollector) ZipData() (data []byte, err error) {

	if err = f.Zip.Close(); err != nil {
		return
	}

	data = f.ZipBuf.Bytes()
	return
}
