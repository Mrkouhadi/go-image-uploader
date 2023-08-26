package main

import (
	"fmt"
	"io"
	"os"
	"net/http"
	"path/filepath"
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
	ext := filepath.Ext(handler.Filename)
	fmt.Printf("Extension of The file : %+v \n", ext)

	// 3. write the file on our server
	newName := "My-new-image"	
	newFile, err := os.Create("static/images/"+ newName + ext)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer newFile.Close()
	// fileBytes, err := ioutil.ReadAll(file)
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	newFile.Write(fileBytes)

	w.Write([]byte("You have Successfully uploaded the file \n"))
}
