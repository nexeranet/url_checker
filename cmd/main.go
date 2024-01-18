package main

import (
	"log"

	"github.com/nexeranet/url_checker/pkg/urlchecker"
)

func main() {
	checker := urlchecker.NewURLChecker()
	log.Println("Lister and server :3000")
	if err := checker.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
