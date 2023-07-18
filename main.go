package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/upload", UploadFiles)
	fmt.Println("Listening to port:8080")
	http.ListenAndServe(":8080", nil)
}
