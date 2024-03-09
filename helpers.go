package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

func GetBatch(list []map[string]any, size int) [][]map[string]any {
	var groups [][]map[string]any

	for i := 0; i < len(list); i += size {
		fin := i + size

		if fin > len(list) {
			fin = len(list)
		}

		groups = append(groups, list[i:fin])
	}

	return groups
}

func contains(slice []error, element error) bool {
	for _, value := range slice {
		if value.Error() == element.Error() {
			return true
		}
	}
	return false
}

func removeDuplication(list []error) []error {
	var newList []error

	for _, element := range list {
		if !contains(newList, element) {
			newList = append(newList, element)
		}
	}
	return newList
}

func doPost(url string, body []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(err)
		return []byte(""), err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic YWRtaW46Q29tcGxleHBhc3MjMTIz")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err, res)
		return []byte(""), err
	}
	defer res.Body.Close()
	content, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err, "Error reading response content")
		return []byte(""), err
	}
	if res.StatusCode == http.StatusOK {
		return content, nil
	} else {
		return []byte(""), errors.New(string(content))
	}
}
