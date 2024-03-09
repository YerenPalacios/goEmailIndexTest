package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

func HandleGetMessages(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	search := query.Get("search")
	from, err := strconv.Atoi(query.Get("from"))
	if err != nil {
		from = 0
	}
	to, err := strconv.Atoi(query.Get("to"))
	if err != nil {
		to = 10
	}

	service := MessagesService{}
	payload := service.createPayload(search, []string{"-Date"}, from, to)
	body, _ := json.Marshal(payload)
	fmt.Println(payload)
	result, err := service.getMessages(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Write(result)
}
