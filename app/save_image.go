package app

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func FileSave(r *http.Request) string {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return ""
	}
	n := r.Form.Get("name")
	f, h, err := r.FormFile("file")
	if err != nil {
		return ""
	}
	defer f.Close()
	path := filepath.Join(".", "image")
	_ = os.MkdirAll(path, os.ModePerm)
	fullPath := path + "/" + n
	file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return ""
	}
	defer file.Close()
	// Copy the file to the destination path
	_, err = io.Copy(file, f)
	if err != nil {
		return ""
	}
	return n + filepath.Ext(h.Filename)
}

func writeResponse(w http.ResponseWriter, bytes []byte) {
	_, _ = w.Write(bytes)
}


func Respond(w http.ResponseWriter, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		RespondError(w, err, http.StatusInternalServerError)
		return
	}
	writeResponse(w, response)
}


func RespondError(w http.ResponseWriter, err error, code int) {
	response, err := json.Marshal(map[string]string{"error": err.Error()})
	w.WriteHeader(code)
	if err != nil {
		writeResponse(w, []byte(err.Error()))
		return
	}
	writeResponse(w, response)
}


func (a *App) PostFile(w http.ResponseWriter, r *http.Request) {
	path := FileSave(r)
	if path == "" {
		RespondError(w, errors.New("an error occurred"), http.StatusInternalServerError)
		return
	}
	Respond(w, map[string]string{"path": path})
}
