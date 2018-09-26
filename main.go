package main

import (
	"Go/utils/"
	"log"
)

func main() {
	output := "done.zip"

	err := zipDir.ZipFiles(output, "filesDir")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Zipped File: " + output)
}
