package main

import (
	"flag"
	"log"
)

func main() {
	server := flag.Bool("server", false, "Server mode")
	imagePathPtr := flag.String("image", "", "Image file path")
	watermarkPathPtr := flag.String("watermark", "", "Watermark file path")

	if *server {
		log.Print("Server mode on")
	} else {
		if *imagePathPtr == "" || *watermarkPathPtr == "" {
			log.Panic("-image or -watermark paths did not present")
		}
	}
}
