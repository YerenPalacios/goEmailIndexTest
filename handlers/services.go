package handlers

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

const ZYNCSARCH_URL = "http://172.26.32.1:4080/api/_bulkv2/"

func ImportFileService(file multipart.File) (string, error) {
	gzf, err := gzip.NewReader(file)
	if err != nil {
		fmt.Println(err)
		return "", errors.New("error opening file")
	}
	contentListMap, err := getFileAsMapList(gzf)
	if err != nil {
		fmt.Println(err)
		return "", errors.New("error reading file")
	}
	if err != nil {
		fmt.Println(err)
		return "", errors.New("error reading file")
	}
	err = sendtoZyncsearch(contentListMap)
	if err != nil {
		fmt.Println(err)
		return "", errors.New("error sending data")
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
		// TODO: review files with emails in more than 1 line
		if len(lineInfo) > 0 {
			splitedline := strings.Split(lineInfo[0], ": ")
			if len(splitedline) > 1 {
				item[v] = splitedline[1]
			}
		}
	}
	// remove lines that are no part of the message
	item["content"] = strings.Join(contentLines[15:], "\n")
	item["_id"] = item["Message-ID"]

	return item
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

		name := header.Name

		var contentBuffer bytes.Buffer

		switch header.Typeflag {
		case tar.TypeDir:
			continue
		case tar.TypeReg:
			if _, err := io.Copy(&contentBuffer, tarReader); err != nil {
				fmt.Println("Error...")
				return fileContents, err
			}
			rawContent := contentBuffer.String()
			rawContent = strings.Replace(rawContent, "\r\n", "\n", -1)

			fileContents = append(fileContents, getMapContent(rawContent))
		default:
			fmt.Printf("%s : %c %s %s\n",
				"Yikes! Unable to figure out type",
				header.Typeflag,
				"in file",
				name,
			)
		}
	}
	return fileContents, nil
}

func send(body []byte) error {
	payload := bytes.NewBuffer(body)

	req, err := http.NewRequest(
		"POST",
		ZYNCSARCH_URL,
		payload,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic YWRtaW46Q29tcGxleHBhc3MjMTIz")

	client := &http.Client{}
	x, err := client.Do(req)

	if err != nil {
		fmt.Println("*****E ", err, x)
		return err
	}
	return nil
}

func sendtoZyncsearch(body []map[string]string) error {
	for _, batch := range GetBatch(body, 1000) {
		requestBody := map[string]interface{}{
			"index":   "Messages3",
			"records": batch,
		}
		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			return err
		}
		go send(jsonData)
	}
	return nil
}
