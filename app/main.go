package main

import (
	"log"

	"github.com/holman_dw/ukjent/http"
	"github.com/holman_dw/ukjent/store/mem"
)

func main() {
	ms := mem.New()
	log.Fatal(http.Run(ms))
}
