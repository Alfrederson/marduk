package main

import (
	"os"

	"github.com/Alfrederson/crebitos/pilar"
	"github.com/Alfrederson/crebitos/viga"
)

func main() {
	modo := os.Args[1]
	switch modo {
	case "viga":
		viga.Viga(os.Args[2:])
	case "pilar":
		pilar.Pilar(os.Args[2])
	}
}
