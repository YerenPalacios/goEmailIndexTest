package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"sync"
	"time"
)

// TODO: use environment variable
const ZYNCSARCH_URL = "http://172.26.32.1:4080/api/_bulkv2/"

func ImportFileService(file multipart.File) (string, error) {
	gzf, err := gzip.NewReader(file)
	if err != nil {
		fmt.Println(err)
		return "", errors.New("error opening file")
	}
	contentListMap, err := getFileAsMapList(gzf)
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

func getMapContent(rawContent string) map[string]string {
	var item = make(map[string]string)
	contentLines := strings.Split(rawContent, "\n")

	fields := [16]string{
		"Message-ID",
		"Date",
		"From",
		"To",
		"Subject",
		"Cc",
		"Mime-Version",
		"Mime-Version",
		"Content-Transfer-Encoding",
		"Content-Transfer-Encoding",
		"X-From",
		"X-To",
		"X-cc",
		"X-Folder",
		"X-Origin",
		"X-FileName",
	}
	for _, v := range fields {
		lineInfo := Filter(contentLines, func(i string) bool {
			return strings.HasPrefix(i, v)
		})

		if len(lineInfo) > 0 {
			splitedline := strings.Split(lineInfo[0], ": ")

			if len(splitedline) > 1 {

				// obtains all "To" emails when are in more than one line
				if v == "To" {
					to := splitedline[1]
					index := 3
					for {
						nextline := contentLines[index+1]
						if !strings.HasPrefix(nextline, "\t") {
							break
						} else {
							to = to + strings.Replace(nextline, "\t", "", -1)
						}
						index++
					}
					item[v] = to
				} else {
					item[v] = splitedline[1]
				}

			}
		}
	}
	// TODO: remove lines that are no part of the message
	item["content"] = strings.Join(contentLines[15:], "\n")
	item["_id"] = item["Message-ID"]

	return item
}

func getMapFile(header *tar.Header, tarReader *tar.Reader) (map[string]string, error) {
	name := header.Name

	var contentBuffer bytes.Buffer

	var item map[string]string

	switch header.Typeflag {
	case tar.TypeDir:
		return item, errors.New("not file")
	case tar.TypeReg:
		if _, err := io.Copy(&contentBuffer, tarReader); err != nil {
			fmt.Println("Error...")
			return item, err
		}
		rawContent := contentBuffer.String()
		rawContent = strings.Replace(rawContent, "\r\n", "\n", -1)

		item = getMapContent(rawContent)
		return item, nil
	default:
		fmt.Printf("%s : %c %s %s\n",
			"Yikes! Unable to figure out type",
			header.Typeflag,
			"in file",
			name,
		)
		return item, errors.New("not file")
	}
}

func getFileAsMapList(tarFile *gzip.Reader) ([]map[string]string, error) {
	tarReader := tar.NewReader(tarFile)
	var fileContents []map[string]string

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println(err)
		}

		item, err := getMapFile(header, tarReader)

		if err != nil {
			if err.Error() == "not file" {
				continue
			} else {
				fmt.Println(err)

			}
		} else {
			fileContents = append(fileContents, item)
		}
	}
	return fileContents, nil
}

func send(wg *sync.WaitGroup, body []byte, errorChannel chan error) {
	defer wg.Done()
	payload := bytes.NewBuffer(body)

	req, err := http.NewRequest(
		"POST",
		ZYNCSARCH_URL,
		payload,
	)
	if err != nil {
		fmt.Println(err)
		errorChannel <- err
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic YWRtaW46Q29tcGxleHBhc3MjMTIz")

	client := &http.Client{}
	x, err := client.Do(req)

	if err != nil {
		log.Println(err, x)
		errorChannel <- err
		return
	}
	errorChannel <- nil
}

func sendToZyncSearch(body []map[string]string) error {
	wg := new(sync.WaitGroup)
	errorChannel := make(chan error)
	defer close(errorChannel)

	batches := GetBatch(body, 1000)

	fmt.Println("Sending started at: ", time.Now(), "\n", len(batches), " Batches")

	for _, batch := range batches {
		wg.Add(1)
		requestBody := map[string]interface{}{
			"index":   "Messages",
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
