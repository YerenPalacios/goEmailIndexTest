package main

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"sync"
	"time"
)

// TODO: use environment variable
const ZYNCSARCH_URL = "http://172.26.32.1:4080"

func ImportFileService(file multipart.File) (string, error) {
	gzf, err := gzip.NewReader(file)
	if err != nil {
		fmt.Println(err)
		return "", errors.New("error opening file")
	}
	contentListMap, err := GetFileAsMapList(gzf)
	if err != nil {
		log.Println(err)
		return "", errors.New("error reading file")
	}
	if err != nil {
		log.Println(err)
		return "", errors.New("error reading file")
	}
	err = sendToZyncSearch(contentListMap)
	if err != nil {
		log.Println(err)
		return "", errors.New("error sending data: " + err.Error())
	} else {
		return "{\"status\": \"Ok\"}", nil
	}
}

func send(wg *sync.WaitGroup, body []byte, errorChannel chan error) {
	defer wg.Done()
	_, err := doPost(ZYNCSARCH_URL+"/api/_bulkv2", body)
	errorChannel <- err
}

const IndexName = "Messages"

func sendToZyncSearch(body []map[string]any) error {
	wg := new(sync.WaitGroup)
	errorChannel := make(chan error)
	defer close(errorChannel)

	batches := GetBatch(body, 1000)

	fmt.Println("Sending started at: ", time.Now(), "\n", len(batches), " Batches")

	for _, batch := range batches {
		wg.Add(1)
		requestBody := map[string]any{
			"index":   IndexName,
			"records": batch,
		}
		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			return err
		}
		go send(wg, jsonData, errorChannel)
	}
	var errorChannelList []error
	for range batches {
		errorChannelList = append(errorChannelList, <-errorChannel)
	}
	wg.Wait()
	var errorList []error
	for _, i := range errorChannelList {
		if i != nil {
			errorList = append(errorList, i)
		}
	}
	errorList = removeDuplication(errorList)
	if len(errorList) == 1 {
		return errorList[0]
	} else if len(errorList) > 1 {
		fmt.Println(errorList)
		return errors.New("different errors returned")
	}
	fmt.Println("Data sent at: ", time.Now())
	return nil
}
