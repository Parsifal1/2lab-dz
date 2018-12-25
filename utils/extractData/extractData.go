package extractData

import (
	"archive/zip"
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"lab1/utils/checkSzp"
	"lab1/utils/getMeta"
	"os"
	"path/filepath"

	"github.com/fullsailor/pkcs7"
)

func Extract(destination string, hash string) error {
	var sign *pkcs7.PKCS7
	var fileMetas []getMeta.FileMeta
	var err error

	if sign, err = checkSzp.CheckSzp("./szip.szp", hash); err != nil {
		return err
	}

	if fileMetas, err = getMeta.GetMeta(sign); err != nil {
		return err
	}

	metaSize := int32(binary.LittleEndian.Uint32(sign.Content[:4]))

	archivedFiles := bytes.NewReader(sign.Content[4+metaSize:])

	err = UnarchiveFiles(archivedFiles, fileMetas, destination)
	if err != nil {
		return err
	}
	return nil
}

func UnarchiveFiles(archive *bytes.Reader, fileMetas []getMeta.FileMeta, destination string) error {
	zipReader, err := zip.NewReader(archive, archive.Size())
	if err != nil {
		return err
	}

	// Creating folder to extract to
	if err = os.MkdirAll(destination, 077); err != nil {
		fmt.Println("Couldn't create a folder to extract to")
		return err
	}

	for _, file := range zipReader.File {
		fileInfo := file.FileInfo()
		dirName, _ := filepath.Split(fileInfo.Name())

		if dirName != "" {
			if err = os.MkdirAll(filepath.Join(destination, "/", dirName), 077); err != nil {
				fmt.Println("Couldn't extract a folder")
				return err
			}
		}

		accessFile, err := file.Open() // gives io.ReadCloser
		if err != nil {
			fmt.Println("Unable to access a file")
			return err
		}

		fileGuts, err := ioutil.ReadAll(accessFile) // read file's bytes to buffer
		if err != nil {
			fmt.Println("Unable to read a file")
			return err
		}

		// Verifying hash for each file
		for _, metaData := range fileMetas {
			if metaData.Name == fileInfo.Name() {
				if metaData.Sha1Hash != sha1.Sum(fileGuts) {
					return errors.New(filepath.Join(file.Name, "'s hash is corrupted. The archive can't be fully unszipped"))
				}
			}
		}

		if err = ioutil.WriteFile(filepath.Join(destination, "/", fileInfo.Name()), fileGuts, 077); err != nil {
			return err
		}
	}

	return nil
}
