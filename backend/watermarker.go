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
	rootURI string = "/watermark/"
	maxSize uint   = 1024
	folder  string = "./dist"
	host    string = ":3210"
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
		http.StripPrefix(rootURI, fs).ServeHTTP(w, r)
	}
}

func process(w http.ResponseWriter, r *http.Request) error {
	// Get multipart form readers
	baseReader, _, err := r.FormFile("imageFile")
	if err != nil {
		log.Printf("imageFile opening error: %s", err)
		return err
	}
	tileReader, _, err := r.FormFile("watermarkFile")
	if err != nil {
		log.Printf("watermarkFile opening error: %s", err)
		return err
	}

	// Extract images
	base, err := imageutils.OpenImageReader(baseReader)
	if err != nil {
		log.Printf("Error on decoding base image: %s", err)
		return err
	}
	tile, err := imageutils.OpenImageReader(tileReader)
	if err != nil {
		log.Printf("Error on decoding tile image: %s", err)
		return err
	}

	// Process images
	result, err := imageutils.ResizeNWatermark(base, tile, maxSize)
	if err != nil {
		log.Printf("Error on image processing: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	// Send response
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Disposition", "attachment; filename=result.jpg")
	w.WriteHeader(http.StatusOK)

	if err = imageutils.WriteJpegImageWriter(result, w); err != nil {
		log.Printf("Error on sending image back: %s", err)
		return err
	}

	return nil
}

func runServer() {
	fs := http.FileServer(http.Dir(folder))
	mimeFs := MimeFileServer(fs)

	http.HandleFunc(rootURI, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			mimeFs(w, r)
		case http.MethodPost:
			process(w, r)
		default:
			log.Printf("Unsupported method: %s", r.Method)
		}
	})

	http.ListenAndServe(host, nil)
}
