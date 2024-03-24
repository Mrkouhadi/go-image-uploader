package main

import (
	"fmt"
	"io"
	"os"
	"net/http"
	"path/filepath"
	"time"
)

func UploadFiles(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "index.html")

	// Define allowed extensions
	allowedExtensions := map[string]bool{
		".png":  true,
		".jpeg": true,
		".jpg":  true,
	}
	// 1. Parse input, type multipart/form-data
	r.Body = http.MaxBytesReader(w, r.Body, 32<<20+512) // 32 mb as maximum bytes to be read

	// 2. retrieve file from posted form-data
	file, handler, err := r.FormFile("myFile")
	if err != nil && err != http.ErrMissingFile {
		return ImageDetails{}, err // Return error if other than missing file
	} else if err == http.ErrMissingFile {
		// No file provided, return nil for ImageDetails
		return ImageDetails{}, nil
	}
	defer file.Close()

	fmt.Printf("Uploaded File Name : %+v \n", handler.Filename)
	fmt.Printf("File Size : %+v KB \n", handler.Size/1000)
	fmt.Printf("MIME Header : %+v \n", handler.Header)
	ext := filepath.Ext(handler.Filename)
	fmt.Printf("Extension of The file : %+v \n", ext)

	if !allowedExtensions[ext] {
		w.Write([]byte("this extension is not allowed \n"))
		return
	}

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

func MultipleFilesUploader(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 32 MB is the default used by FormFile()
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get a reference to the fileHeaders.
	// They are accessible only after ParseMultipartForm is called
	files := r.MultipartForm.File["img"]

	for _, fileHeader := range files {
		// Restrict the size of each uploaded file to 1MB.
		// To prevent the aggregate size from exceeding
		// a specified value, use the http.MaxBytesReader() method
		// before calling ParseMultipartForm()
		const MAX_UPLOAD_SIZE = 1024 * 1024 // 1MB

		if fileHeader.Size > MAX_UPLOAD_SIZE {
			http.Error(w, fmt.Sprintf("The uploaded image is too big: %s. Please use an image less than 1MB in size", fileHeader.Filename), http.StatusBadRequest)
			return
		}

		// Open the file
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer file.Close()
		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		filetype := http.DetectContentType(buff)
		if filetype != "image/jpeg" && filetype != "image/png" {
			http.Error(w, "The provided file format is not allowed. Please upload a JPEG or PNG image", http.StatusBadRequest)
			return
		}

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = os.MkdirAll("./uploads", os.ModePerm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		f, err := os.Create(fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		defer f.Close()

		_, err = io.Copy(f, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	fmt.Fprintf(w, "Upload successful")
}
