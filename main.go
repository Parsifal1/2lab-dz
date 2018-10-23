package main

import (
	"flag"
	"fmt"
	"lab1/utils/zipData"
	"runtime/debug"
)

func main() {

	output := "done.zip"

	var mode string

	flag.StringVar(&mode, "mode", "z", "Режим работы приложения")

	flag.Parse()

	var err error

	switch mode {
	case "z":
		err = zipData.ZipFiles(output, "filesDir")
	}

	if err != nil {
		fmt.Printf("Произошла неведомая ересь: %s\nПричина тут:\n%s", err, debug.Stack())
	}
}
