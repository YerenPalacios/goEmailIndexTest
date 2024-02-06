package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/import_file", handleImportFile)
	http.ListenAndServe(":8000", r)
}

func handleImportFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error reading file", http.StatusBadGateway)
	}
	fmt.Println(handler.Filename, file)
	gzf, err := gzip.NewReader(file)
	if err != nil {
		http.Error(w, "Error opening file", http.StatusBadGateway)
	}
	contentListMap, err := getFileAsMapList(gzf)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error reading file", http.StatusBadGateway)
	} else {
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error transforming file to json", http.StatusBadGateway)
		} else {
			err = sendtoZyncsearch(contentListMap)
			if err != nil {
				http.Error(w, "Error sending data", http.StatusBadGateway)

			} else {
				w.Write([]byte("OK"))
			}

		}
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
		lineInfo := filter(contentLines, func(i string) bool {
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

func filter(numbers []string, condition func(string) bool) []string {
	var result []string

	for _, num := range numbers {
		if condition(num) {
			result = append(result, num)
		}
	}

	return result
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

func getBatch(list []map[string]string, size int) [][]map[string]string {
	var grupos [][]map[string]string

	for i := 0; i < len(list); i += size {
		fin := i + size

		if fin > len(list) {
			fin = len(list)
		}

		grupos = append(grupos, list[i:fin])
	}

	return grupos
}

func send(body []byte) error {
	payload := bytes.NewBuffer(body)

	req, err := http.NewRequest(
		"POST",
		"http://172.26.32.1:4080/api/_bulkv2/",
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
	for _, batch := range getBatch(body, 1000) {
		requestBody := map[string]interface{}{
			"index":   "Messages2",
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
