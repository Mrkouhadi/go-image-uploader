package main

import (
	"fmt"
	"io"
	"os"
	"net/http"
)

func UploadFiles(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "index.html")

	// 1. Parse input, type multipart/form-data
	r.Body = http.MaxBytesReader(w, r.Body, 32<<20+512)
	// r.ParseMultipartForm(32 << 20) // means max 32 MB    32 << 20 means max 10 MB
	r.ParseMultipartForm(10 * 1024 * 1024) // max is 10MB

	// 2. retrieve file from posted form-data
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Failed to retrieve file from form-data")
		fmt.Println(err)
		return
	}
	defer file.Close()

	fmt.Printf("Uploaded File Name : %+v \n", handler.Filename)
	fmt.Printf("File Size : %+v KB \n", handler.Size/1000)
	fmt.Printf("MIME Header : %+v \n", handler.Header)

	// 3. write temprorary file on our server
	// tempoFile, err := ioutil.TempFile("tempo-images", "uploaded-*.jpeg")
	tempoFile, err := os.CreateTemp("tempo-images", "uploaded-*.jpeg")

	if err != nil {
		fmt.Println(err)
		return
	}
	defer tempoFile.Close()
	// fileBytes, err := ioutil.ReadAll(file)
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	tempoFile.Write(fileBytes)

	w.Write([]byte("You have Successfully uploaded the file \n"))
}
