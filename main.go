package main

import (
	"log"

	"github.com/junos-streaming-analytics/app"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	app.Execute()
}
