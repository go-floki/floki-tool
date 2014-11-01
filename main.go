package main

import (
	"flag"
	"log"
	"os"
)

var (
	projectDir = flag.String("dir", ".", "path to the project directory")
)

func main() {
	flag.Usage = usage
	flag.Parse()

	models := ParseModels(*projectDir)

	RemoveControllerFiles(*projectDir, models)
	RemoveServiceFiles(*projectDir, models)

	serviceSymbols := ParseSymbols(*projectDir + "/services")
	controllerSymbols := ParseSymbols(*projectDir + "/controllers")

	log.Println(serviceSymbols)
	//log.Println(controllerSymbols)

	GenerateServices(*projectDir, models, serviceSymbols)
	GenerateControllers(*projectDir, models, controllerSymbols)
}

func usage() {
	log.Printf("usage: floki-tool [cmd] arguments")
	flag.PrintDefaults()
	os.Exit(2)
}
