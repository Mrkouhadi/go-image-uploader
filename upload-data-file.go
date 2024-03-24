package main

import (
	"errors"
	"io"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// ////// TYPES
type ImageDetails struct {
	Filename  string
	Extension string
	Size      int64
	Header    textproto.MIMEHeader
}
type Employee struct {
	ID            uuid.UUID `json:"id"`   // github.com/google/uuid
	Role          string    `json:"role"` // admin, editor, author
	First_name    string    `json:"first_name"`
	Last_name     string    `json:"last_name"`
	Gender        string    `json:"gender"`
	Birth_date    time.Time `json:"birth_date"`
	Nationality   string    `json:"nationality"`
	Email         string    `json:"email"`
	Password      string    `json:"password"`
	Profile_image string    `json:"profile_image"`
	Created_at    time.Time `json:"created_at"`
	Updated_at    time.Time `json:"updated_at"`
}

///////////////////////////////// the handler

func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	role := r.FormValue("role")
	email := r.FormValue("email")
	gender := r.FormValue("gender")
	birthdate := r.FormValue("birth_date")
	nationality := r.FormValue("nationality")
	password := r.FormValue("password")

	birthDate, err := time.Parse("2006-01-02", birthdate)
	if err != nil {
		// errors handling here
		return
	}
	// hash the password before storing it in the database
	psw, err := HashPassword(password)
	if err != nil {
		// errors handling here
		return
	}
	user := Employee{
		ID:          uuid.New(),
		First_name:  firstName,
		Last_name:   lastName,
		Gender:      gender,
		Role:        role,
		Email:       email,
		Password:    psw,
		Nationality: nationality,
		Birth_date:  birthDate,
		Created_at:  time.Now(),
		Updated_at:  time.Now(),
	}
	_, err = UploadFile(w, r, role+"-"+lastName, &user)
	if err != nil {

		return
	}
	// store it in the databse

	// write json feedback
}

// ////////////////////////////////// upload a file fuction
func UploadFile(w http.ResponseWriter, r *http.Request, name string, employeeData *Employee) (ImageDetails, error) {
	// Define allowed extensions
	allowedExtensions := map[string]bool{
		".png":  true,
		".jpeg": true,
		".jpg":  true,
	}
	// 1. Parse input, type multipart/form-data
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20+512) // // max is 10MB
	// the commented below line of code is for alllowing user to upload multiple files, so we don't need that...
	// r.ParseMultipartForm(10 * 1024 * 1024)              // max is 10MB

	// 2. retrieve file from posted form-data
	file, handler, err := r.FormFile("img")
	if err != nil && err != http.ErrMissingFile {
		return ImageDetails{}, err // Return error if other than missing file
	} else if err == http.ErrMissingFile {
		// No file provided, return nil for ImageDetails
		return ImageDetails{}, nil
	}
	defer file.Close()
	ext := filepath.Ext(handler.Filename)
	// Check if the file extension is allowed
	if !allowedExtensions[ext] {
		return ImageDetails{}, errors.New("invalid file extension")
	}
	// 3. write temprorary file on our server
	newFile, err := os.Create("data/images/employees" + name + ext)
	if err != nil {
		return ImageDetails{}, err
	}
	defer newFile.Close()
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return ImageDetails{}, err
	}
	newFile.Write(fileBytes)

	// Fill employeeData struct
	employeeData.Profile_image = name + ext

	return ImageDetails{
		Filename:  handler.Filename,
		Extension: ext,
		Size:      handler.Size / 1000, // KB
		Header:    handler.Header,
	}, nil
}
//////////// html code
<form action="http://localhost:8080/employee/register-employee" method="POST" enctype="multipart/form-data">
      <input type="text" name="first_name" placeholder="First Name" />
      <input type="text" name="last_name" placeholder="Last Name" />
      <input type="text" name="role" placeholder="Role" />
      <input type="date" name="birth_date" placeholder="Birth Date" />
      <input type="text" name="email" placeholder="Email" />
      <input type="text" name="gender" placeholder="Gender" />
      <input type="password" name="password" placeholder="Password" />
      <input type="file" name="img" accept="image/png, image/jpeg, image/jpg"/>
      <input type="submit" value="Upload file.." />
    </form>