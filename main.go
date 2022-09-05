package main

import (
	"log"

	"github.com/jhonnyesquivel/quasar-op/cmd/api/bootstrap"
)

func main() {
	if err := bootstrap.Run(); err != nil {
		log.Fatal(err)
	}
}
