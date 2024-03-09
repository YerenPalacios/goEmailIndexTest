package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

func getMapContent(rawContent string) map[string]any {
	var item = make(map[string]any)
	contentLines := strings.Split(rawContent, "\n")

	bodyIndex := 0
	hasMultiLineKey := ""
	for index, line := range contentLines {
		splitLine := strings.Split(line, ": ")
		nextLine := contentLines[index+1]
		if len(splitLine) > 2 {
			splitLine[1] = strings.Join(splitLine[1:], ": ")
		}
		if len(splitLine) > 1 {
			if splitLine[0] == "Date" {
				currentFormat := "Mon, _2 Jan 2006 15:04:05 -0700 (MST)"
				date, err := time.Parse(currentFormat, splitLine[1])
				if err != nil {
					fmt.Println(err)
					date = time.Now()
				}
				item[splitLine[0]] = date
			} else {
				item[splitLine[0]] = splitLine[1]
			}

			if strings.HasPrefix(nextLine, "\t") {
				hasMultiLineKey = splitLine[0]
			}

		}

		if strings.HasPrefix(line, "\t") && hasMultiLineKey != "" {
			item[hasMultiLineKey] = item[hasMultiLineKey].(string) + strings.Replace(line, "\t", "", -1)
		}

		if nextLine == "" {
			bodyIndex = index + 1
			break
		}
	}

	item["content"] = strings.Join(contentLines[bodyIndex:], "\n")
	item["_id"] = item["Message-ID"] // TODO: review nil
	return item
}

func getMapFile(header *tar.Header, tarReader *tar.Reader) (map[string]any, error) {
	name := header.Name

	var contentBuffer bytes.Buffer

	var item map[string]any

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
		item["fileName"] = header.Name
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

func GetFileAsMapList(tarFile *gzip.Reader) ([]map[string]any, error) {
	tarReader := tar.NewReader(tarFile)
	var fileContents []map[string]any

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
