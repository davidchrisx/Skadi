package main

import (
	"log"
	"os"
	"github.com/davidchrisx/skadi/proccess"
	"github.com/davidchrisx/skadi/sheet"
)

func main() {
	f, err := os.Open(os.Args[1])
	log.Println("parsing: ", os.Args[1])
	if err != nil {
		log.Fatalf("unable to open file: %s", err)
	}
	defer f.Close()

	m, err := proccess.Run(f)

	if err != nil {
		log.Fatalf("unable to create parser: %s", err)
	}
	err = sheet.Run(m)

	if err != nil {
		log.Fatalf("unabe to write to sheet: %s", err)
	}

}
