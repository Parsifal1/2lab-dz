package main

import (
	"flag"
	"fmt"
	"lab1/utils"
	"os"
	"runtime/debug"
)

func main() {

	fmt.Println("test")
	output := "done.zip"

	var mode string

	flag.StringVar(&mode, "mode", "z", "Режим работы приложения")

	flag.Parse()

	var err error

	switch mode {
	case "z":
		err = zipDir.ZipFiles(output, "filesDir")
	}

	if err != nil {
		fmt.Printf("Произошла неведомая ересь: %s\nПричина тут:\n%s", err, debug.Stack())
	}
}

//FileMeta - единица передачи метаданных файла
type FileMeta struct {
	Name string `json: "filename"`
}

func fileMeta(info os.FileInfo) (meta *FileMeta, err error) {

	meta = &FileMeta{
		Name: info.Name(),
	}

	return
}
