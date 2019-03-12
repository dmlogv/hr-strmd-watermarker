package main

import (
	"flag"
	"log"
	"mime"
	"net/http"
	"path/filepath"

	"./imageutils"
)

const (
	maxSize uint = 1024
)

func main() {
	server := flag.Bool("server", false, "Server mode")
	imagePathPtr := flag.String("image", "", "Image file path")
	watermarkPathPtr := flag.String("watermark", "", "Watermark file path")
	outputPath := flag.String("output", "", "Output file path")

	flag.Parse()

	if *server {
		log.Print("Server mode on")

		runServer()
	} else {
		if *imagePathPtr == "" || *watermarkPathPtr == "" || *outputPath == "" {
			log.Fatal("-image, -watermark or -output paths did not present")
		}

		base, _ := imageutils.OpenImage(*imagePathPtr)
		tile, _ := imageutils.OpenImage(*watermarkPathPtr)

		result, err := imageutils.ResizeNWatermark(base, tile, maxSize)
		if err != nil {
			log.Fatal(err)
		}

		if err = imageutils.WriteJpegImage(result, *outputPath); err != nil {
			log.Fatal(err)
		}
	}
}

// MimeFileServer implements FileServer with valid MIME-types
func MimeFileServer(fs http.Handler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		uri := r.RequestURI
		ext := filepath.Ext(uri)
		mimeType := mime.TypeByExtension(ext)

		if mimeType == "" {
			mimeType = "text/html"
		}

		log.Printf("%s %s %s", uri, ext, mimeType)

		w.Header().Set("Content-Type", mimeType)
		http.StripPrefix("/", fs).ServeHTTP(w, r)
	}
}

func runServer() {
	fs := http.FileServer(http.Dir("./dist"))
	mimeFs := MimeFileServer(fs)

	http.HandleFunc("/", mimeFs)

	http.ListenAndServe(":3210", nil)
}
