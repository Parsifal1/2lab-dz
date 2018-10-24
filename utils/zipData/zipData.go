package zipData

import (
	"fmt"
	"io/ioutil"
	"lab1/utils/collectData"
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

	fmt.Println("Метаданные: ", js)

	var zipData []byte

	if zipData, err = collector.ZipData(); err != nil {
		return
	}

	if err = ioutil.WriteFile("archive.zip", zipData, 0644); err != nil {
		return
	}

	return
}
