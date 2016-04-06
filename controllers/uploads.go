package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/denisbakhtin/blog/shared"
)

//Upload handles POST /upload route
func Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		err := r.ParseMultipartForm(32 << 20) // ~32MB
		if err != nil {
			log.Printf("ERROR: %s\n", err)
			http.Error(w, err.Error(), 500)
			return
		}
		var uris []string
		fmap := r.MultipartForm.File
		for k := range fmap {
			file, fileHeader, err := r.FormFile(k)
			if err != nil {
				log.Printf("ERROR: %s\n", err)
				http.Error(w, err.Error(), 500)
				return
			}
			uri, err := saveFile(fileHeader, file)
			if err != nil {
				log.Printf("ERROR: %s\n", err)
				http.Error(w, err.Error(), 500)
				return
			}
			uris = append(uris, uri)
		}
		json.NewEncoder(w).Encode(uris)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		http.Error(w, err.Error(), 405)
	}
}

//saveFile saves file to disc and returns its relative uri
func saveFile(fh *multipart.FileHeader, f multipart.File) (string, error) {
	fileExt := filepath.Ext(fh.Filename)
	newName := fmt.Sprint(time.Now().Unix()) + fileExt //unique file name ;D
	uri := "/public/uploads/" + newName
	fullName := filepath.Join(shared.UploadsPath(), newName)

	file, err := os.OpenFile(fullName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = io.Copy(file, f)
	if err != nil {
		return "", err
	}
	return uri, nil
}
