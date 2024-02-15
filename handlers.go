package main

import (
	"net/http"
)

func HandleImportFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error reading file", http.StatusBadRequest)
		return
	}

	status, err := ImportFileService(file)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(status))
	}
}
