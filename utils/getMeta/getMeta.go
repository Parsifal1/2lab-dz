package getMeta

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/fullsailor/pkcs7"
)

func GetMeta(p *pkcs7.PKCS7) ([]FileMeta, error) {
	//Read meta
	metaSize := int32(binary.LittleEndian.Uint32(p.Content[:4]))
	fmt.Println(metaSize)
	bytedMeta := bytes.NewReader(p.Content[4 : metaSize+4])

	readableMeta, err := zip.NewReader(bytedMeta, bytedMeta.Size())
	if err != nil {
		//fmt.Println("ошибка 2")
		return nil, err
	}

	if len(readableMeta.File) < 1 {
		return nil, errors.New("File doesn't have meta")
	}

	metaCompressed := readableMeta.File[0] //meta.xml

	metaUncompressed, err := metaCompressed.Open()
	if err != nil {
		//fmt.Println("ошибка 3")
		return nil, err
	}
	defer metaUncompressed.Close()

	var fileMetas []FileMeta
	metaUncompressedBody, err := ioutil.ReadAll(metaUncompressed)
	if err != nil {
		//fmt.Println("ошибка 4")
		return nil, err
	}
	err = xml.Unmarshal(metaUncompressedBody, &fileMetas)
	if err != nil {
		//fmt.Println("ошибка 4")
		return nil, err
	}

	return fileMetas, err
}

type FileMeta struct {
	Name           string   `json: "filename"`
	OriginalSize   uint64   `json:"original_size"`
	CompressedSize uint64   `json:"compressed_size"`
	ModTime        string   `json:"mod_time"`
	Sha1Hash       [20]byte `json:"sha1_hash"`
}
