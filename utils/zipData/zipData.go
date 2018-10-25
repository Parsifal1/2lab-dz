package zipData

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"lab1/utils/collectData"
	"lab1/utils/signData"
)

//ZipFiles - создаем сборщик метаданных, передаем его методу перебора указанной директории,
//в котором также упаковываем файлы в архив, выводим метаданные в json, закрываем writer,
//сохраняем архив
func ZipFiles() (err error) {

	collector := collectData.NewFileCollector()

	if err = collector.WalkFiles("./filesDir"); err != nil {
		return
	}

	var js []byte

	if js, err = collector.Meta2json(); err != nil {
		return
	}

	metaCollector := collectData.NewFileCollector()

	if err = metaCollector.PackFile("meta.zip", bytes.NewReader(js)); err != nil {
		return
	}

	var metaZip []byte
	if metaZip, err = metaCollector.ZipData(); err != nil {
		return
	}

	fmt.Println("Метаданные: ", metaZip)

	var zipData []byte

	if zipData, err = collector.ZipData(); err != nil {
		return
	}

	resultBuf := new(bytes.Buffer)

	if err = binary.Write(resultBuf, binary.LittleEndian, uint32(len(metaZip))); err != nil {
		return
	}

	if _, err = resultBuf.Write(metaZip); err != nil {
		return
	}
	if _, err = resultBuf.Write(zipData); err != nil {
		return
	}

	var signedData []byte

	if signedData, err = signData.SignData(resultBuf.Bytes()); err != nil {
		return
	}

	if err = ioutil.WriteFile("archive.szp", signedData, 0644); err != nil {
		return
	}
	return
}
