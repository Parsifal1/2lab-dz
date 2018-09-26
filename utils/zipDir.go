package zipDir

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"log"
	"os"
)

// ZipFiles - Сжатие в один архив нескольких файлов
// zipName - параметр для названия zip архива
// filesDir - имя папки с файлами, которые нужно добавить в архив
func ZipFiles(zipName string, filesDir string) error {
	// вытаскиваем все файлы из папки в массив files, сохраняя их в типе os.FileInfo
	files, err := ioutil.ReadDir(filesDir)
	if err != nil {
		log.Fatal(err)
	}
	// создаем переменную newZipFile с заданным именем : тип File
	newZipFile, err := os.Create(zipName)
	if err != nil {
		return err
	}
	defer newZipFile.Close() //в конце программы закрываем создание

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Добавление файлов в архив
	for _, file := range files {
		filename := filesDir + "/" + file.Name()
		zippingfile, err := os.Open(filename) //открываем архивируемый файл для:
		if err != nil {
			return err
		}
		defer zippingfile.Close() //закрываем архивируемый файл

		info, err := zippingfile.Stat() //получаем fileInfo архивируемого файла
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info) //из полученного fileInfo делаем fileHeader
		if err != nil {
			return err
		}

		header.Name = filename //записываем название файла в fileHeader

		header.Method = zip.Deflate //записываем метод сжатия файла в fileHeader

		writer, err := zipWriter.CreateHeader(header) //создаем Writer из полученных данных в fileHeader для записи в архив
		if err != nil {
			return err
		}
		if _, err = io.Copy(writer, zippingfile); err != nil { //копируем файл в полученный Writer
			return err
		}
	}
	return nil
}
