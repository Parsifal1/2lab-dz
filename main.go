package main

import (
	"flag"
	"fmt"
	"lab1/utils/extractData"
	"lab1/utils/zipData"
	"path/filepath"
	"runtime/debug"
)

func main() {

	var hash, mode, destination, source string

	flag.StringVar(&mode, "mode", "z", "Режим работы приложения")
	flag.StringVar(&hash, "hash", "UNDEF", "hash")
	flag.StringVar(&destination, "d", "./unszipped/", "destination to extract to")
	flag.StringVar(&source, "s", ".", "source of the archive")
	flag.Parse()

	var err error

	switch mode {
	case "z":
		err = zipData.ZipFiles()
	case "x":
		if err = extractData.Extract(destination, hash); err != nil {
			return
		}
		fmt.Println(filepath.Join("Your files have been successfully extracted to folder ", destination))
	default:
		fmt.Println("Uknown command. Please read manual and restart the application")
	}

	if err != nil {
		fmt.Printf("Произошла неведомая ересь: %s\nПричина тут:\n%s", err, debug.Stack())
	}
}
